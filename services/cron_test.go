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
)

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
	}
	c := newCron()
	c.Subscribe(func(t *task.Task) {
		fmt.Println("cron notify:", t)
	})
	for _, tt := range tests {
		got := c.Add(tt.cron, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
		l := len(c.tasks.Tasks)
		assert.Equal(t, tt.wantTasksLen, l, tt.name)
		assert.Equal(t, tt.cron, c.tasks.Tasks[l-1].Cron, tt.name)
		assert.Equal(t, tt.service, c.tasks.Tasks[l-1].Service, tt.name)
		assert.Equal(t, false, c.tasks.Tasks[l-1].Disable, tt.name)
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
	c.Add("cron", "service")
	for _, tt := range tests {
		c.Subscribe(tt.f)
		l := len(c.events)
		assert.Equal(t, tt.wantEeventsLen, l, tt.name)
		//assert.Equal(t, tt.f, c.events[l-1].f, tt.name)
		select {
		case task := <-c.events[l-1].msg:
			assert.Equal(t, c.tasks.Tasks[0], task, tt.name)
		default:
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
		{name: "移除任务", cron: "task:task1", service: "service1", wantTasksLen: 1, wantDisable: true},
		{name: "再次移除任务", cron: "task:task2", service: "service2", wantTasksLen: 2, wantDisable: true},
	}
	c := newCron()
	for _, tt := range tests {
		c.Remove(tt.cron, tt.service)
		assert.Equal(t, tt.wantTasksLen, len(c.tasks.Tasks), tt.name)
		for _, v := range c.tasks.Tasks {
			if v.Cron == tt.cron {
				assert.Equal(t, tt.wantDisable, v.Disable, tt.name)
			}
		}
	}
}
