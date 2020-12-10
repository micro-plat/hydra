package cron

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/robfig/cron"
)

//CronTask 定时任务
type CronTask struct {
	*task.Task
	Counter  *Counter
	Round    *Round
	schedule cron.Schedule
	method   string
	form     map[string]interface{}
	header   map[string]string
}

//NewCronTask 构建定时任务
func NewCronTask(t *task.Task) (r *CronTask, err error) {
	r = &CronTask{
		Task:    t,
		Counter: &Counter{},
		Round:   &Round{},
		method:  "GET",
		form:    make(map[string]interface{}),
		header:  map[string]string{"Client-IP": "127.0.0.1"},
	}
	if t.IsImmediately() {
		return r, nil
	}

	r.schedule, err = cron.ParseStandard(t.Cron)
	if err != nil {
		return r, fmt.Errorf("%s的cron表达式(%s)配置有误 %w", t.Service, t.Cron, err)
	}
	return r, nil
}

//GetName 获取任务名称
func (m *CronTask) GetName() string {
	return m.Task.GetUNQ()
}

//NextTime 下次执行时间
func (m *CronTask) NextTime(t time.Time) time.Time {
	if m.IsImmediately() {
		return time.Now()
	}
	return m.schedule.Next(t)
}

//GetService 服务名
func (m *CronTask) GetService() string {
	return m.Task.Service
}

//GetMethod 方法名
func (m *CronTask) GetMethod() string {
	return m.method
}

//GetForm 输入参数
func (m *CronTask) GetForm() map[string]interface{} {
	return m.form
}

//GetHeader 头信息
func (m *CronTask) GetHeader() map[string]string {
	return m.header
}
