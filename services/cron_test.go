/*
处理cron任务，包括任务注册、已注册任务获取、服务器订阅
*/

package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/security/md5"
)

func Test_cron_Add(t *testing.T) {

	tests := []struct {
		name         string
		cron         string
		service      string
		wantTasksLen int
		want         ICRON
	}{
		{name: "1.1 添加任务", cron: "task:task1", service: "service1", wantTasksLen: 1},
		{name: "1.2 添加不存在任务", cron: "task:task2", service: "service2", wantTasksLen: 2},
		{name: "1.3 添加已存在任务", cron: "task:task2", service: "service2", wantTasksLen: 2},
	}
	c := newCron()

	//添加任务
	for _, tt := range tests {
		got := c.Add(tt.cron, tt.service)
		assert.Equal(t, true, got != nil, tt.name)
	}

	//获取任务
	tasks := c.dynamicTasks.Tasks
	assert.Equal(t, len(tests)-1, len(tasks), "任务长度")
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

}

func Test_cron_Add_WithMultithread(t *testing.T) {

	c := newCron()
	var lock1, lock2 sync.Mutex
	subscriber1 := 0
	subscriber2 := 0
	c.Subscribe(func(t *task.Task) {
		lock1.Lock()
		defer lock1.Unlock()
		subscriber1++
	})
	c.Subscribe(func(t *task.Task) {
		lock2.Lock()
		defer lock2.Unlock()
		subscriber2++
	})

	coroutine := 100
	for i := 0; i < coroutine; i++ {
		go func() {
			got := c.Add("task"+fmt.Sprint(time.Now().UnixNano()), "service")
			assert.Equal(t, true, got != nil, "添加任务")
		}()
	}
	time.Sleep(time.Second * 2)
	//获取任务
	tasks := c.dynamicTasks.Tasks
	assert.Equal(t, coroutine, len(tasks), "任务长度")
	for _, v := range tasks {
		assert.Equal(t, false, v.Disable, "任务状态")
	}

	time.Sleep(time.Second)
	assert.Equal(t, coroutine, subscriber1, "订阅者接收任务数量1")
	assert.Equal(t, coroutine, subscriber2, "订阅者接收任务数量2")
}

func Test_cron_Subscribe(t *testing.T) {

	//添加任务
	c := newCron()
	c.Add("cron1", "service1")
	c.Add("cron2", "service2")
	c.Add("cron3", "service3")
	time.Sleep(time.Second * 1)

	coroutine := 100
	for i := 0; i < coroutine; i++ {
		go func() {
			c.Subscribe(func(t *task.Task) {})
		}()
	}

	time.Sleep(time.Second * 2)
	assert.Equal(t, coroutine, len(c.subscribers), "订阅者长度")
	for _, v := range c.subscribers {
		assert.Equal(t, 0, len(v.taskChan), "订阅者任务长度")
	}

}

func Test_cron_Remove(t *testing.T) {
	c := newCron()
	//添加任务
	addTasksLen := 50
	for i := 0; i < addTasksLen; i++ {
		got := c.Add("task"+fmt.Sprint(i), "service")
		assert.Equal(t, true, got != nil, "添加任务")
	}

	//订阅
	var lock1, lock2 sync.Mutex
	subscriber1 := 0
	subscriber2 := 0
	c.Subscribe(func(t *task.Task) {
		lock1.Lock()
		defer lock1.Unlock()
		subscriber1++
	})
	c.Subscribe(func(t *task.Task) {
		lock2.Lock()
		defer lock2.Unlock()
		subscriber2++
	})

	var lock3 sync.Mutex
	all := map[string]bool{}
	noNotify := 0
	coroutine := 100
	for i := 0; i < coroutine; i++ {
		go func() {
			//移除存在的任务或者不存在的任务
			s := rand.Intn(150)
			got := c.Remove("task"+fmt.Sprint(s), "service")
			assert.Equal(t, true, got != nil, "移除任务")
			lock3.Lock()
			defer lock3.Unlock()
			if _, ok := all[fmt.Sprint(s)]; ok {
				noNotify++
			}
			all[fmt.Sprint(s)] = true
		}()
	}

	time.Sleep(time.Second * 2)

	nonExist := 0
	for k := range all {
		index, _ := strconv.Atoi(k)
		if index > addTasksLen-1 {
			nonExist++
		}
	}
	//获取任务
	tasks := c.dynamicTasks.Tasks
	assert.Equal(t, nonExist+addTasksLen, len(tasks), "任务长度")

	for _, v := range tasks {
		index := strings.TrimLeft(v.Cron, "task")
		if _, ok := all[index]; ok {
			assert.Equal(t, true, v.Disable, "任务状态1")
			continue
		}
		assert.Equal(t, false, v.Disable, "任务状态2")
	}
	//订阅者长度
	assert.Equal(t, coroutine+addTasksLen-noNotify, subscriber1, "订阅者接收任务通知数量1")
	assert.Equal(t, coroutine+addTasksLen-noNotify, subscriber2, "订阅者接收任务通知数量2")
}
