package services

import (
	"sync"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

type mqcSubscriber struct {
	f   func(t *queue.Queue)
	msg chan *queue.Queue
}

//MQC mqc消息
var MQC = newMQC()

type mqc struct {
	queues *queue.Queues
	events []*mqcSubscriber
	lock   sync.Mutex
	n      chan struct{}
}

func newMQC() *mqc {
	c := &mqc{
		queues: queue.NewEmptyQueues(),
		events: make([]*mqcSubscriber, 0, 0),
		n:      make(chan struct{}, 100),
	}
	go c.notify()
	return c
}

//Add 添加任务
func (c *mqc) Add(mqName string, service string, concurrency ...int) *mqc {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := queue.NewQueueByConcurrency(mqName, service, types.GetIntByIndex(concurrency, 0, 10))
	c.queues.Append(task)
	for _, s := range c.events {
		s.msg <- task
	}
	c.n <- struct{}{}
	return c
}

//Remove 移除任务
func (c *mqc) Remove(mqName string, service string) *mqc {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := queue.NewQueue(mqName, service)
	task.Disable = true
	c.queues.Append(task)
	for _, s := range c.events {
		s.msg <- task
	}
	c.n <- struct{}{}
	return c
}

//Subscribe 订阅任务
func (c *mqc) Subscribe(f func(t *queue.Queue)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	subscriber := &mqcSubscriber{
		f:   f,
		msg: make(chan *queue.Queue, len(c.queues.Queues)+100),
	}
	for _, t := range c.queues.Queues {
		subscriber.msg <- t
	}
	c.events = append(c.events, subscriber)

}

//notify 通知任务
func (c *mqc) notify() {
BREAK:
	for {
		select {
		case <-global.Current().ClosingNotify():
			break BREAK
		case <-c.n:
			c.lock.Lock()
			for _, e := range c.events {
			LOOP:
				for {
					select {
					case t := <-e.msg:
						e.f(t)
					default:
						break LOOP
					}
				}
			}
			c.lock.Unlock()
		}
	}
}
