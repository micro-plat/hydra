package mqc

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/mq"
)

type Processor struct {
	*dispatcher.Dispatcher
	mq.MQConsumer
	queues        []*conf.Queue
	handles       map[string]dispatcher.HandlerFunc
	isConsume     bool
	lock          sync.Mutex
	once          sync.Once
	done          bool
	address       string
	raw           string
	hasAddRouters bool
}

//NewProcessor 创建processor
func NewProcessor(address, raw string, queues []*conf.Queue) (p *Processor, err error) {
	p = &Processor{
		Dispatcher: dispatcher.New(),
		handles:    make(map[string]dispatcher.HandlerFunc),
		address:    address,
		raw:        raw,
		queues:     queues,
	}
	if p.MQConsumer, err = mq.NewMQConsumer(address, mq.WithRaw(raw)); err != nil {
		return
	}
	return p, nil
}

func (s *Processor) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.done = true
	s.isConsume = false
	if s.MQConsumer != nil {
		s.once.Do(func() {
			s.MQConsumer.Close()
			s.MQConsumer = nil
		})
	}
}

func (s *Processor) Consumes() (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.isConsume {
		return nil
	}
	if s.MQConsumer == nil {
		s.once = sync.Once{}
		s.MQConsumer, err = mq.NewMQConsumer(s.address, mq.WithRaw(s.raw), mq.WithQueueCount(len(s.queues)))
		if err != nil {
			return err
		}
	}
	for _, queue := range s.queues {
		err = s.Consume(queue)
		if err != nil {
			return err
		}
	}
	s.isConsume = len(s.queues) > 0
	return nil
}

//Consume 浪费指定的队列数据
func (s *Processor) Consume(r *conf.Queue) error {
	return s.MQConsumer.Consume(r.Queue, r.Concurrency, func(m mq.IMessage) {
		request := newMQRequest(r.Name, "GET", m.GetMessage())
		s.HandleRequest(request)
		request = nil
	})
}
