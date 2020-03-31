package conf

type Tasks struct {
	Setting map[string]string `json:"args,omitempty"`
	Tasks   []*Task           `json:"tasks"`
}
type Task struct {
	Name    string                 `json:"name,omitempty" valid:"ascii"`
	Cron    string                 `json:"cron" valid:"ascii,required"`
	Input   map[string]interface{} `json:"-"`
	Engine  string                 `json:"engine,omitempty"  valid:"ascii,uppercase,in(*|RPC)"`
	Service string                 `json:"service"  valid:"ascii,required"`
	Setting map[string]string      `json:"args,omitempty"`
	Next    string                 `json:"-"`
	Last    string                 `json:"-"`
	Handler interface{}            `json:"-" valid:"-"`
	Disable bool                   `json:"disable,omitempty"`
}

//NewTasks 构建CRON任务列表
func NewTasks() *Tasks {
	return &Tasks{
		Tasks: make([]*Task, 0),
	}
}

//Append 添加路由信息
func (h *Tasks) Append(cron string, service string) *Tasks {
	h.Tasks = append(h.Tasks, &Task{
		Cron:    cron,
		Engine:  "*",
		Service: service,
	})
	return h
}

//NewTask 构建CRON任务
func NewTask(cron string, service string) *Task {
	return &Task{
		Cron:    cron,
		Engine:  "*",
		Service: service,
	}
}

//WithDisable 禁用任务
func (t *Task) WithDisable() *Task {
	t.Disable = true
	return t
}

//WithEnable 启用任务
func (t *Task) WithEnable() *Task {
	t.Disable = false
	return t
}
