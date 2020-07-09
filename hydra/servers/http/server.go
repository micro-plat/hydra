package http

import (
	"context"
	"fmt"
	xnet "net"
	x "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/types"
)

//Server api服务器
type Server struct {
	*option
	server  *x.Server
	engine  *gin.Engine
	running bool
	ip      string
	proto   string
	host    string
	port    string
}

//NewServer 创建http api服务嚣
func NewServer(name string, addr string, routers []*router.Router, opts ...Option) (t *Server, err error) {
	t, err = new(name, addr, opts...)
	if err != nil {
		return
	}
	t.addHttpRouters(routers...)
	return
}

//NewWSServer 创建web socket服务嚣
func NewWSServer(name string, addr string, routers []*router.Router, opts ...Option) (t *Server, err error) {
	t, err = new(name, addr, opts...)
	if err != nil {
		return
	}
	t.proto = "ws"
	t.addWSRouters(routers...)
	return
}

//new 创建http api服务嚣
func new(name string, addr string, opts ...Option) (t *Server, err error) {
	t = &Server{
		proto: "http",
		ip:    net.GetLocalIPAddress(),
		option: &option{
			readHeaderTimeout: 6,
			readTimeout:       6,
			writeTimeout:      6,
			metric:            middleware.NewMetric(),
		},
	}
	for _, opt := range opts {
		opt(t.option)
	}
	t.host, t.port, err = global.GetHostPort(addr)
	if err != nil {
		return nil, err
	}
	t.server = &x.Server{
		Addr:              xnet.JoinHostPort(t.host, t.port),
		ReadHeaderTimeout: time.Second * time.Duration(t.option.readHeaderTimeout),
		ReadTimeout:       time.Second * time.Duration(t.option.readTimeout),
		WriteTimeout:      time.Second * time.Duration(t.option.writeTimeout),
		MaxHeaderBytes:    1 << 20,
	}
	return
}

// Start the http server
func (s *Server) Start() error {
	s.running = true
	errChan := make(chan error, 1)
	switch len(s.tls) {
	case 2:
		s.proto = "https"
		go func(ch chan error) {
			if err := s.server.ListenAndServeTLS(s.tls[0], s.tls[1]); err != nil {
				ch <- err
			}
		}(errChan)
	default:
		go func(ch chan error) {
			if err := s.server.ListenAndServe(); err != nil {
				ch <- err
			}
		}(errChan)

	}
	select {
	case <-time.After(time.Millisecond * 500):
		return nil
	case err := <-errChan:
		s.running = false
		return err
	}
}

//Shutdown 关闭服务器
func (s *Server) Shutdown() error {
	if s.server != nil && s.running {
		s.running = false
		defer s.metric.Stop()
		ctx, cannel := context.WithTimeout(context.Background(), time.Second*10)
		defer cannel()
		if err := s.server.Shutdown(ctx); err != nil {
			if err == x.ErrServerClosed {
				return nil
			}
			return fmt.Errorf("关闭出现错误:%w", err)
		}
	}
	return nil
}

//GetAddress 获取当前服务地址
func (s *Server) GetAddress(h ...string) string {
	if len(h) > 0 && h[0] != "" {
		return fmt.Sprintf("%s://%s:%s", s.proto, h[0], s.port)
	}
	if s.host == "0.0.0.0" {
		return fmt.Sprintf("%s://%s:%s", s.proto, s.ip, s.port)
	}
	return fmt.Sprintf("%s://%s:%s", s.proto, s.host, s.port)
}

//GetStatus 获取当前服务器状态
func (s *Server) GetStatus() string {
	return types.DecodeString(s.running, true, "运行中", "停止")
}
