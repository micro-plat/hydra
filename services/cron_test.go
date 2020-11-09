/*
处理cron任务，包括任务注册、已注册任务获取、服务器订阅
*/

package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/security/md5"
)

//@todo 多线程测试订阅
func Test_cron_Add(t *testing.T) {

	tests := []struct {
		name         string
		cron         string
		service      string
		wantTasksLen int
		want         ICRON
	}{
		{name: "添加任务", cron: "task:task1", service: "service1", wantTasksLen: 1},
		{name: "再次添加任务", cron: "task:task2", service: "service2", wantTasksLen: 2},
		{name: "再次添加已存在任务", cron: "task:task2", service: "service2", wantTasksLen: 2},
	}
	c := newCron()
	c.Subscribe(func(t *task.Task) {
		fmt.Println("cron notify:", t)
	})

	//添加任务
	for _, tt := range tests {
		got := c.Add(tt.cron, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//获取任务
	tasks := c.GetTasks().Tasks
	assert.Equal(t, 2, len(tasks), "任务长度")
	keyMap := map[string]*task.Task{}
	for _, v := range tasks {
		keyMap[v.GetUNQ()] = v
	}

	//验证结果
	for _, tt := range tests {
		task := keyMap[md5.Encrypt(fmt.Sprintf("%s(%s)", tt.service, tt.cron))]
		assert.Equal(t, tt.cron, task.Cron, tt.name)
		assert.Equal(t, tt.service, task.Service, tt.name)
		assert.Equal(t, false, task.Disable, tt.name)
	}

	time.Sleep(time.Second)
}

func Test_cron_Subscribe(t *testing.T) {
	f := func(t *task.Task) {
		fmt.Println("cron notify:", t)
	}
	tests := []struct {
		name           string
		f              func(t *task.Task)
		wantEeventsLen int
	}{
		{name: "加入首个订阅者", f: f, wantEeventsLen: 1},
		{name: "再加入一个订阅者", f: f, wantEeventsLen: 2},
	}
	c := newCron()
	c.Add("cron1", "service1")
	c.Add("cron2", "service2")
	for _, tt := range tests {
		c.Subscribe(tt.f)
		l := len(c.subscribers)
		assert.Equal(t, tt.wantEeventsLen, l, tt.name)
		assert.Equal(t, 2, len(c.subscribers[l-1].taskChan), tt.name)
		i := 0
	LOOP:
		for {
			select {
			case task := <-c.subscribers[l-1].taskChan:
				assert.Equal(t, c.tasks.Tasks[i], task, tt.name)
				i++
			default:
				break LOOP
			}
		}
	}
}

func Test_cron_Remove(t *testing.T) {
	tests := []struct {
		name         string
		cron         string
		service      string
		wantTasksLen int
		want         ICRON
		wantDisable  bool
	}{
		{name: "移除存在任务", cron: "cron1", service: "service1", wantTasksLen: 1, wantDisable: true},
		{name: "移除不存在任务", cron: "cron3", service: "service3", wantTasksLen: 2, wantDisable: true},
	}
	c := newCron()
	c.Subscribe(func(t *task.Task) {
		fmt.Println("cron notify:", t)
	})

	//添加任务
	c.Add("cron1", "service1")
	c.Add("cron2", "service2")

	//移除任务
	for _, tt := range tests {
		got := c.Remove(tt.cron, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//获取任务
	tasks := c.GetTasks().Tasks
	assert.Equal(t, 3, len(tasks), "任务长度")
	keyMap := map[string]*task.Task{}
	for _, v := range tasks {
		keyMap[v.GetUNQ()] = v
	}

	//验证结果
	for _, tt := range tests {
		task := keyMap[md5.Encrypt(fmt.Sprintf("%s(%s)", tt.service, tt.cron))]
		assert.Equal(t, tt.cron, task.Cron, tt.name)
		assert.Equal(t, tt.service, task.Service, tt.name)
		assert.Equal(t, true, task.Disable, tt.name)
	}

	time.Sleep(time.Second * 2)
}
