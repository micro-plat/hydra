package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_mqc_Subscribe(t *testing.T) {
	f := func(t *queue.Queue) {
		fmt.Println("queue notify:", t)
	}
	tests := []struct {
		name           string
		f              func(t *queue.Queue)
		wantEeventsLen int
	}{
		{name: "加入首个订阅者", f: f, wantEeventsLen: 1},
		{name: "再加入一个订阅者", f: f, wantEeventsLen: 2},
	}
	c := newMQC()
	c.Add("mqc", "service")
	for _, tt := range tests {
		c.Subscribe(tt.f)
		l := len(c.events)
		assert.Equal(t, tt.wantEeventsLen, l, tt.name)
		//assert.Equal(t, tt.f, c.events[l-1].f, tt.name)
		select {
		case queue := <-c.events[l-1].msg:
			assert.Equal(t, c.queues.Queues[0], queue, tt.name)
		default:
		}
	}
}

func Test_mqc_Add(t *testing.T) {

	tests := []struct {
		name          string
		mqName        string
		service       string
		concurrency   int
		wantQueyesLen int
		want          IMQC
	}{
		{name: "添加队列", mqName: "mqc1", service: "service1", wantQueyesLen: 1, concurrency: 1},
		{name: "再次添加队列", mqName: "mqc2", service: "service2", wantQueyesLen: 2, concurrency: 2},
	}
	c := newMQC()
	c.Subscribe(func(t *queue.Queue) {
		fmt.Println("mqc notify:", t)
	})
	global.MQConf.PlatNameAsPrefix(false)
	for _, tt := range tests {
		got := c.Add(tt.mqName, tt.service, tt.concurrency)
		assert.Equal(t, true, got != nil, tt.name)
		l := len(c.GetQueues().Queues)
		assert.Equal(t, tt.wantQueyesLen, l, tt.name)
		assert.Equal(t, tt.mqName, c.queues.Queues[l-1].Queue, tt.name)
		assert.Equal(t, tt.service, c.queues.Queues[l-1].Service, tt.name)
		assert.Equal(t, tt.concurrency, c.queues.Queues[l-1].Concurrency, tt.name)
		assert.Equal(t, false, c.queues.Queues[l-1].Disable, tt.name)
	}
	time.Sleep(time.Second)
}

func Test_mqc_Remove(t *testing.T) {

	tests := []struct {
		name          string
		mqName        string
		service       string
		wantQueyesLen int
		want          IMQC
	}{
		{name: "添加移除队列", mqName: "mqc1", service: "service1", wantQueyesLen: 1},
		{name: "再次添加移除队列", mqName: "mqc2", service: "service2", wantQueyesLen: 2},
	}
	c := newMQC()
	c.Subscribe(func(t *queue.Queue) {
		fmt.Println("mqc notify:", t)
	})
	global.MQConf.PlatNameAsPrefix(false)
	for _, tt := range tests {
		got := c.Remove(tt.mqName, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
		l := len(c.GetQueues().Queues)
		assert.Equal(t, tt.wantQueyesLen, l, tt.name)
		assert.Equal(t, tt.mqName, c.queues.Queues[l-1].Queue, tt.name)
		assert.Equal(t, tt.service, c.queues.Queues[l-1].Service, tt.name)
		assert.Equal(t, 0, c.queues.Queues[l-1].Concurrency, tt.name)
		assert.Equal(t, true, c.queues.Queues[l-1].Disable, tt.name)
	}
	time.Sleep(time.Second)
}
