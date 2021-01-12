package task

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName 分类节点名
const TypeNodeName = "task"

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
	t, _ := NewEmptyTasks().Append(tasks...)
	return t
}

//Append 增加任务列表
func (t *Tasks) Append(tasks ...*Task) (*Tasks, []*Task) {
	keyMap := map[string]*Task{}
	for _, v := range t.Tasks {
		keyMap[v.GetUNQ()] = v
	}
	notifyTasks := []*Task{}
	for _, v := range tasks {
		if task, ok := keyMap[v.GetUNQ()]; ok {
			if task.Disable != v.Disable {
				notifyTasks = append(notifyTasks, v)
				task.Disable = v.Disable
			}
			continue
		}
		notifyTasks = append(notifyTasks, v)
		t.Tasks = append(t.Tasks, v)
	}
	return t, notifyTasks
}

//GetConf 根据服务嚣配置获取task
func GetConf(cnf conf.IServerConf) (tasks *Tasks, err error) {
	tasks = &Tasks{}
	_, err = cnf.GetSubObject(TypeNodeName, tasks)
	if errors.Is(err, conf.ErrNoSetting) {
		return &Tasks{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("task:%v", err)
	}

	for _, task := range tasks.Tasks {
		if b, err := govalidator.ValidateStruct(task); !b {
			return nil, fmt.Errorf("task配置有误:%v", err)
		}
	}
	return tasks, nil
}
