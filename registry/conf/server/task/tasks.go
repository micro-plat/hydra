package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/registry/conf"
)

//Tasks cron任务的task配置信息
type Tasks struct {
	Tasks []*Task `json:"tasks" toml:"tasks,omitempty"`
}

//NewEmptyTasks 构建空的tasks
func NewEmptyTasks() *Tasks {
	return &Tasks{
		Tasks: make([]*Task, 0),
	}
}

//NewTasks 构建任务列表
func NewTasks(tasks ...*Task) *Tasks {
	t := NewEmptyTasks()
	return t.Append(tasks...)
}

//Append 增加任务列表
func (t *Tasks) Append(tasks ...*Task) *Tasks {
	for _, task := range tasks {
		t.Tasks = append(t.Tasks, task)
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
