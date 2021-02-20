package rpc

import (
	"fmt"
	xnet "net"
	"strings"
	"time"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/global/compatible"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/lib4go/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Server cron服务器
type Server struct {
	*Processor
	engine  *grpc.Server
	running bool
	addr    string
}

//NewServer 创建mqc服务器
//未使用压缩，由于传输数据默认限制为4M(已修改为20M)压缩后会影响系统并发能力
// grpc.RPCDecompressor(grpc.NewGZIPDecompressor())
func NewServer(addr string, routers []*router.Router, maxRecvSize, maxSendSize int) (t *Server, err error) {
	t = &Server{
		Processor: NewProcessor(routers...),
		engine: grpc.NewServer(
			grpc.MaxRecvMsgSize(maxRecvSize),
			grpc.MaxSendMsgSize(maxSendSize),
		),
	}

	if t.addr, err = GetAddress(addr); err != nil {
		return nil, err
	}

	return t, nil
}

//Start 启动rpc服务嚣
func (s *Server) Start() error {
	if s.running {
		return nil
	}
	s.running = true
	pb.RegisterRPCServer(s.engine, s.Processor)
	errChan := make(chan error, 1)
	go func(ch chan error) {
		lis, err := xnet.Listen("tcp", s.addr)
		if err != nil {
			ch <- err
			return
		}
		//debug模式，将注册的服务反射到server上，方便调试RPC接口
		if global.Def.IsDebug() {
			reflection.Register(s.engine)
		}
		if err := s.engine.Serve(lis); err != nil {
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

func GetAddress(addr string) (string, error) {
	host := "0.0.0.0"
	port := "8090"
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
		host = global.LocalIP()
	case "127.0.0.1", "localhost":
		break
	default:
		if net.GetLocalIPAddress(host) != host {
			return "", fmt.Errorf("%s地址不合法", addr)
		}
	}

	if !govalidator.IsPort(port) {
		return "", fmt.Errorf("%s端口不合法", addr)
	}
	if port == "80" {
		if err := compatible.CheckPrivileges(); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}

//Shutdown 关闭服务器
func (s *Server) Shutdown() {
	defer s.Processor.Close()
	if s.running {
		s.running = false
		s.engine.GracefulStop()
	}
}

//GetAddress 获取当前服务地址
func (s *Server) GetAddress() string {
	return fmt.Sprintf("tcp://%s", s.addr)
}
