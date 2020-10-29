package queue

//Queue 配置参数
type Queue struct {
	Queue       string `json:"queue,omitempty" valid:"ascii,required" toml:"queue,omitempty"`
	Service     string `json:"service,omitempty" valid:"ascii,required" toml:"service,omitempty"`
	Concurrency int    `json:"concurrency,omitempty" toml:"concurrency,omitempty"`
	Disable     bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//NewQueue 构建queue任务信息
func NewQueue(queue string, service string) *Queue {
	return &Queue{
		Queue:   queue,
		Service: service,
	}
}

//NewQueueByConcurrency 构建queue任务信息
func NewQueueByConcurrency(queue string, service string, concurrency int) *Queue {
	return &Queue{
		Queue:       queue,
		Service:     service,
		Concurrency: concurrency,
	}
}
