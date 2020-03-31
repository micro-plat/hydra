package conf

type Queues struct {
	Setting map[string]string `json:"args,omitempty"`
	Queues  []*Queue          `json:"queues"`
}

type Queue struct {
	Name        string            `json:"name,omitempty" valid:"ascii"`
	Queue       string            `json:"queue" valid:"ascii,required"`
	Engine      string            `json:"engine,omitempty"  valid:"ascii,uppercase,in(*|RPC)"`
	Service     string            `json:"service" valid:"ascii,required"`
	Setting     map[string]string `json:"args,omitempty"`
	Concurrency int               `json:"concurrency,omitempty"`
	Disable     bool              `json:"disable,omitemptye"`
	Handler     interface{}       `json:"-"`
}

//NewQueues 构建Queue注册列表
func NewQueues() *Queues {
	return &Queues{
		Queues: make([]*Queue, 0),
	}
}

//Append 添加Queue注册信息
func (h *Queues) Append(queue string, service string) *Queues {
	h.Queues = append(h.Queues, &Queue{
		Queue:   queue,
		Service: service,
	})
	return h
}

//NewQueue 构建Queue注册项
func NewQueue(queue string, service string) *Queue {
	return &Queue{
		Queue:   queue,
		Service: service,
	}
}

//NewQueueWithConcurrency 构建Queue注册项
func NewQueueWithConcurrency(queue string, service string, concurrency int) *Queue {
	return &Queue{
		Queue:       queue,
		Service:     service,
		Concurrency: concurrency,
	}
}

//WithConcurrency 设置并发协程数
func (t *Queue) WithConcurrency(c int) *Queue {
	t.Concurrency = c
	return t
}

//WithEnable 启用任务
func (t *Queue) WithEnable() *Queue {
	t.Disable = false
	return t
}

//WithDisable 禁用任务
func (t *Queue) WithDisable() *Queue {
	t.Disable = false
	return t
}
