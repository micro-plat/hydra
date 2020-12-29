package apikey

//Option 配置选项
type Option func(*APIKeyAuth)

//WithSecret 设置密钥
func WithSecret(secret string) Option {
	return func(a *APIKeyAuth) {
		a.Secret = secret
	}
}

//WithMode 设置为验证模式
func WithMode(m string) Option {
	return func(a *APIKeyAuth) {
		a.Mode = m
	}
}

//WithMD5Mode 设置为MD5验证模式
func WithMD5Mode() Option {
	return func(a *APIKeyAuth) {
		a.Mode = ModeMD5
	}
}

//WithSHA1Mode 设置为SHA1验证模式
func WithSHA1Mode() Option {
	return func(a *APIKeyAuth) {
		a.Mode = ModeSHA1
	}
}

//WithSHA256Mode 设置为SHA256验证模式
func WithSHA256Mode() Option {
	return func(a *APIKeyAuth) {
		a.Mode = ModeSHA256
	}
}

//WithExcludes 指定需要排除验证的服务名列表
func WithExcludes(p ...string) Option {
	return func(a *APIKeyAuth) {
		a.Excludes = p
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *APIKeyAuth) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *APIKeyAuth) {
		a.Disable = false
	}
}

//WithInvoker 排除的服务或请求
func WithInvoker(v string) Option {
	return func(b *APIKeyAuth) {
		b.Invoker = v
	}
}
