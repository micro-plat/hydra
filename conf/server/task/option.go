package task

//Option 配置选项
type Option func(*Task)

//WithDisable 禁用任务
func WithDisable() Option {
	return func(a *Task) {
		a.Disable = true
	}
}

//WithEnable 启用任务
func WithEnable() Option {
	return func(a *Task) {
		a.Disable = false
	}
}
