package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//Tasks cron任务的task配置信息
type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

//Task cron任务的task明细
type Task struct {
	*option
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

//GetTasks 根据服务嚣配置获取task
func GetTasks(cnf conf.IMainConf) (tasks *Tasks, err error) {
	if _, err = cnf.GetSubObject("task", &tasks); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("task:%v", err)
		return nil, err
	}
	if len(tasks.Tasks) > 0 {
		if b, err := govalidator.ValidateStruct(&tasks); !b {
			err = fmt.Errorf("task配置有误:%v", err)
			return nil, err
		}
	}
	return
}
