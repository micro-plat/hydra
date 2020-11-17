package render

//Option 渲染可选参数
type Option func(*Render)

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Render) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Render) {
		a.Disable = false
	}
}
