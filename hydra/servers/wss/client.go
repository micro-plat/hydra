package wss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	wssconf "github.com/micro-plat/hydra/conf/server/wss"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/lib4go/logger"
)

type Client struct {
	conf    *wssconf.Server
	engine  *adapter.DispatcherEngine
	conn    *websocket.Conn
	stop    chan struct{}
	running bool
	id      string
	mu      sync.Mutex
	log     logger.ILogger
}

func NewClient(conf *wssconf.Server, engine *adapter.DispatcherEngine, log logger.ILogger) *Client {
	if log == nil {
		log = logger.New("wss.client")
	}
	return &Client{conf: conf, engine: engine, stop: make(chan struct{}), log: log}
}

func (c *Client) Start() error {
	if c.conf.Server == "" {
		return fmt.Errorf("wss.client server is required")
	}
	if c.conf.Group == "" {
		return fmt.Errorf("wss.client group is required")
	}
	c.running = true
	c.log.Info("wss.client.start", c.conf.Server, c.conf.Group, c.clientID())
	go c.loop()
	return nil
}

func (c *Client) GetAddress() string {
	return c.conf.Server
}

func (c *Client) ServiceNum() int {
	routers := c.engine.Routes()
	serverMap := map[string]string{}
	for _, item := range routers {
		if _, ok := serverMap[item.Path]; !ok {
			serverMap[item.Path] = item.Path
		}
	}
	return len(serverMap)
}

func (c *Client) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	close(c.stop)
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) loop() {
	for {
		select {
		case <-c.stop:
			return
		default:
		}
		if err := c.connectAndServe(); err != nil {
			c.log.Warn("wss.client.connect.failed", c.conf.Group, c.clientID(), c.conf.GetReconnect(), err)
			time.Sleep(time.Second * time.Duration(c.conf.GetReconnect()))
		}
	}
}

func (c *Client) connectAndServe() error {
	u, err := url.Parse(c.conf.Server)
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Set("X-Hydra-Group", c.conf.Group)
	header.Set("X-Hydra-Client-ID", c.clientID())
	setAuthHeader(header, c.conf.AuthType, c.conf.AuthSecret)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}
	c.conn = conn
	c.log.Info("wss.client.connected", c.conf.Group, c.clientID())
	defer func() {
		conn.Close()
		c.log.Warn("wss.client.disconnected", c.conf.Group, c.clientID())
	}()
	if err := c.writeTo(conn, &Frame{Type: frameRegister, Group: c.conf.Group, Client: c.clientID()}); err != nil {
		return err
	}
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		frame := &Frame{}
		if err := json.Unmarshal(data, frame); err != nil {
			continue
		}
		switch frame.Type {
		case frameRegistered:
			c.log.Info("wss.client.registered", frame.Group, frame.Client)
		case frameRequest:
			go c.handle(conn, frame)
		case framePing:
			c.writeTo(conn, &Frame{Type: framePong, ID: frame.ID})
		case frameError:
			c.log.Warn("wss.client.error", c.conf.Group, c.clientID(), frame.Error)
		}
	}
}

func (c *Client) handle(conn *websocket.Conn, frame *Frame) {
	if isStreamFrame(frame) {
		c.handleStream(conn, frame)
		return
	}
	begin := time.Now()
	fullPath := pathWithQuery(frame.Path, frame.Query)
	c.log.Info("wss.client.request:", shortID(frame.ID), frame.Method, fullPath, frame.Group)
	c.log.Info("wss.client.dispatch:", shortID(frame.ID), frame.Method, fullPath, "local")
	status, header, body, err := dispatch(c.engine, frame.Method, frame.Path, frame.Query, frame.Body, frame.Header)
	resp := &Frame{Type: frameResponse, ID: frame.ID, Status: status, Header: headers(header), Body: body}
	if err != nil {
		resp.Error = err.Error()
	}
	if writeErr := c.writeTo(conn, resp); writeErr != nil {
		c.log.Error("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusInternalServerError), costSince(begin), writeErr)
		return
	}
	if err != nil {
		c.log.Error("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusInternalServerError), costSince(begin), err)
		return
	}
	c.log.Info("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusOK), costSince(begin))
}

func (c *Client) handleStream(conn *websocket.Conn, frame *Frame) {
	begin := time.Now()
	fullPath := pathWithQuery(frame.Path, frame.Query)
	c.log.Info("wss.client.request:", shortID(frame.ID), frame.Method, fullPath, frame.Group)
	c.log.Info("wss.client.dispatch:", shortID(frame.ID), frame.Method, fullPath, "local")
	writer := newFrameStreamWriter(frame.ID, func(resp *Frame) error {
		return c.writeTo(conn, resp)
	})
	status, _, err := dispatchStream(c.engine, frame.Method, frame.Path, frame.Query, frame.Body, frame.Header, writer)
	if writeErr := writeStreamEnd(writer, status, err); writeErr != nil {
		c.log.Error("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusInternalServerError), costSince(begin), writeErr)
		return
	}
	if err != nil {
		c.log.Error("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusInternalServerError), costSince(begin), err)
		return
	}
	c.log.Info("wss.client.response:", frame.Method, fullPath, shortID(frame.ID), statusOr(status, http.StatusOK), costSince(begin))
}

func (c *Client) writeTo(conn *websocket.Conn, frame *Frame) error {
	if conn == nil {
		return fmt.Errorf("wss client connection is not ready")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := json.Marshal(frame)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *Client) clientID() string {
	if c.id != "" {
		return c.id
	}
	if c.conf.ClientID != "" {
		c.id = c.conf.ClientID
		return c.id
	}
	c.id = uuid.New().String()[:8]
	return c.id
}
