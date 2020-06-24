package rlog

//Option 配置选项
type Option func(*Layout)

//WithLayout 设置日志级别为info
func WithLayout(v string) Option {
	return func(a *Layout) {
		a.Layout = v
	}
}

//WithDisable 禁用任务
func WithDisable() Option {
	return func(a *Layout) {
		a.Disable = true
	}
}

//WithEnable 启用任务
func WithEnable() Option {
	return func(a *Layout) {
		a.Disable = false
	}
}

//WithInfo 设置日志级别为info
func WithInfo() Option {
	return func(a *Layout) {
		a.Level = "Info"
	}
}

//WithOff 设置日志级别为info
func WithOff() Option {
	return func(a *Layout) {
		a.Level = "Off"
	}
}

//WithWarn 设置日志级别为Warn
func WithWarn() Option {
	return func(a *Layout) {
		a.Level = "Warn"
	}
}

//WithError 设置日志级别为Error
func WithError() Option {
	return func(a *Layout) {
		a.Level = "Error"
	}
}

//WithFatal 设置日志级别为Fatal
func WithFatal() Option {
	return func(a *Layout) {
		a.Level = "Fatal"
	}
}

//WithDebug 设置日志级别为Debug
func WithDebug() Option {
	return func(a *Layout) {
		a.Level = "Debug"
	}
}

//WithAll 设置日志级别为All
func WithAll() Option {
	return func(a *Layout) {
		a.Level = "All"
	}
}
