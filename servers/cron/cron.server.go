package cron

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
)

//CronServer cron服务器
type CronServer struct {
	*option
	conf *conf.MetadataConf
	*Processor
	running string
	addr    string
}

//NewCronServer 创建mqc服务器
func NewCronServer(name string, config string, tasks []*conf.Task, opts ...Option) (t *CronServer, err error) {
	t = &CronServer{conf: &conf.MetadataConf{Name: name, Type: "cron"}}
	t.option = &option{metric: middleware.NewMetric(t.conf)}
	for _, opt := range opts {
		opt(t.option)
	}
	if t.Logger == nil {
		t.Logger = logger.GetSession(name, logger.CreateSession())
	}
	if tasks != nil && len(tasks) > 0 {
		err = t.SetTasks(config, tasks)
	}
	t.SetTrace(t.showTrace)
	return
}
func (s *CronServer) Start() error {
	return s.Run()
}

func (s *CronServer) Run() error {
	if s.running == servers.ST_RUNNING {
		return nil
	}
	s.running = servers.ST_RUNNING
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
		s.running = servers.ST_STOP
		return err
	}
}

//Shutdown 关闭服务器
func (s *CronServer) Shutdown(time.Duration) {
	if s.Processor != nil {
		s.running = servers.ST_STOP
		s.Processor.Close()
	}
}

//Pause 暂停服务器
func (s *CronServer) Pause() {
	if s.Processor != nil {
		s.running = servers.ST_PAUSE
		s.Processor.Pause()
		time.Sleep(time.Second)
	}
}

//Resume 恢复执行
func (s *CronServer) Resume() error {
	if s.Processor != nil {
		s.running = servers.ST_RUNNING
		s.Processor.Resume()
	}
	return nil
}

//GetAddress 获取当前服务地址
func (s *CronServer) GetAddress() string {
	return fmt.Sprintf("cron://%s", net.GetLocalIPAddress())
}

//GetStatus 获取当前服务器状态
func (s *CronServer) GetStatus() string {
	return s.running
}
