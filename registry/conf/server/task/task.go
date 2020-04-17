package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/lib4go/security/md5"
)

//Tasks cron任务的task配置信息
type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

//Task cron任务的task明细
type Task struct {
	*option
}

//GetUNQ 获取任务的唯一标识
func (t *Task) GetUNQ() string {
	return md5.Encrypt(fmt.Sprintf("%s-%s", t.Cron, t.Service))
}

//Validate 验证任务参数
func (t *Task) Validate() error {
	if b, err := govalidator.ValidateStruct(t); !b && err != nil {
		return fmt.Errorf("task配置有误:%v", err)
	}
	return nil
}

//NewTasks 构建任务列表
func NewTasks(first Option, tasks ...Option) *Tasks {
	t := &Tasks{
		Tasks: make([]*Task, 0),
	}
	ft := &Task{option: &option{}}
	first(ft.option)
	t.Tasks = append(t.Tasks, ft)
	for _, opt := range tasks {
		nopt := &Task{option: &option{}}
		opt(nopt.option)
		t.Tasks = append(t.Tasks, nopt)
	}
	return t
}

//GetConf 根据服务嚣配置获取task
func GetConf(cnf conf.IMainConf) (tasks *Tasks, err error) {
	if _, err = cnf.GetSubObject("task", &tasks); err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("task:%v", err)
	}
	if len(tasks.Tasks) > 0 {
		if b, err := govalidator.ValidateStruct(&tasks); !b {
			return nil, fmt.Errorf("task配置有误:%v", err)
		}
	}
	return
}
