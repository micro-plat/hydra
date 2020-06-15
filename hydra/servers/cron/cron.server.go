package cron

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/lib4go/net"
)

//Server cron服务器
type Server struct {
	*Processor
	running bool
	addr    string
}

//NewServer 创建mqc服务器
func NewServer(name string, tasks ...*task.Task) (t *Server, err error) {
	t = &Server{Processor: NewProcessor()}
	if err := t.Processor.Add(tasks...); err != nil {
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
	case <-time.After(time.Millisecond * 500):
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
func (s *Server) Pause() error {
	if s.running {
		s.running = false
		s.Processor.Pause()
		time.Sleep(time.Second)
	}
	return nil

}

//Resume 恢复执行
func (s *Server) Resume() error {
	if !s.running {
		s.running = true
		s.Processor.Resume()
	}
	return nil
}

//GetAddress 获取当前服务地址
func (s *Server) GetAddress() string {
	return fmt.Sprintf("cron://%s", net.GetLocalIPAddress())
}
