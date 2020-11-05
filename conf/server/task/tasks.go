package task

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//Tasks cron任务的task配置信息
type Tasks struct {
	Tasks []*Task `json:"tasks,omitempty" toml:"tasks,omitempty"`
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

//Append 增加任务列表 @fix:存在的数据进行修改,不存在则添加 @hj
func (t *Tasks) Append(tasks ...*Task) *Tasks {
	keyMap := map[string]*Task{}
	for _, v := range t.Tasks {
		keyMap[v.GetUNQ()] = v
	}
	nonExistTask := []*Task{}
	for _, v := range tasks {
		if task, ok := keyMap[v.GetUNQ()]; ok {
			task.Disable = v.Disable
			continue
		}
		nonExistTask = append(nonExistTask, v)
	}
	t.Tasks = append(t.Tasks, nonExistTask...)
	return t
}

//GetConf 根据服务嚣配置获取task
func GetConf(cnf conf.IServerConf) (tasks *Tasks, err error) {
	tasks = &Tasks{}
	_, err = cnf.GetSubObject("task", tasks)
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("task:%v", err)
	}
	if err == conf.ErrNoSetting {
		return &Tasks{}, nil
	}

	for _, task := range tasks.Tasks {
		if b, err := govalidator.ValidateStruct(task); !b {
			return nil, fmt.Errorf("task配置有误:%v", err)
		}
	}
	return tasks, nil
}
