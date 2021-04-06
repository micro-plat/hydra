package apm

//Option 配置选项
type Option func(*APM)

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *APM) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *APM) {
		a.Disable = false
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *APM) {
		a.EnableEncryption = true
	}
}
