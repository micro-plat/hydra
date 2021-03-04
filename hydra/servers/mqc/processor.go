package mqc

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
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
	//*dispatcher.Engine
	lock          sync.Mutex
	done          bool
	closeChan     chan struct{}
	queues        cmap.ConcurrentMap
	metric        *middleware.Metric
	startTime     time.Time
	customer      mq.IMQC
	status        int
	adapterEngine *adapter.Engine
}

//NewProcessor 创建processor
func NewProcessor(proto string, confRaw string) (p *Processor, err error) {
	p = &Processor{
		status:    unstarted,
		closeChan: make(chan struct{}),
		startTime: time.Now(),
		queues:    cmap.New(4),
		metric:    middleware.NewMetric(),
	}

	p.customer, err = mq.NewMQC(proto, confRaw)
	if err != nil {
		return nil, fmt.Errorf("构建mqc服务失败(proto:%s,raw:%s) %v", proto, confRaw, err)
	}
	p.adapterEngine = adapter.New(adapter.NewEngineWrapperDisp(dispatcher.New(), MQC))
	//p.Engine = p.adapterEngine.DispEngine()

	p.adapterEngine.Use(middleware.Recovery())
	p.adapterEngine.Use(middleware.Logging())
	p.adapterEngine.Use(middleware.Recovery())
	p.adapterEngine.Use(p.metric.Handle())
	p.adapterEngine.Use(middleware.Trace()) //跟踪信息
	p.adapterEngine.Use(middlewares...)

	return p, nil
}

//Done Done
func (s *Processor) Done() bool {
	return s.done
}

//QueueItems QueueItems
func (s *Processor) QueueItems() map[string]interface{} {
	return s.queues.Items()
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
		if ok, _ := s.queues.SetIfAbsent(queue.Queue, queue); ok && s.status == running {
			if err := s.consume(queue); err != nil {
				return err
			}
		}
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
			if err := s.consume(queue); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}
func (s *Processor) consume(queue *queue.Queue) error {
	if !s.adapterEngine.Find(queue.Service) {
		s.adapterEngine.Handle(queue)
	}
	if err := s.customer.Consume(queue.Queue, queue.Concurrency, s.handle(queue)); err != nil {
		return err
	}
	return nil
}

//Close 退出
func (s *Processor) Close() {
	defer s.metric.Stop()
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
		s.adapterEngine.HandleRequest(req)
	}
}
