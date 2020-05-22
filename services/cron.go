package services

import (
	"sync"

	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/registry/conf/server/task"
)

type subscriber struct {
	f   func(t *task.Task)
	msg chan *task.Task
}

//CRON cron消息
var CRON = newCron()

type cron struct {
	tasks  *task.Tasks
	events []*subscriber
	lock   sync.Mutex
	n      chan struct{}
}

func newCron() *cron {
	c := &cron{
		tasks:  task.NewEmptyTasks(),
		events: make([]*subscriber, 0, 0),
		n:      make(chan struct{}, 100),
	}
	go c.notify()
	return c
}

//Add 添加任务
func (c *cron) Add(cron string, service string) *cron {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := task.NewTask(cron, service)
	c.tasks.Append(task)
	for _, s := range c.events {
		s.msg <- task
	}
	c.n <- struct{}{}
	return c
}

//Remove 移除任务
func (c *cron) Remove(cron string, service string) *cron {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := task.NewTask(cron, service)
	task.Disable = true
	c.tasks.Append(task)
	for _, s := range c.events {
		s.msg <- task
	}
	c.n <- struct{}{}
	return c
}

//Subscribe 订阅任务
func (c *cron) Subscribe(f func(t *task.Task)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	subscriber := &subscriber{
		f:   f,
		msg: make(chan *task.Task, len(c.tasks.Tasks)+100),
	}
	for _, t := range c.tasks.Tasks {
		subscriber.msg <- t
	}
	c.events = append(c.events, subscriber)
}

//notify 通知任务
func (c *cron) notify() {
BREAK:
	for {
		select {
		case <-application.Current().ClosingNotify():
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
