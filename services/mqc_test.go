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
	global.MQConf.PlatNameAsPrefix(false)
	c := newMQC()
	c.Add("mqc1", "service1")
	c.Add("mqc2", "service2")
	for _, tt := range tests {
		c.Subscribe(tt.f)
		l := len(c.subscribers)
		assert.Equal(t, tt.wantEeventsLen, l, tt.name)
		assert.Equal(t, 2, len(c.subscribers[l-1].queueChan), tt.name)
		i := 0
	LOOP:
		for {
			select {
			case queue := <-c.subscribers[l-1].queueChan:
				assert.Equal(t, c.queues.Queues[i], queue, tt.name)
				i++
			default:
				break LOOP
			}
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
		{name: "添加已存在的队列", mqName: "mqc2", service: "service2", wantQueyesLen: 2, concurrency: 2},
	}
	c := newMQC()
	c.Subscribe(func(t *queue.Queue) {
		fmt.Println("mqc notify:", t)
	})
	global.MQConf.PlatNameAsPrefix(false)

	//添加队列
	for _, tt := range tests {
		got := c.Add(tt.mqName, tt.service, tt.concurrency)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//获取队列
	queues := c.GetQueues().Queues
	assert.Equal(t, 2, len(queues), "队列长度")
	keyMap := map[string]*queue.Queue{}
	for _, v := range queues {
		keyMap[v.Queue] = v
	}

	//验证结果
	for _, tt := range tests {
		queue := keyMap[tt.mqName]
		assert.Equal(t, tt.mqName, queue.Queue, tt.name)
		assert.Equal(t, tt.service, queue.Service, tt.name)
		assert.Equal(t, tt.concurrency, queue.Concurrency, tt.name)
		assert.Equal(t, false, queue.Disable, tt.name)
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
		{name: "移除存在队列", mqName: "mqc1", service: "service1", wantQueyesLen: 1},
		{name: "移除不存在队列", mqName: "mqc3", service: "service3", wantQueyesLen: 2},
	}
	c := newMQC()
	c.Subscribe(func(t *queue.Queue) {
		fmt.Println("mqc notify:", t)
	})
	global.MQConf.PlatNameAsPrefix(false)

	//添加队列
	c.Add("mqc1", "service1")
	c.Add("mqc2", "service2")

	//移除队列
	for _, tt := range tests {
		got := c.Remove(tt.mqName, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//获取队列
	queues := c.GetQueues().Queues
	assert.Equal(t, 3, len(queues), "队列长度")
	keyMap := map[string]*queue.Queue{}
	for _, v := range queues {
		keyMap[v.Queue] = v
	}

	//验证结果
	for _, tt := range tests {
		q := keyMap[tt.mqName]
		assert.Equal(t, tt.mqName, q.Queue, tt.name)
		assert.Equal(t, tt.service, q.Service, tt.name)
		assert.Equal(t, 0, q.Concurrency, tt.name)
		assert.Equal(t, true, q.Disable, tt.name)
	}

	time.Sleep(time.Second)
}
