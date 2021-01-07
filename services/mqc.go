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
	callback  func(t *queue.Queue)
	queueChan chan *queue.Queue
}

//MQC mqc消息
var MQC = newMQC()

type mqc struct {
	dynamicQueues *queue.Queues
	staticQueues  *queue.Queues
	subscribers   []*mqcSubscriber
	lock          sync.Mutex
	staticLock    sync.Mutex
	signalChan    chan struct{}
}

func newMQC() *mqc {
	c := &mqc{
		dynamicQueues: queue.NewEmptyQueues(),
		staticQueues:  queue.NewEmptyQueues(),
		subscribers:   make([]*mqcSubscriber, 0, 0),
		signalChan:    make(chan struct{}, 100),
	}
	go c.notify()
	return c
}

//Static 添加静态配置
func (c *mqc) Static(mqName string, service string, concurrency ...int) IMQC {
	return c.static(mqName, service, false, concurrency...)
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
	subscriber := &mqcSubscriber{
		callback:  f,
		queueChan: make(chan *queue.Queue, len(c.dynamicQueues.Queues)+100),
	}

	//锁定消息与订阅者，避免消息丢失
	c.lock.Lock()
	for _, t := range c.dynamicQueues.Queues {
		subscriber.queueChan <- t
	}
	c.subscribers = append(c.subscribers, subscriber)
	c.lock.Unlock()

	//立即向订阅者推送消息
	c.sendNow(subscriber)
}

//GetTasks 获取任务列表
func (c *mqc) GetQueues() *queue.Queues {
	return c.staticQueues
}

//notify 通知任务
func (c *mqc) notify() {
	for {
		select {
		case <-global.Current().ClosingNotify():
			return
		case <-c.signalChan:
			c.lock.Lock()
			for _, e := range c.subscribers {
				c.sendNow(e) //向订阅者推送消息
			}
			c.lock.Unlock()
		}
	}
}

//sendNow 向订阅者推送消息
func (c *mqc) sendNow(e *mqcSubscriber) {
SUBFOR:
	for {
		select {
		case t := <-e.queueChan:
			e.callback(t)
		default:
			break SUBFOR
		}
	}
}

func (c *mqc) add(mqName string, service string, disable bool, concurrency ...int) *mqc {
	f := func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		mqName = global.MQConf.GetQueueName(mqName)
		queue := queue.NewQueue(mqName, service, queue.WithConcurrency(types.GetIntByIndex(concurrency, 0, 10)))
		queue.Disable = disable
		_, notifyQueues := c.dynamicQueues.Append(queue)
		for _, q := range notifyQueues {
			for _, s := range c.subscribers {
				s.queueChan <- q
			}
		}
		c.signalChan <- struct{}{}
	}
	if !global.MQConf.NeedAddPrefix() {
		f()
		return c
	}
	global.OnReady(f)
	return c
}
func (c *mqc) static(mqName string, service string, disable bool, concurrency ...int) *mqc {
	f := func() {
		c.staticLock.Lock()
		defer c.staticLock.Unlock()
		mqName = global.MQConf.GetQueueName(mqName)
		queue := queue.NewQueue(mqName, service, queue.WithConcurrency(types.GetIntByIndex(concurrency, 0, 10)))
		queue.Disable = disable
		c.staticQueues.Append(queue)
	}
	if !global.MQConf.NeedAddPrefix() {
		f()
		return c
	}
	global.OnReady(f)
	return c
}
