package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/assert"
)

func Test_mqc_Subscribe(t *testing.T) {
	//添加队列
	c := newMQC()
	global.MQConf.PlatNameAsPrefix(false)
	c.Add("mq1", "service1")
	c.Add("mq2", "service2")
	c.Add("mq3", "service3")
	time.Sleep(time.Second * 1)

	coroutine := 100
	for i := 0; i < coroutine; i++ {
		go func() {
			c.Subscribe(func(t *queue.Queue) {})
		}()
	}
	time.Sleep(time.Second)
	assert.Equal(t, coroutine, len(c.subscribers), "订阅者长度")
	for _, v := range c.subscribers {
		assert.Equal(t, 0, len(v.queueChan), "订阅者队列长度")
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
		{name: "1.1 添加队列", mqName: "mqc1", service: "service1", wantQueyesLen: 1, concurrency: 1},
		{name: "1.2 添加不存在队列", mqName: "mqc2", service: "service2", wantQueyesLen: 2, concurrency: 2},
		{name: "1.3 添加已存在的队列", mqName: "mqc2", service: "service2", wantQueyesLen: 2, concurrency: 2},
	}
	c := newMQC()
	//添加队列
	global.MQConf.PlatNameAsPrefix(false)

	for _, tt := range tests {
		got := c.Add(tt.mqName, tt.service, tt.concurrency)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//验证结果
	keyMap := map[string]*queue.Queue{}
	for _, v := range c.dynamicQueues.Queues {
		keyMap[v.Queue] = v
	}

	for _, tt := range tests {
		queue := keyMap[tt.mqName]
		assert.Equal(t, tt.mqName, queue.Queue, tt.name)
		assert.Equal(t, tt.service, queue.Service, tt.name)
		assert.Equal(t, tt.concurrency, queue.Concurrency, tt.name)
		assert.Equal(t, false, queue.Disable, tt.name)
	}
}

func Test_mqc_Add_WithMultithread(t *testing.T) {

	c := newMQC()
	var lock1, lock2 sync.Mutex
	subscriber1 := 0
	subscriber2 := 0
	c.Subscribe(func(t *queue.Queue) {
		lock1.Lock()
		defer lock1.Unlock()
		subscriber1++
	})
	c.Subscribe(func(t *queue.Queue) {
		lock2.Lock()
		defer lock2.Unlock()
		subscriber2++
	})

	//添加队列
	global.MQConf.PlatNameAsPrefix(false)
	//多线程添加
	coroutine := 100
	for i := 0; i < coroutine; i++ {
		go func() {
			got := c.Add("mqc"+fmt.Sprint(time.Now().UnixNano()), "service")
			assert.Equal(t, true, got != nil, "添加队列")
		}()
	}
	time.Sleep(time.Second * 2)

	//获取队列
	queues := c.dynamicQueues.Queues
	assert.Equal(t, coroutine, len(queues), "队列长度")

	time.Sleep(time.Second)
	assert.Equal(t, coroutine, subscriber1, "订阅者接收队列数量1")
	assert.Equal(t, coroutine, subscriber2, "订阅者接收队列数量2")
}

func Test_mqc_Remove(t *testing.T) {
	c := newMQC()

	//添加队列
	addQueuesLen := 50
	global.MQConf.PlatNameAsPrefix(false)
	for i := 0; i < addQueuesLen; i++ {
		got := c.Add("queue"+fmt.Sprint(i), "service")
		assert.Equal(t, true, got != nil, "添加队列")
	}

	//订阅
	var lock1, lock2 sync.Mutex
	subscriber1 := 0
	subscriber2 := 0
	c.Subscribe(func(t *queue.Queue) {
		lock1.Lock()
		defer lock1.Unlock()
		subscriber1++
	})
	c.Subscribe(func(t *queue.Queue) {
		lock2.Lock()
		defer lock2.Unlock()
		subscriber2++
	})

	var lock3 sync.Mutex
	all := map[string]bool{}
	coroutine := 100
	noNotify := 0
	for i := 0; i < coroutine; i++ {
		go func() {
			//移除存在的队列或者不存在的队列
			s := rand.Intn(150)
			got := c.Remove("queue"+fmt.Sprint(s), "service")
			assert.Equal(t, true, got != nil, "移除队列")
			lock3.Lock()
			defer lock3.Unlock()
			if _, ok := all[fmt.Sprint(s)]; ok {
				noNotify++
			}
			all[fmt.Sprint(s)] = true
		}()
	}

	time.Sleep(time.Second)

	nonExist := 0
	for k := range all {
		index, _ := strconv.Atoi(k)
		if index > addQueuesLen-1 {
			nonExist++
		}
	}
	//获取队列
	queues := c.dynamicQueues.Queues
	assert.Equal(t, nonExist+addQueuesLen, len(queues), "队列长度")

	for _, v := range queues {
		index := strings.TrimPrefix(v.Queue, "queue")
		if _, ok := all[index]; ok {
			assert.Equal(t, true, v.Disable, "队列状态1")
			assert.Equal(t, 0, v.Concurrency, "队列并发1")
			continue
		}
		assert.Equal(t, false, v.Disable, "队列状态2")
	}
	//订阅者长度
	assert.Equal(t, coroutine+addQueuesLen-noNotify, subscriber1, "订阅者接收队列通知数量1")
	assert.Equal(t, coroutine+addQueuesLen-noNotify, subscriber2, "订阅者接收队列通知数量2")
}
