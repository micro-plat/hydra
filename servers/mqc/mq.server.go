package mqc

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
)

//MqcServer mqc服务器
type MqcServer struct {
	*option
	conf *conf.MetadataConf
	*Processor
	running string
	addr    string
}

//NewMqcServer 创建mqc服务器
func NewMqcServer(name string, proto string, config string, queues []*conf.Queue, opts ...Option) (t *MqcServer, err error) {
	t = &MqcServer{conf: &conf.MetadataConf{Name: name, Type: "mqc"}}
	t.option = &option{metric: middleware.NewMetric(t.conf)}
	for _, opt := range opts {
		opt(t.option)
	}
	t.conf.Name = fmt.Sprintf("%s.%s.%s", t.platName, t.systemName, t.clusterName)
	if t.Logger == nil {
		t.Logger = logger.GetSession(name, logger.CreateSession())
	}
	if queues != nil && proto != "" && len(queues) > 0 {
		err = t.SetQueues(proto, config, queues)
	}
	t.SetTrace(t.showTrace)
	return
}

// Run the http server
func (s *MqcServer) Run() error {
	if s.running == servers.ST_RUNNING {
		return nil
	}
	s.running = servers.ST_RUNNING
	errChan := make(chan error, 1)
	if err := s.Processor.Consumes(); err != nil {
		return err
	}
	go func(ch chan error) {
		if err := s.Processor.Connect(); err != nil {
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
func (s *MqcServer) Shutdown(timeout time.Duration) {
	if s.Processor != nil {
		s.running = servers.ST_STOP
		s.Processor.Close()
		time.Sleep(time.Second)
	}
}

//Pause 暂停服务器
func (s *MqcServer) Pause(timeout time.Duration) {
	if s.Processor != nil {
		s.running = servers.ST_PAUSE
		s.Processor.Close()
		time.Sleep(time.Second)
	}
}

//GetAddress 获取当前服务地址
func (s *MqcServer) GetAddress() string {
	return fmt.Sprintf("mqc://%s", net.GetLocalIPAddress())
}

//GetStatus 获取当前服务器状态
func (s *MqcServer) GetStatus() string {
	return s.running
}
