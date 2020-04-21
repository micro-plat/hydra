package cron

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry/conf/server/task"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/zkfy/cron"
)

var _ ICronTask = &CronTask{}

type ICronTask interface {
	GetName() string
	ReduceRound(int)
	GetRound() int
	UpdateRound(int)
	GetExecuted() int
	IncreaseExecuted()
	NextTime(time.Time) time.Time
	IsEnable() bool
	SetDisable()
	dispatcher.IRequest
}

//CronTask 定时任务
type CronTask struct {
	*task.Task
	schedule cron.Schedule
	Executed int
	round    int
	method   string
	form     map[string]interface{}
	header   map[string]string
	status   int
	result   []byte
}

//NewCronTask 构建定时任务
func NewCronTask(t *task.Task) (r *CronTask, err error) {
	r = &CronTask{
		Task:   t,
		method: "GET",
		form:   make(map[string]interface{}),
		header: make(map[string]string),
	}
	r.schedule, err = cron.ParseStandard(t.Cron)
	if err != nil {
		return nil, fmt.Errorf("%s的cron表达式(%s)配置有误", t.Service, t.Cron)
	}
	return
}

//GetName 获取任务名称
func (m *CronTask) GetName() string {
	return m.Task.GetUNQ()
}

//ReduceRound 减少任务执行轮数
func (m *CronTask) ReduceRound(v int) {
	m.round -= v
}

//GetRound 获取任务执行轮数
func (m *CronTask) GetRound() int {
	return m.round
}

//UpdateRound 获取任务的轮数
func (m *CronTask) UpdateRound(v int) {
	m.round = v
}

//GetExecuted 获取执行次数
func (m *CronTask) GetExecuted() int {
	return m.Executed
}

//IncreaseExecuted 累加执行次数
func (m *CronTask) IncreaseExecuted() {
	if m.Executed >= math.MaxInt32 {
		m.Executed = 1
	} else {
		m.Executed++
	}
}

//NextTime 下次执行时间
func (m *CronTask) NextTime(t time.Time) time.Time {
	return m.schedule.Next(t)
}

//GetService 服务名
func (m *CronTask) GetService() string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(m.GetUNQ(), "/"))
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

//IsEnable 是否启动
func (m *CronTask) IsEnable() bool {
	return !m.Disable
}

//SetDisable 禁用
func (m *CronTask) SetDisable() {
	m.Disable = true
}
