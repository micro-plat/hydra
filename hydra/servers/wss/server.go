package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	componentwss "github.com/micro-plat/hydra/components/wss"
	wssconf "github.com/micro-plat/hydra/conf/server/wss"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/lib4go/logger"
)

type Server struct {
	conf    *wssconf.Server
	server  *http.Server
	pool    *tunnelPool
	engine  *adapter.DispatcherEngine
	routes  []wssconf.Route
	log     logger.ILogger
	running bool
}

var activeServer *Server

func NewServer(conf *wssconf.Server, engine *adapter.DispatcherEngine, routes []wssconf.Route, log logger.ILogger) *Server {
	if log == nil {
		log = logger.New("wss.server")
	}
	return &Server{conf: conf, engine: engine, routes: routes, pool: newTunnelPool(), log: log}
}

func (s *Server) Start() error {
	host, port, err := global.GetHostPort(s.conf.GetAddress())
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveHTTP)
	s.server = &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           mux,
		ReadHeaderTimeout: time.Second * time.Duration(s.conf.GetWriteTimeout()),
	}
	errCh := make(chan error, 1)
	s.running = true
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		s.running = false
		return err
	case <-time.After(500 * time.Millisecond):
		activeServer = s
		componentwss.Register(func(method string, path string, query string, body []byte, header map[string]string) (int, map[string]string, []byte, error) {
			return s.Invoke(method, path, query, body, header)
		})
		componentwss.RegisterStream(func(req componentwss.Request, onChunk componentwss.ChunkHandler) (*componentwss.Response, error) {
			return s.InvokeStream(req, onChunk)
		})
		componentwss.RegisterOpenStream(func(req componentwss.Request, onStart componentwss.StartHandler, onChunk componentwss.ChunkHandler) (*componentwss.Response, error) {
			return s.InvokeOpenStream(req, onStart, onChunk)
		})
		return nil
	}
}

func (s *Server) Shutdown() {
	if s.pool != nil {
		s.pool.closeAll()
	}
	if s.server != nil && s.running {
		s.running = false
		if activeServer == s {
			activeServer = nil
			componentwss.Register(nil)
			componentwss.RegisterStream(nil)
			componentwss.RegisterOpenStream(nil)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

func (s *Server) GetAddress(hosts ...string) string {
	host, port, _ := global.GetHostPort(s.conf.GetAddress())
	if len(hosts) > 0 && hosts[0] != "" {
		host = hosts[0]
	}
	if host == "0.0.0.0" || host == "" {
		host = global.LocalIP()
	}
	return fmt.Sprintf("ws://%s:%s%s", host, port, s.conf.GetPath())
}

func (s *Server) ServiceNum() int {
	routers := s.engine.Routes()
	serverMap := map[string]string{}
	for _, item := range routers {
		if _, ok := serverMap[item.Path]; !ok {
			serverMap[item.Path] = item.Path
		}
	}
	return len(serverMap)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == s.conf.GetPath() {
		s.serveTunnel(w, r)
		return
	}
	if websocket.IsWebSocketUpgrade(r) {
		s.serveBusinessWSS(w, r)
		return
	}
	s.serveProxy(w, r)
}

func (s *Server) serveTunnel(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r, s.conf.AuthType, s.conf.AuthSecret) {
		s.log.Warn("wss.server.auth.failed", r.URL.Path, r.RemoteAddr)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Warn("wss.server.upgrade.failed", r.RemoteAddr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.log.Info("wss.server.connected", r.URL.Path, r.RemoteAddr)
	go s.readTunnel(conn)
}

func (s *Server) readTunnel(conn *websocket.Conn) {
	var sess *session
	defer func() {
		if sess != nil {
			s.pool.remove(sess)
			sess.close()
			s.log.Warn("wss.server.disconnected", conn.RemoteAddr(), sess.group, sess.id)
		} else {
			conn.Close()
			s.log.Warn("wss.server.disconnected", conn.RemoteAddr(), "unregistered")
		}
	}()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		frame := &Frame{}
		if err := json.Unmarshal(data, frame); err != nil {
			continue
		}
		switch frame.Type {
		case frameHello, frameRegister:
			if sess != nil {
				s.log.Warn("wss.server.register.duplicate", conn.RemoteAddr(), sess.group, sess.id)
				s.writeRaw(conn, &Frame{Type: frameError, Error: "session already registered"})
				continue
			}
			id := frame.Client
			if id == "" {
				id = uuid.New().String()
			}
			if frame.Group == "" {
				s.log.Warn("wss.server.register.failed", conn.RemoteAddr(), "missing group", id)
				s.writeRaw(conn, &Frame{Type: frameError, Error: "missing group"})
				continue
			}
			sess = newSession(id, frame.Group, conn)
			s.pool.add(sess)
			sess.write(&Frame{Type: frameRegistered, Client: id, Group: frame.Group})
			s.log.Info("wss.server.registered", conn.RemoteAddr(), frame.Group, id)
			go s.heartbeat(sess)
		case frameResponse:
			s.pool.done(frame.ID, frame)
		case frameStreamStart, frameStreamChunk:
			s.pool.send(frame.ID, frame)
		case frameStreamEnd:
			s.pool.done(frame.ID, frame)
		case framePong:
			if sess != nil {
				sess.pong()
			}
		}
	}
}

func (s *Server) heartbeat(sess *session) {
	interval := time.Duration(s.conf.GetPingInterval()) * time.Second
	timeout := time.Duration(s.conf.GetPongTimeout()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-sess.closed:
			return
		case <-ticker.C:
			if !sess.alive(timeout) {
				s.log.Warn("wss.server.heartbeat.timeout", sess.group, sess.id, timeout)
				s.pool.remove(sess)
				sess.close()
				return
			}
			if err := sess.write(&Frame{Type: framePing, ID: uuid.New().String()}); err != nil {
				s.log.Warn("wss.server.heartbeat.failed", sess.group, sess.id, err)
				s.pool.remove(sess)
				sess.close()
				return
			}
		}
	}
}

func (s *Server) writeRaw(conn *websocket.Conn, frame *Frame) {
	data, _ := json.Marshal(frame)
	conn.WriteMessage(websocket.BinaryMessage, data)
}

func (s *Server) serveBusinessWSS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()
	for {
		tp, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if tp != websocket.TextMessage && tp != websocket.BinaryMessage {
			continue
		}
		begin := time.Now()
		fullPath := pathWithQuery(r.URL.Path, r.URL.RawQuery)
		s.log.Info("wss.server.request:", http.MethodGet, fullPath, "from", r.RemoteAddr)
		status, _, body, err := dispatch(s.engine, http.MethodGet, r.URL.Path, r.URL.RawQuery, data, headers(r.Header))
		if err != nil {
			s.log.Error("wss.server.response:", http.MethodGet, fullPath, status, "ws", costSince(begin), err)
			resp, _ := json.Marshal(map[string]interface{}{"code": status, "err": err.Error()})
			conn.WriteMessage(websocket.TextMessage, resp)
			continue
		}
		s.log.Info("wss.server.response:", http.MethodGet, fullPath, statusOr(status, http.StatusOK), "ws", costSince(begin))
		conn.WriteMessage(websocket.TextMessage, body)
	}
}

func (s *Server) serveProxy(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	if wantsStream(r, body) {
		s.serveProxyStream(w, r, body)
		return
	}
	status, header, body, err := s.Invoke(r.Method, r.URL.Path, r.URL.RawQuery, body, headers(r.Header), r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), statusOr(status, http.StatusBadGateway))
		return
	}
	copyStringHeader(w.Header(), header)
	w.WriteHeader(statusOr(status, http.StatusOK))
	w.Write(body)
}

func (s *Server) serveProxyStream(w http.ResponseWriter, r *http.Request, body []byte) {
	begin := time.Now()
	fullPath := pathWithQuery(r.URL.Path, r.URL.RawQuery)
	s.log.Info("wss.server.request:", r.Method, fullPath, "from", r.RemoteAddr)
	group, path, ok := s.resolve(r)
	if !ok {
		writer := newHTTPStreamWriter(w)
		status, _, err := dispatchStream(s.engine, r.Method, r.URL.Path, r.URL.RawQuery, body, headers(r.Header), writer)
		if err != nil {
			if !writer.Written() {
				http.Error(w, err.Error(), statusOr(status, http.StatusInternalServerError))
			}
			s.log.Error("wss.server.response:", r.Method, fullPath, statusOr(status, http.StatusInternalServerError), "local", costSince(begin), err)
			return
		}
		s.log.Info("wss.server.response:", r.Method, fullPath, statusOr(status, http.StatusOK), "local", costSince(begin))
		return
	}
	sess, err := s.pool.pick(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		s.log.Error("wss.server.response:", r.Method, fullPath, http.StatusServiceUnavailable, group, costSince(begin), err)
		return
	}
	id := uuid.New().String()
	ch := s.pool.wait(id)
	defer s.pool.cancel(id)
	s.log.Info("wss.server.dispatch:", shortID(id), r.Method, fullPath, "to", group+":"+path, "client", sess.id)
	if err := sess.write(&Frame{Type: frameRequest, ID: id, Group: group, Method: r.Method, Path: path, Query: r.URL.RawQuery, Header: headers(r.Header), Body: body, Extra: map[string]interface{}{"stream": true}}); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		s.log.Error("wss.server.response:", r.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
		return
	}
	timeout := time.Duration(s.conf.GetRequestTimeout()) * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	started := false
	status := http.StatusOK
	for {
		select {
		case frame, ok := <-ch:
			if !ok || frame == nil {
				err := fmt.Errorf("request canceled")
				if !started {
					http.Error(w, err.Error(), http.StatusGatewayTimeout)
				}
				s.log.Error("wss.server.response:", r.Method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
				return
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(timeout)
			switch frame.Type {
			case frameStreamStart:
				copyStringHeader(w.Header(), frame.Header)
				status = statusOr(frame.Status, http.StatusOK)
				if !started {
					w.WriteHeader(status)
					started = true
				}
			case frameStreamChunk:
				if !started {
					status = statusOr(frame.Status, http.StatusOK)
					w.WriteHeader(status)
					started = true
				}
				if len(frame.Body) > 0 {
					if _, err := w.Write(frame.Body); err != nil {
						s.log.Error("wss.server.response:", r.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
						return
					}
					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
				}
			case frameStreamEnd:
				status = statusOr(frame.Status, status)
				if frame.Error != "" {
					err := fmt.Errorf(frame.Error)
					if !started {
						http.Error(w, err.Error(), statusOr(status, http.StatusBadGateway))
					}
					s.log.Error("wss.server.response:", r.Method, fullPath, statusOr(status, http.StatusBadGateway), group, costSince(begin), err)
					return
				}
				if !started {
					copyStringHeader(w.Header(), frame.Header)
					w.WriteHeader(statusOr(status, http.StatusOK))
				}
				s.log.Info("wss.server.response:", r.Method, fullPath, statusOr(status, http.StatusOK), group, costSince(begin))
				return
			case frameResponse:
				copyStringHeader(w.Header(), frame.Header)
				status = statusOr(frame.Status, http.StatusOK)
				if frame.Error != "" {
					err := fmt.Errorf(frame.Error)
					http.Error(w, err.Error(), statusOr(status, http.StatusBadGateway))
					s.log.Error("wss.server.response:", r.Method, fullPath, statusOr(status, http.StatusBadGateway), group, costSince(begin), err)
					return
				}
				w.WriteHeader(status)
				w.Write(frame.Body)
				s.log.Info("wss.server.response:", r.Method, fullPath, status, group, costSince(begin))
				return
			}
		case <-timer.C:
			err := fmt.Errorf("request timeout")
			if !started {
				http.Error(w, err.Error(), http.StatusGatewayTimeout)
			}
			s.log.Error("wss.server.response:", r.Method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
			return
		case <-r.Context().Done():
			err := r.Context().Err()
			s.log.Warn("wss.server.response:", r.Method, fullPath, http.StatusRequestTimeout, group, costSince(begin), err)
			return
		}
	}
}

func (s *Server) Invoke(method string, requestPath string, query string, body []byte, header map[string]string, remote ...string) (int, map[string]string, []byte, error) {
	begin := time.Now()
	fullPath := pathWithQuery(requestPath, query)
	from := "internal"
	if len(remote) > 0 && remote[0] != "" {
		from = remote[0]
	}
	s.log.Info("wss.server.request:", method, fullPath, "from", from)
	group, path, ok := s.resolvePath("", requestPath)
	if !ok {
		status, respHeader, respBody, err := dispatch(s.engine, method, requestPath, query, body, header)
		if err != nil {
			s.log.Error("wss.server.response:", method, fullPath, status, "local", costSince(begin), err)
			return status, nil, nil, err
		}
		s.log.Info("wss.server.response:", method, fullPath, statusOr(status, http.StatusOK), "local", costSince(begin))
		return status, headers(respHeader), respBody, nil
	}
	sess, err := s.pool.pick(group)
	if err != nil {
		s.log.Error("wss.server.response:", method, fullPath, http.StatusServiceUnavailable, group, costSince(begin), err)
		return http.StatusServiceUnavailable, nil, nil, err
	}
	id := uuid.New().String()
	ch := s.pool.wait(id)
	defer s.pool.cancel(id)
	s.log.Info("wss.server.dispatch:", shortID(id), method, fullPath, "to", group+":"+path, "client", sess.id)
	if err := sess.write(&Frame{Type: frameRequest, ID: id, Group: group, Method: method, Path: path, Query: query, Header: header, Body: body}); err != nil {
		s.log.Error("wss.server.response:", method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
		return http.StatusBadGateway, nil, nil, err
	}
	select {
	case resp, ok := <-ch:
		if !ok || resp == nil {
			err := fmt.Errorf("request canceled")
			s.log.Error("wss.server.response:", method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
			return http.StatusGatewayTimeout, nil, nil, err
		}
		if resp.Error != "" {
			err := fmt.Errorf(resp.Error)
			s.log.Error("wss.server.response:", method, fullPath, statusOr(resp.Status, http.StatusBadGateway), group, costSince(begin), err)
			return statusOr(resp.Status, http.StatusBadGateway), resp.Header, nil, err
		}
		s.log.Info("wss.server.response:", method, fullPath, statusOr(resp.Status, http.StatusOK), group, costSince(begin))
		return statusOr(resp.Status, http.StatusOK), resp.Header, resp.Body, nil
	case <-time.After(time.Duration(s.conf.GetRequestTimeout()) * time.Second):
		err := fmt.Errorf("request timeout")
		s.log.Error("wss.server.response:", method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
		return http.StatusGatewayTimeout, nil, nil, err
	}
}

func (s *Server) InvokeStream(req componentwss.Request, onChunk componentwss.ChunkHandler) (*componentwss.Response, error) {
	return s.InvokeOpenStream(req, nil, onChunk)
}

func (s *Server) InvokeOpenStream(req componentwss.Request, onStart componentwss.StartHandler, onChunk componentwss.ChunkHandler) (*componentwss.Response, error) {
	begin := time.Now()
	fullPath := pathWithQuery(req.Path, req.Query)
	if req.Method == "" {
		req.Method = http.MethodGet
	}
	if req.Header == nil {
		req.Header = componentwss.Header{}
	}
	s.log.Info("wss.server.request:", req.Method, fullPath, "from", "internal")
	group, path, ok := s.resolvePath("", req.Path)
	if !ok {
		writer := newComponentStreamWriter(onStart, onChunk)
		status, header, err := dispatchStream(s.engine, req.Method, req.Path, req.Query, req.Body, map[string]string(req.Header), writer)
		resp := &componentwss.Response{Status: statusOr(status, http.StatusOK), Header: componentwss.Header(headers(header))}
		if err != nil {
			s.log.Error("wss.server.response:", req.Method, fullPath, statusOr(status, http.StatusInternalServerError), "local", costSince(begin), err)
			return resp, err
		}
		s.log.Info("wss.server.response:", req.Method, fullPath, resp.Status, "local", costSince(begin))
		return resp, nil
	}
	sess, err := s.pool.pick(group)
	if err != nil {
		s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusServiceUnavailable, group, costSince(begin), err)
		return &componentwss.Response{Status: http.StatusServiceUnavailable}, err
	}
	id := uuid.New().String()
	ch := s.pool.wait(id)
	defer s.pool.cancel(id)
	s.log.Info("wss.server.dispatch:", shortID(id), req.Method, fullPath, "to", group+":"+path, "client", sess.id)
	if err := sess.write(&Frame{Type: frameRequest, ID: id, Group: group, Method: req.Method, Path: path, Query: req.Query, Header: map[string]string(req.Header), Body: req.Body, Extra: map[string]interface{}{"stream": true}}); err != nil {
		s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
		return &componentwss.Response{Status: http.StatusBadGateway}, err
	}
	timeout := time.Duration(s.conf.GetRequestTimeout()) * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	resp := &componentwss.Response{Status: http.StatusOK, Header: componentwss.Header{}}
	for {
		select {
		case frame, ok := <-ch:
			if !ok || frame == nil {
				err := fmt.Errorf("request canceled")
				s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
				resp.Status = http.StatusGatewayTimeout
				return resp, err
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(timeout)
			switch frame.Type {
			case frameStreamStart:
				resp.Status = statusOr(frame.Status, http.StatusOK)
				resp.Header = componentwss.Header(frame.Header)
				if onStart != nil {
					if err := onStart(resp); err != nil {
						s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
						return resp, err
					}
				}
			case frameStreamChunk:
				if len(frame.Body) > 0 && onChunk != nil {
					if err := onChunk(frame.Body); err != nil {
						s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
						return resp, err
					}
				}
			case frameStreamEnd:
				resp.Status = statusOr(frame.Status, resp.Status)
				if len(frame.Header) > 0 {
					resp.Header = componentwss.Header(frame.Header)
				}
				if frame.Error != "" {
					err := fmt.Errorf(frame.Error)
					s.log.Error("wss.server.response:", req.Method, fullPath, statusOr(resp.Status, http.StatusBadGateway), group, costSince(begin), err)
					return resp, err
				}
				s.log.Info("wss.server.response:", req.Method, fullPath, statusOr(resp.Status, http.StatusOK), group, costSince(begin))
				return resp, nil
			case frameResponse:
				resp.Status = statusOr(frame.Status, http.StatusOK)
				resp.Header = componentwss.Header(frame.Header)
				resp.Body = frame.Body
				if frame.Error != "" {
					err := fmt.Errorf(frame.Error)
					s.log.Error("wss.server.response:", req.Method, fullPath, statusOr(resp.Status, http.StatusBadGateway), group, costSince(begin), err)
					return resp, err
				}
				if onStart != nil {
					if err := onStart(resp); err != nil {
						s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
						return resp, err
					}
				}
				if len(frame.Body) > 0 && onChunk != nil {
					if err := onChunk(frame.Body); err != nil {
						s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusBadGateway, group, costSince(begin), err)
						return resp, err
					}
				}
				s.log.Info("wss.server.response:", req.Method, fullPath, resp.Status, group, costSince(begin))
				return resp, nil
			}
		case <-timer.C:
			err := fmt.Errorf("request timeout")
			s.log.Error("wss.server.response:", req.Method, fullPath, http.StatusGatewayTimeout, group, costSince(begin), err)
			resp.Status = http.StatusGatewayTimeout
			return resp, err
		}
	}
}

func (s *Server) resolve(r *http.Request) (group string, path string, ok bool) {
	host := stripPort(r.Host)
	return s.resolvePath(host, r.URL.Path)
}

func (s *Server) resolvePath(host string, requestPath string) (group string, path string, ok bool) {
	for _, route := range s.routes {
		if route.Host != "" && !strings.EqualFold(stripPort(route.Host), host) {
			continue
		}
		if route.PathPrefix != "" && !matchPrefix(requestPath, route.PathPrefix) {
			continue
		}
		if route.Group == "" {
			continue
		}
		path = requestPath
		if route.StripPrefix != "" && matchPrefix(path, route.StripPrefix) {
			path = strings.TrimPrefix(path, route.StripPrefix)
			if path == "" {
				path = "/"
			}
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
		}
		return route.Group, path, true
	}
	parts := strings.Split(strings.Trim(requestPath, "/"), "/")
	if len(parts) < 2 || parts[0] == "" {
		return "", "", false
	}
	return parts[0], "/" + strings.Join(parts[1:], "/"), true
}

func Call(method string, path string, query string, body []byte, header map[string]string) (int, map[string]string, []byte, error) {
	if activeServer == nil {
		return http.StatusServiceUnavailable, nil, nil, fmt.Errorf("wss.server is not started")
	}
	return activeServer.Invoke(method, path, query, body, header)
}

func stripPort(host string) string {
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return strings.Trim(h, "[]")
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return strings.Trim(host, "[]")
	}
	if strings.Count(host, ":") > 1 {
		return host
	}
	if i := strings.LastIndex(host, ":"); i > -1 {
		return host[:i]
	}
	return host
}

func checkAuth(r *http.Request, authType string, secret string) bool {
	if secret == "" {
		return true
	}
	token := getAuthToken(r.Header, r.URL.Query(), authType)
	return token == secret
}

func setAuthHeader(header http.Header, authType string, secret string) {
	if secret == "" {
		return
	}
	switch strings.ToLower(authType) {
	case "apikey", "api-key":
		header.Set("X-API-Key", secret)
	case "bearer":
		header.Set("Authorization", "Bearer "+secret)
	default:
		header.Set("Authorization", "Hydra-WSS "+secret)
		header.Set("X-Hydra-WSS-Token", secret)
	}
}

func getAuthToken(header http.Header, query url.Values, authType string) string {
	switch strings.ToLower(authType) {
	case "apikey", "api-key":
		return firstNonEmpty(header.Get("X-API-Key"), header.Get("X-Hydra-WSS-Token"), query.Get("token"))
	case "bearer":
		return trimAuthPrefix(header.Get("Authorization"), "Bearer")
	default:
		auth := header.Get("Authorization")
		if token := trimAuthPrefix(auth, "Hydra-WSS"); token != "" {
			return token
		}
		return firstNonEmpty(header.Get("X-Hydra-WSS-Token"), header.Get("X-API-Key"), query.Get("token"))
	}
}

func trimAuthPrefix(auth string, scheme string) string {
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], scheme) {
		return ""
	}
	return parts[1]
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func matchPrefix(path string, prefix string) bool {
	if prefix == "" || prefix == "/" {
		return true
	}
	if path == prefix {
		return true
	}
	return strings.HasPrefix(path, strings.TrimRight(prefix, "/")+"/")
}

func readBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	data, _ := io.ReadAll(r.Body)
	return data
}

func headers(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}

func copyHeader(dst http.Header, src http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}

func copyStringHeader(dst http.Header, src map[string]string) {
	for k, v := range src {
		dst.Set(k, v)
	}
}

func statusOr(status int, def int) int {
	if status <= 0 {
		return def
	}
	return status
}

func pathWithQuery(path string, query string) string {
	if query == "" {
		return path
	}
	return path + "?" + query
}

func shortID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

func costSince(start time.Time) string {
	d := time.Since(start)
	switch {
	case d >= time.Second:
		return strconv.FormatFloat(float64(d)/float64(time.Second), 'f', 3, 64) + "s"
	case d >= time.Millisecond:
		return strconv.FormatFloat(float64(d)/float64(time.Millisecond), 'f', 3, 64) + "ms"
	case d >= time.Microsecond:
		return strconv.FormatFloat(float64(d)/float64(time.Microsecond), 'f', 3, 64) + "us"
	default:
		return strconv.FormatInt(int64(d), 10) + "ns"
	}
}
