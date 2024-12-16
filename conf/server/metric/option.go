package metric

//Option 配置选项
type Option func(*Metric)

//WithUPName 设置用户名密码
func WithUPName(userName string, password string) Option {
	return func(a *Metric) {
		a.UserName = userName
		a.Password = password
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Metric) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Metric) {
		a.Disable = false
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *Metric) {
		a.EnableEncryption = true
	}
}
