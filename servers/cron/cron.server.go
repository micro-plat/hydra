package cron

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server/task"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
)

//CronServer cron服务器
type CronServer struct {
	*option
	conf   *conf.Metadata
	engine servers.IRegistryEngine
	*Processor
	running string
	addr    string
}

//NewCronServer 创建mqc服务器
func NewCronServer(name string, engine servers.IRegistryEngine, tasks []*task.Task, opts ...Option) (t *CronServer, err error) {
	t = &CronServer{
		engine: engine,
		conf:   conf.NewMetadata(name, "cron"),
	}
	t.option = &option{
		metric: middleware.NewMetric(t.conf),
		Logger: logger.GetSession(name, logger.CreateSession()),
	}
	for _, opt := range opts {
		opt(t.option)
	}

	if err = t.SetTasks(tasks); err != nil {
		return nil, err
	}

	t.ShowTrace(t.showTrace)
	return
}

//Start 启动cron服务嚣
func (s *CronServer) Start() error {
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

//Dynamic 动态注册或撤销cron任务
func (s *CronServer) Dynamic(engine servers.IRegistryEngine, c chan *task.Task) {
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			if s.running != servers.ST_RUNNING {
				return
			}
		case task := <-c:
			if task.Disable {
				s.Debugf("[取消定时任务(%s)]", task.GetUNQ())
				s.Processor.Remove(task.GetUNQ())
				continue
			}
			if err := task.Validate(); err != nil {
				s.Logger.Error(err)
				continue
			}
			if err := s.Processor.Add(task); err != nil {
				s.Logger.Error("添加cron到任务列表失败:", err)
			}
			s.Debugf("[注册定时任务(%s)(%s)]", task.Cron, task.Service)
		}
	}
}
