package mqc

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/router"
)

//Server cron服务器
type Server struct {
	*Processor
	running bool
	addr    string
}

//NewServer 创建mqc服务器
func NewServer(proto string, raw []byte, queues []*queue.Queue, routers ...*router.Router) (t *Server, err error) {
	p, err := NewProcessor(proto, string(raw), routers...)
	if err != nil {
		return nil, err
	}
	t = &Server{Processor: p}
	if err := t.Processor.Add(queues...); err != nil {
		return nil, err
	}
	return
}

//Start 启动cron服务嚣
func (s *Server) Start() error {
	if s.running {
		return nil
	}
	s.running = true
	errChan := make(chan error, 1)
	go func(ch chan error) {
		if err := s.Processor.Start(); err != nil {
			ch <- err
		}
	}(errChan)
	select {
	case <-time.After(time.Millisecond * 200):
		return nil
	case err := <-errChan:
		s.running = false
		return err
	}
}

//Shutdown 关闭服务器
func (s *Server) Shutdown() {
	s.running = false
	s.Processor.Close()
}

//Pause 暂停服务器
func (s *Server) Pause() (bool, error) {
	ok, err := s.Processor.Pause()
	if ok {
		s.running = false
	}
	return ok, err
}

//Resume 恢复执行
func (s *Server) Resume() (bool, error) {
	ok, err := s.Processor.Resume()
	if ok {
		s.running = true
	}
	return ok, err

}

//GetAddress 获取当前服务地址
func (s *Server) GetAddress() string {
	if s.addr == "" {
		s.addr = fmt.Sprintf("mqc://%s", global.LocalIP())
	}
	return s.addr
}
