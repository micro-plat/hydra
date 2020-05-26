package fsa

type fixedOption struct {
}

//FixedOption 配置选项
type FixedOption func(*FixedSecretAuth)

//WithSecret 设置密钥
func WithSecret(secret string) FixedOption {
	return func(a *FixedSecretAuth) {
		a.Secret = secret
	}
}

//WithMD5Mode 设置为MD5验证模式
func WithMD5Mode() FixedOption {
	return func(a *FixedSecretAuth) {
		a.Mode = "MD5"
	}
}

//WithSHA1Mode 设置为SHA1验证模式
func WithSHA1Mode() FixedOption {
	return func(a *FixedSecretAuth) {
		a.Mode = "SHA1"
	}
}

//WithSHA256Mode 设置为SHA256验证模式
func WithSHA256Mode() FixedOption {
	return func(a *FixedSecretAuth) {
		a.Mode = "SHA256"
	}
}

//WithInclude 指定需要进行签名验证的服务名列表
func WithInclude(p ...string) FixedOption {
	return func(a *FixedSecretAuth) {
		a.Include = p
	}
}

//WithDisable 禁用配置
func WithDisable() FixedOption {
	return func(a *FixedSecretAuth) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() FixedOption {
	return func(a *FixedSecretAuth) {
		a.Disable = false
	}
}
