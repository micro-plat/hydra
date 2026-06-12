package aigw

import (
	"context"
	"fmt"
	xnet "net"
	xhttp "net/http"
	"time"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/types"
)

// Server AIGW服务
type Server struct {
	*option
	server  *xhttp.Server
	running bool
	ip      string
	proto   string
	host    string
	port    string
	engine  *adapter.GinEngine
}

// NewServer 创建AIGW服务
func NewServer(name string, addr string, routers []*router.Router, opts ...Option) (t *Server, err error) {
	t = &Server{
		proto: "http",
		ip:    global.LocalIP(),
		option: &option{
			readHeaderTimeout: 30,
			readTimeout:       0,
			writeTimeout:      0,
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
	t.server = &xhttp.Server{
		Addr:              xnet.JoinHostPort(t.host, t.port),
		ReadHeaderTimeout: time.Second * time.Duration(t.option.readHeaderTimeout),
		MaxHeaderBytes:    1 << 20,
	}
	if t.option.readTimeout > 0 {
		t.server.ReadTimeout = time.Second * time.Duration(t.option.readTimeout)
	}
	if t.option.writeTimeout > 0 {
		t.server.WriteTimeout = time.Second * time.Duration(t.option.writeTimeout)
	}
	t.addAIRouters(routers...)
	return t, nil
}

// Start 启动服务
func (s *Server) Start() error {
	s.running = true
	errChan := make(chan error, 1)
	go func(ch chan error) {
		if err := s.server.ListenAndServe(); err != nil {
			ch <- err
		}
	}(errChan)

	select {
	case <-time.After(time.Millisecond * 500):
		return nil
	case err := <-errChan:
		s.running = false
		return err
	}
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() error {
	if s.server != nil && s.running {
		s.running = false
		defer s.metric.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			if err == xhttp.ErrServerClosed {
				return nil
			}
			return fmt.Errorf("关闭出现错误:%w", err)
		}
	}
	return nil
}

// GetAddress 获取当前服务地址
func (s *Server) GetAddress(h ...string) string {
	if len(h) > 0 && h[0] != "" {
		return fmt.Sprintf("%s://%s:%s", s.proto, h[0], s.port)
	}
	if s.host == "0.0.0.0" {
		return fmt.Sprintf("%s://%s:%s", s.proto, s.ip, s.port)
	}
	return fmt.Sprintf("%s://%s:%s", s.proto, s.host, s.port)
}

// GetStatus 获取当前服务器状态
func (s *Server) GetStatus() string {
	return types.DecodeString(s.running, true, "运行中", "停止")
}
