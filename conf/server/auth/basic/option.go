package basic

//Option 配置选项
type Option func(*BasicAuth)

//WithUP 添加用户名密码
func WithUP(userName string, pwd string) Option {
	return func(b *BasicAuth) {
		b.Members[userName] = pwd
	}
}

//WithInvoker 排除的服务或请求
func WithInvoker(v string) Option {
	return func(b *BasicAuth) {
		b.Invoker = v
	}
}

//WithExcludes 排除的服务或请求
func WithExcludes(p ...string) Option {
	return func(b *BasicAuth) {
		b.Excludes = p
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(b *BasicAuth) {
		b.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(b *BasicAuth) {
		b.Disable = false
	}
}
