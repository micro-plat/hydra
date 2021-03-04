package queue

import "net/http"

//Queue 配置参数
type Queue struct {
	Queue       string `json:"queue,omitempty" valid:"ascii,required" toml:"queue,omitempty" label:"队列名"`
	Service     string `json:"service,omitempty" valid:"ascii,spath,required" toml:"service,omitempty" label:"队列服务"`
	Concurrency int    `json:"concurrency,omitempty" toml:"concurrency,omitempty"`
	Disable     bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//GetActions GetActions
func (q *Queue) GetActions() []string {
	return []string{http.MethodGet}
}

//GetPath GetPath
func (q *Queue) GetPath() string {
	return q.Service
}

//GetService GetService
func (q *Queue) GetService() string {
	return q.Service
}

//NewQueue 构建queue任务信息
func NewQueue(queue string, service string, opts ...Option) *Queue {
	q := &Queue{
		Queue:   queue,
		Service: service,
	}
	for i := range opts {
		opts[i](q)
	}
	return q
}

//Option Option
type Option func(q *Queue)

//WithConcurrency 并发数
func WithConcurrency(concurrency int) Option {
	return func(q *Queue) {
		q.Concurrency = concurrency
	}
}

//WithDisable 禁用
func WithDisable() Option {
	return func(q *Queue) {
		q.Disable = true
	}
}

//WithEnable 启用
func WithEnable() Option {
	return func(q *Queue) {
		q.Disable = false
	}
}
