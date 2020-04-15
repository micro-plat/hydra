package jwt

//jwtOption jwt配置参数
type jwtOption struct {
	Name       string   `json:"name" valid:"ascii,required"`
	ExpireAt   int64    `json:"expireAt" valid:"required"`
	Mode       string   `json:"mode" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required"`
	Secret     string   `json:"secret" valid:"ascii,required"`
	Source     string   `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)"`
	Exclude    []string `json:"exclude,omitempty"`
	FailedCode string   `json:"failed-code,omitempty" valid:"numeric,range(400|999)"`
	Redirect   string   `json:"redirect,omitempty" valid:"ascii"`
	Domain     string   `json:"domain,omitempty" valid:"ascii"`
	Disable    bool     `json:"disable,omitempty"`
}

//Option jwt配置选项
type Option func(*jwtOption)

//WithName jwt设置到cookie中或header中的名称
func WithName(name string) Option {
	return func(a *jwtOption) {
		a.Name = name
	}
}

//WithExpireAt 过期时间以秒为单位
func WithExpireAt(expireAt int64) Option {
	return func(a *jwtOption) {
		a.ExpireAt = expireAt
	}
}

//WithMode jwt加解密模式
func WithMode(mode string) Option {
	return func(a *jwtOption) {
		a.Mode = mode
	}
}

//WithSecret jwt加密密钥
func WithSecret(secret string) Option {
	return func(a *jwtOption) {
		a.Secret = secret
	}
}

//WithHeader 将jwt存储到http头中
func WithHeader() Option {
	return func(a *jwtOption) {
		a.Source = "HEADER"
	}
}

//WithCookie 将jwt存储到cookie中
func WithCookie() Option {
	return func(a *jwtOption) {
		a.Source = "COOKIE"
	}
}

//WithExclude 排除的服务或请求
func WithExclude(p ...string) Option {
	return func(a *jwtOption) {
		a.Exclude = p
	}
}

//WithFailedCode 设置验证失败的错误码（400-999）
func WithFailedCode(code string) Option {
	return func(a *jwtOption) {
		a.FailedCode = code
	}
}

//WithRedirect 设置转跳URL
func WithRedirect(url string) Option {
	return func(a *jwtOption) {
		a.Redirect = url
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *jwtOption) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *jwtOption) {
		a.Disable = false
	}
}
