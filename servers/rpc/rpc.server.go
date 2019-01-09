package rpc

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/rpc/pb"
	xnet "github.com/micro-plat/lib4go/net"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
	"google.golang.org/grpc"
)

//RpcServer rpc服务器
type RpcServer struct {
	*option
	conf   *conf.MetadataConf
	engine *grpc.Server
	*Processor
	running string
	proto   string
	port    string
	addr    string
	host    string
}

//NewRpcServer 创建rpc服务器
func NewRpcServer(name string, address string, routers []*conf.Router, opts ...Option) (t *RpcServer, err error) {
	t = &RpcServer{conf: &conf.MetadataConf{
		Name: name,
		Type: "rpc",
	}}
	if t.addr, err = t.getAddress(address); err != nil {
		return nil, err
	}
	t.option = &option{metric: middleware.NewMetric(t.conf)}
	for _, opt := range opts {
		opt(t.option)
	}
	t.conf.Name = fmt.Sprintf("%s.%s.%s", t.platName, t.systemName, t.clusterName)
	if t.Logger == nil {
		t.Logger = logger.GetSession(name, logger.CreateSession())
	}
	t.engine = grpc.NewServer()
	if routers != nil {
		t.Processor, err = t.getProcessor(routers)
		if err != nil {
			return
		}
	}
	t.SetTrace(t.showTrace)
	return
}

// Run the http server
func (s *RpcServer) Run() error {
	pb.RegisterRPCServer(s.engine, s.Processor)
	s.proto = "tcp"
	s.running = servers.ST_RUNNING
	errChan := make(chan error, 1)
	go func(ch chan error) {
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			ch <- err
			return
		}
		if err := s.engine.Serve(lis); err != nil {
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
func (s *RpcServer) Shutdown(timeout time.Duration) {
	if s.engine != nil {
		s.running = servers.ST_STOP
		s.engine.GracefulStop()
		time.Sleep(time.Second)

	}
}

//GetAddress 获取当前服务地址
func (s *RpcServer) GetAddress() string {
	return fmt.Sprintf("%s://%s:%s", s.proto, s.host, s.port)
}

//GetStatus 获取当前服务器状态
func (s *RpcServer) GetStatus() string {
	return s.running
}

func (s *RpcServer) getAddress(addr string) (string, error) {
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
		s.host = xnet.GetLocalIPAddress()
	case "127.0.0.1", "localhost":
		s.host = host
	default:
		if xnet.GetLocalIPAddress(host) != host {
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
