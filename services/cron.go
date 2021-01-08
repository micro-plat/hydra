/*
处理cron任务，包括任务注册、已注册任务获取、服务器订阅
*/

package services

import (
	"sync"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/global"
)

type subscriber struct {
	callback func(t *task.Task)
	taskChan chan *task.Task
}

//CRON cron消息
var CRON = newCron()

var _ ICRON = CRON

//ICRON CRON动态服务
type ICRON interface {
	Add(cron string, service string) ICRON
	Remove(cron string, service string) ICRON
}

type cron struct {
	dynamicTasks *task.Tasks
	staticTasks  *task.Tasks
	subscribers  []*subscriber
	lock         sync.Mutex
	staticLock   sync.Mutex
	signalChan   chan struct{}
}

func newCron() *cron {
	c := &cron{
		dynamicTasks: task.NewEmptyTasks(),
		staticTasks:  task.NewEmptyTasks(),
		subscribers:  make([]*subscriber, 0, 0),
		signalChan:   make(chan struct{}, 100),
	}
	go c.notify()
	return c
}

//GetTasks 获取任务列表
func (c *cron) GetTasks() *task.Tasks {
	return c.staticTasks
}

//Add 添加任务
func (c *cron) Static(cron string, service string) ICRON {
	c.staticLock.Lock()
	defer c.staticLock.Unlock()
	task := task.NewTask(cron, service)
	c.staticTasks.Append(task)
	return c
}

//Add 添加任务
func (c *cron) Add(cron string, service string) ICRON {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := task.NewTask(cron, service)
	_, notifyTasks := c.dynamicTasks.Append(task)
	for _, t := range notifyTasks {
		for _, s := range c.subscribers {
			s.taskChan <- t
		}
	}
	c.signalChan <- struct{}{}
	return c
}

//Remove 移除任务
func (c *cron) Remove(cron string, service string) ICRON {
	c.lock.Lock()
	defer c.lock.Unlock()
	task := task.NewTask(cron, service)
	task.Disable = true
	_, notifyTasks := c.dynamicTasks.Append(task)
	for _, t := range notifyTasks {
		for _, s := range c.subscribers {
			s.taskChan <- t
		}
	}
	c.signalChan <- struct{}{}
	return c
}

//Subscribe 订阅任务
func (c *cron) Subscribe(callback func(t *task.Task)) {

	subscriber := &subscriber{
		callback: callback,
		taskChan: make(chan *task.Task, 255),
	}
	c.lock.Lock()
	for _, t := range c.dynamicTasks.Tasks {
		subscriber.taskChan <- t
	}
	c.subscribers = append(c.subscribers, subscriber)

	c.lock.Unlock()
	c.sendNow(subscriber)
}

//notify 通知任务
func (c *cron) notify() {
	for {
		select {
		case <-global.Current().ClosingNotify():
			return
		case <-c.signalChan:
			c.lock.Lock()
			for _, e := range c.subscribers {
				c.sendNow(e)
			}
			c.lock.Unlock()
		}
	}
}

//sendNow 向订阅者推送消息
func (c *cron) sendNow(e *subscriber) {
SUBFOR:
	for {
		select {
		case t := <-e.taskChan:
			e.callback(t)
		default:
			break SUBFOR
		}
	}
}
