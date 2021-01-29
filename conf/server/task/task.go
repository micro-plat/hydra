package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/lib4go/security/md5"
)

//CronExecuteImmediately 立即执行
const CronExecuteImmediately = "@immediately"

//CronExecuteNow 立即执行
const CronExecuteNow = "@now"

//Task cron任务的task明细
type Task struct {
	Cron    string `json:"cron,omitempty" valid:"ascii,required" toml:"cron,omitempty" label:"任务名称"`
	Service string `json:"service,omitempty" valid:"ascii,spath,required" toml:"service,omitempty" label:"任务服务"`
	Disable bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//NewTask 创建任务信息
func NewTask(cron string, service string, opts ...Option) *Task {
	t := &Task{
		Cron:    cron,
		Service: service,
	}

	for _, f := range opts {
		f(t)
	}
	return t
}

//GetUNQ 获取任务的唯一标识
func (t *Task) GetUNQ() string {
	return md5.Encrypt(fmt.Sprintf("%s(%s)", t.Service, t.Cron))
}

//IsImmediately 是否立即
func (t *Task) IsImmediately() bool {
	return t.Cron == CronExecuteNow || t.Cron == CronExecuteImmediately
}

//Validate 验证任务参数
func (t *Task) Validate() error {
	if b, err := govalidator.ValidateStruct(t); !b && err != nil {
		return fmt.Errorf("task配置有误:%v", err)
	}
	return nil
}
