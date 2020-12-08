package cron

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/zkfy/cron"
)

func TestNewCronTask(t *testing.T) {
	tests := []struct {
		name    string
		args    *task.Task
		ok      bool
		wantErr bool
	}{
		{name: "1. cron任务对象-初始化空任务", ok: false, args: &task.Task{}, wantErr: false},
		{name: "2. cron任务对象-初始化一次性任务", ok: false, args: &task.Task{Cron: task.CronExecuteImmediately, Service: "/servderi1"}, wantErr: true},
		{name: "3. cron任务对象-初始化一次性任务1", ok: false, args: &task.Task{Cron: task.CronExecuteNow, Service: "/servderi1"}, wantErr: true},
		{name: "4. cron任务对象-初始化错误的cron表达式", ok: false, args: &task.Task{Cron: "sdsd", Service: "/servderi1"}, wantErr: false},
		{name: "5. cron任务对象-初始化正确的任务", ok: true, args: &task.Task{Cron: "@every 10s", Service: "/servderi1"}, wantErr: true},
	}
	for _, tt := range tests {
		m := &CronTask{
			Task:    tt.args,
			Counter: &Counter{},
			Round:   &Round{},
			method:  "GET",
			form:    make(map[string]interface{}),
			header:  map[string]string{"Client-IP": "127.0.0.1"},
		}
		if tt.ok {
			s, err := cron.ParseStandard(tt.args.Cron)
			assert.Equalf(t, true, err == nil, tt.name, err)
			m.schedule = s
		}
		gotR, err := NewCronTask(tt.args)
		assert.Equalf(t, tt.wantErr, err == nil, tt.name, err)
		assert.Equal(t, m, gotR, tt.name)
	}
}

func TestCronTask_GetName(t *testing.T) {
	m := &CronTask{
		Task:    &task.Task{Cron: "@every 10h", Service: "/xxxxx1"},
		Counter: &Counter{},
		Round:   &Round{},
		method:  "get post",
		form:    map[string]interface{}{"taosy": "test"},
		header:  map[string]string{"Client-IP": "192.168.0.101", "Host": "www.baidu.com"},
	}

	s, err := cron.ParseStandard(m.Task.Cron)
	assert.Equalf(t, true, err == nil, "schedule creontask 错误")
	m.schedule = s

	got := m.GetName()
	assert.Equal(t, md5.Encrypt(fmt.Sprintf("%s(%s)", m.Task.Service, m.Task.Cron)), got, "获取任务的唯一索引名失败")

	got = m.GetService()
	assert.Equal(t, "/xxxxx1", got, "获取任务的services失败")

	got = m.GetMethod()
	assert.Equal(t, "get post", got, "获取任务的method失败")

	got1 := m.GetForm()
	assert.Equal(t, map[string]interface{}{"taosy": "test"}, got1, "获取任务的GetForm失败")

	got2 := m.GetHeader()
	assert.Equal(t, map[string]string{"Client-IP": "192.168.0.101", "Host": "www.baidu.com"}, got2, "获取任务的GetForm失败")
}
