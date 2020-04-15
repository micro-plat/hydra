package queue

//option 配置参数
type option struct {
	Queue       string `json:"queue" valid:"ascii,required"`
	Service     string `json:"service" valid:"ascii,required"`
	Concurrency int    `json:"concurrency,omitempty"`
	Disable     bool   `json:"disable,omitemptye"`
}

//Option 配置选项
type Option func(*option)

//WithQueue 构建queue任务信息
func WithQueue(queue string, service string) Option {
	return func(a *option) {
		a.Queue = queue
		a.Service = service
	}
}

//WithQueueByConcurrency 构建queue任务信息
func WithQueueByConcurrency(queue string, service string, concurrency int) Option {
	return func(a *option) {
		a.Queue = queue
		a.Service = service
		a.Concurrency = concurrency
	}
}
