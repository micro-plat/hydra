package task

type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

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
