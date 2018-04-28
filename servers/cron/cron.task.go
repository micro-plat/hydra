package cron

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/logger"
	"github.com/zkfy/cron"
)

type iCronTask interface {
	GetName() string
	ReduceRound(int)
	GetRound() int
	SetRound(int)
	GetExecuted() int
	AddExecuted()
	NextTime(time.Time) time.Time
	GetHandler() interface{}
	GetTaskExecutionRecord() (string, error)
	SetResult(status int, result []byte)
	dispatcher.IRequest
	logger.ILogger
}
type cronTask struct {
	*conf.Task
	schedule cron.Schedule
	Executed int `json:"executed"`
	round    int
	method   string
	form     map[string]interface{}
	header   map[string]string
	logger.ILogger
	status int
	result []byte
}

func newCronTask(t *conf.Task) (r *cronTask, err error) {
	r = &cronTask{
		Task:    t,
		method:  "GET",
		header:  make(map[string]string),
		ILogger: logger.GetSession(t.Name, logger.CreateSession()),
	}
	r.schedule, err = cron.ParseStandard(t.Cron)
	if err != nil {
		return nil, fmt.Errorf("%s的cron表达式(%s)配置有误", t.Name, t.Cron)
	}
	r.form = t.Input
	if r.form == nil {
		r.form = make(map[string]interface{})
	}
	return
}
func (m *cronTask) GetName() string {
	return m.Task.Name
}
func (m *cronTask) ReduceRound(v int) {
	m.round -= v
}

func (m *cronTask) GetRound() int {
	return m.round
}
func (m *cronTask) SetRound(v int) {
	m.round = v
}
func (m *cronTask) GetExecuted() int {
	return m.Executed
}
func (m *cronTask) AddExecuted() {
	if m.Executed >= math.MaxInt32 {
		m.Executed = 1
	} else {
		m.Executed++
	}
}
func (m *cronTask) GetHandler() interface{} {
	return m.Handler
}
func (m *cronTask) NextTime(t time.Time) time.Time {
	return m.schedule.Next(t)
}
func (m *cronTask) GetService() string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(m.Name, "/"))
}
func (m *cronTask) GetMethod() string {
	return m.method
}
func (m *cronTask) GetForm() map[string]interface{} {
	return m.form
}
func (m *cronTask) GetHeader() map[string]string {
	return m.header
}
func (m *cronTask) SetResult(status int, result []byte) {
	m.status = status
	m.result = result
}
func (m *cronTask) GetTaskExecutionRecord() (string, error) {
	data := map[string]interface{}{
		"name":     m.Name,
		"cron":     m.Cron,
		"service":  m.Service,
		"engine":   m.Engine,
		"executed": m.Executed,
		"result":   fmt.Sprintf("%d,%s", m.status, json.RawMessage(m.result)),
		"next":     m.NextTime(time.Now()).Format("20060102150405"),
		"last":     time.Now().Format("20060102150405"),
	}
	buff, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
