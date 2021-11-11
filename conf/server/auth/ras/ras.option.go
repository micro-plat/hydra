package ras

import (
	"github.com/micro-plat/hydra/conf"
)

//Option 配置选项
type Option func(*RASAuth)

//WithParam 设置扩展参数
func WithAuths(auths ...*Auth) Option {
	return func(ras *RASAuth) {
		for _, a := range auths {
			a.PathMatch = conf.NewPathMatch(a.Requests...)
		}
		ras.Auth = append(ras.Auth, auths...)
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *RASAuth) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *RASAuth) {
		a.Disable = false
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *RASAuth) {
		a.EnableEncryption = true
	}
}

//WithExcludes 排除加密路径
func WithExcludes(p ...string) Option {
	return func(a *RASAuth) {
		a.Excludes = p
	}
}
