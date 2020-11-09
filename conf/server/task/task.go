package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/lib4go/security/md5"
)

//CronExxcuteOnce cron单次执行标识
const CronExxcuteOnce = "@once"

//CronExxcuteNow cron单次立即执行标识
const CronExxcuteNow = "@now"

//Task cron任务的task明细
type Task struct {
	Cron    string `json:"cron,omitempty" valid:"ascii,required" toml:"cron,omitempty"`
	Service string `json:"service,omitempty" valid:"ascii,required" toml:"service,omitempty"`
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

//IsOnce 是否只需要处理一次
func (t *Task) IsOnce() bool {
	return t.Cron == CronExxcuteOnce || t.Cron == CronExxcuteNow
}

//Validate 验证任务参数
func (t *Task) Validate() error {
	if b, err := govalidator.ValidateStruct(t); !b && err != nil {
		return fmt.Errorf("task配置有误:%v", err)
	}
	return nil
}
