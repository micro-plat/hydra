package proxy

//Option 配置选项
type Option func(*Proxy)

//WithScript 设置脚本
func WithScript(script string) Option {
	return func(a *Proxy) {
		a.Script = script
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *Proxy) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *Proxy) {
		a.Disable = false
	}
}
