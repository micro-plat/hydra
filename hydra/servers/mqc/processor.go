package mqc

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/pkgs/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

const (
	unstarted = 1
	pause     = 2
	running   = 4
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	lock      sync.Mutex
	done      bool
	closeChan chan struct{}
	queues    cmap.ConcurrentMap
	startTime time.Time
	customer  mq.IMQC
	status    int
}

//NewProcessor 创建processor
func NewProcessor(addr string, raw []byte) (p *Processor, err error) {
	p = &Processor{
		status:    unstarted,
		closeChan: make(chan struct{}),
		startTime: time.Now(),
		queues:    cmap.New(4),
	}
	opt, err := mq.WithRaw(raw)
	if err != nil {
		return nil, fmt.Errorf("队列配置信息有误:%w", err)
	}
	p.customer, err = mq.NewMQC(addr, opt)
	if err != nil {
		return nil, fmt.Errorf("构建mqc服务失败 %s %w", addr, err)
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(MQC))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Use(middleware.Delay().DispFunc()) //
	return p, nil
}

//Start 所有任务
func (s *Processor) Start(wait ...bool) error {
	if err := s.customer.Connect(); err != nil {
		return err
	}
	if len(wait) > 0 && !wait[0] {
		_, err := s.Resume()
		return err
	}
	return nil
}

//Add 添加队列信息
func (s *Processor) Add(queues ...*queue.Queue) error {
	for _, queue := range queues {
		s.queues.SetIfAbsent(queue.Queue, queue)
	}
	return nil
}

//Remove 除移队列信息
func (s *Processor) Remove(queues ...*queue.Queue) error {
	for _, queue := range queues {
		s.customer.UnConsume(queue.Queue)
		s.queues.Remove(queue.Queue)
	}
	return nil
}

//Pause 暂停所有任务
func (s *Processor) Pause() (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != pause {
		s.status = pause
		items := s.queues.Items()
		for _, v := range items {
			queue := v.(*queue.Queue)
			s.customer.UnConsume(queue.Queue) //取消服务订阅
		}
		return true, nil
	}
	return false, nil
}

//Resume 恢复所有任务
func (s *Processor) Resume() (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.status != running {
		s.status = running
		items := s.queues.Items()
		for _, v := range items {
			queue := v.(*queue.Queue)
			if !s.Engine.Find(queue.Service) {
				s.Engine.Handle(DefMethod, queue.Service, middleware.ExecuteHandler(queue.Service).DispFunc(MQC))
			}
			if err := s.customer.Consume(queue.Queue, queue.Concurrency, s.handle(queue)); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

//Close 退出
func (s *Processor) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.done {
		s.done = true
		close(s.closeChan)
		s.queues.Clear()
		s.customer.Close()
	}
}

func (s *Processor) handle(queue *queue.Queue) func(mq.IMQCMessage) {
	return func(m mq.IMQCMessage) {
		req, err := NewRequest(queue, m)
		if err != nil {
			panic(err)
		}
		s.Engine.HandleRequest(req)
	}
}
