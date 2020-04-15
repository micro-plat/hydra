package task

//option 配置参数
type option struct {
	Cron    string `json:"cron" valid:"ascii,required"`
	Service string `json:"service" valid:"ascii,required"`
	Disable bool   `json:"disable,omitemptye"`
}

//Option 配置选项
type Option func(*option)

//WithTask 构建task任务信息
func WithTask(cron string, service string) Option {
	return func(a *option) {
		a.Cron = cron
		a.Service = service
	}
}
