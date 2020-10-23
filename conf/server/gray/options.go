package gray

type Option func(*Gray)

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Gray) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Gray) {
		a.Disable = false
	}
}
