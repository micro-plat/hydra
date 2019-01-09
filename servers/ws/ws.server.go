package ws

import (
	"context"
	"errors"
	"fmt"
	x "net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
)

//WSServer WSServer服务器
type WSServer struct {
	*option
	conf    *conf.MetadataConf
	engine  *x.Server
	running string
	proto   string
	host    string
	port    string
}

//NewWSServerServer 创建WSServer服务器
func NewWSServerServer(name string, addr string, routers []*conf.Router, opts ...Option) (t *WSServer, err error) {
	t = &WSServer{conf: &conf.MetadataConf{
		Name: name,
		Type: "ws",
	}}
	t.option = &option{
		metric:            middleware.NewMetric(t.conf),
		readHeaderTimeout: 6,
		readTimeout:       6,
		writeTimeout:      6}
	for _, opt := range opts {
		opt(t.option)
	}
	t.conf.Name = fmt.Sprintf("%s.%s.%s", t.platName, t.systemName, t.clusterName)
	if t.Logger == nil {
		t.Logger = logger.GetSession(name, logger.CreateSession())
	}
	naddr, err := t.getAddress(addr)
	if err != nil {
		return nil, err
	}
	t.engine = &x.Server{
		Addr:              naddr,
		ReadHeaderTimeout: time.Second * time.Duration(t.option.readHeaderTimeout),
		ReadTimeout:       time.Second * time.Duration(t.option.readTimeout),
		WriteTimeout:      time.Second * time.Duration(t.option.writeTimeout),
		MaxHeaderBytes:    1 << 20,
	}
	if routers != nil {
		t.engine.Handler, err = t.getHandler(routers)
	}
	t.SetTrace(t.showTrace)
	return
}

// Run the http server
func (s *WSServer) Run() error {
	s.proto = "ws"
	s.running = servers.ST_RUNNING
	errChan := make(chan error, 1)
	go func(ch chan error) {
		if err := s.engine.ListenAndServe(); err != nil {
			ch <- err
		}
	}(errChan)
	select {
	case <-time.After(time.Millisecond * 500):
		return nil
	case err := <-errChan:
		s.running = servers.ST_STOP
		return err
	}
}

//Shutdown 关闭服务器
func (s *WSServer) Shutdown(timeout time.Duration) {
	if s.engine != nil {
		s.metric.Stop()
		s.running = servers.ST_STOP
		ctx, cannel := context.WithTimeout(context.Background(), timeout)
		defer cannel()
		if err := s.engine.Shutdown(ctx); err != nil {
			if err == x.ErrServerClosed {
				s.Infof("%s:已关闭", s.conf.Name)
				return
			}
			s.Errorf("关闭出现错误:%v", err)
		}
	}
}

//GetAddress 获取当前服务地址
func (s *WSServer) GetAddress() string {
	return fmt.Sprintf("%s://%s:%s", s.proto, s.host, s.port)
}

//GetStatus 获取当前服务器状态
func (s *WSServer) GetStatus() string {
	return s.running
}

func (s *WSServer) getAddress(addr string) (string, error) {
	host := "0.0.0.0"
	port := "8081"
	args := strings.Split(addr, ":")
	l := len(args)
	if addr == "" {
		l = 0
	}
	switch l {
	case 0:
	case 1:
		if govalidator.IsPort(args[0]) {
			port = args[0]
		} else {
			host = args[0]
		}
	case 2:
		host = args[0]
		port = args[1]
	default:
		return "", fmt.Errorf("%s地址不合法", addr)
	}
	switch host {
	case "0.0.0.0", "":
		s.host = net.GetLocalIPAddress()
	case "127.0.0.1", "localhost":
		s.host = host
	default:
		if net.GetLocalIPAddress(host) != host {
			return "", fmt.Errorf("%s地址不合法", addr)
		}
		s.host = host
	}

	if !govalidator.IsPort(port) {
		return "", fmt.Errorf("%s端口不合法", addr)
	}
	if port == "80" {
		if err := checkPrivileges(); err != nil {
			return "", err
		}
	}
	s.port = port
	return fmt.Sprintf("%s:%s", host, s.port), nil
}
func checkPrivileges() error {
	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return nil
			}
			return ErrRootPrivileges
		}
	}
	return ErrUnsupportedSystem
}

var ErrUnsupportedSystem = errors.New("Unsupported system")
var ErrRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")
