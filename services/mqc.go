package services

import (
	"sync"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

//IMQC MQC动态服务管理
type IMQC interface {
	Add(mqName string, service string, concurrency ...int) IMQC
	Remove(mqName string, service string) IMQC
}
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

//Add 添加任务,队列名称需要接平台名称时将会自动转为异步注册(系统准备好后注册)
func (c *mqc) Add(mqName string, service string, concurrency ...int) IMQC {
	return c.add(mqName, service, false, concurrency...)
}

//Remove 移除任务
func (c *mqc) Remove(mqName string, service string) IMQC {
	return c.add(mqName, service, true, 0)
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

//GetTasks 获取任务列表
func (c *mqc) GetQueues() *queue.Queues {
	return c.queues
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

func (c *mqc) add(mqName string, service string, disable bool, concurrency ...int) *mqc {
	f := func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		mqName = global.MQConf.GetQueueName(mqName)
		task := queue.NewQueueByConcurrency(mqName, service, types.GetIntByIndex(concurrency, 0, 10))
		task.Disable = disable
		c.queues.Append(task)
		for _, s := range c.events {
			s.msg <- task
		}
		c.n <- struct{}{}
	}
	if !global.MQConf.NeedAddPrefix() {
		f()
		return c
	}
	global.OnReady(f)
	return c
}
