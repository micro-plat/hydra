package jwt

//Option jwt配置选项
type Option func(*JWTAuth)

//WithName jwt设置到cookie中或header中的名称
func WithName(name string) Option {
	return func(a *JWTAuth) {
		a.Name = name
	}
}

//WithExpireAt 过期时间以秒为单位
func WithExpireAt(expireAt int64) Option {
	return func(a *JWTAuth) {
		a.ExpireAt = expireAt
	}
}

//WithMode jwt加解密模式
func WithMode(mode string) Option {
	return func(a *JWTAuth) {
		a.Mode = mode
	}
}

//WithSecret jwt加密密钥
func WithSecret(secret string) Option {
	return func(a *JWTAuth) {
		a.Secret = secret
	}
}

//WithHeader 将jwt存储到http头中
func WithHeader() Option {
	return func(a *JWTAuth) {
		a.Source = "HEADER"
	}
}

//WithCookie 将jwt存储到cookie中
func WithCookie() Option {
	return func(a *JWTAuth) {
		a.Source = "COOKIE"
	}
}

//WithExcludes 排除的服务或请求
func WithExcludes(p ...string) Option {
	return func(a *JWTAuth) {
		a.Excludes = p
	}
}

//WithRedirect 设置转跳URL
func WithRedirect(url string) Option {
	return func(a *JWTAuth) {
		a.Redirect = url
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *JWTAuth) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *JWTAuth) {
		a.Disable = false
	}
}
