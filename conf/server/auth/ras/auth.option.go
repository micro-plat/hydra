package ras

import "github.com/micro-plat/lib4go/types"

//Option 配置选项
type Option func(*Auth)

//WithRequest 设置requests的请求服务路径
func WithRequest(path ...string) Option {
	return func(a *Auth) {
		if len(path) > 0 {
			a.Requests = append(a.Requests, path...)
		}
	}
}

//WithRequired 设置必须字段
func WithRequired(fieldName ...string) Option {
	return func(a *Auth) {
		if len(fieldName) > 0 {
			a.Required = append(a.Required, fieldName...)
		}
	}
}

//WithUIDAlias 设置用户id的字段名
func WithUIDAlias(name string) Option {
	return func(a *Auth) {
		a.Alias["euid"] = name
	}
}

//WithTimestampAlias 设置timestamp的字段名
func WithTimestampAlias(name string) Option {
	return func(a *Auth) {
		a.Alias["timestamp"] = name
	}
}

//WithSignAlias 设置sign的字段名
func WithSignAlias(name string) Option {
	return func(a *Auth) {
		a.Alias["sign"] = name
	}
}

//WithDecryptName 设置需要解密的字段名
func WithDecryptName(name ...string) Option {
	return func(a *Auth) {
		a.Decrypt = append(a.Decrypt, name...)
	}
}

//WithCheckTimestamp 设置需要检查时间戳
func WithCheckTimestamp(e ...bool) Option {
	return func(a *Auth) {
		a.CheckTS = types.GetBoolByIndex(e, 0, true)
	}
}

//WithParam 设置扩展参数
func WithParam(key string, value interface{}) Option {
	return func(a *Auth) {
		a.Params[key] = value
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Auth) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Auth) {
		a.Disable = false
	}
}

//WithConnect 启用配置
func WithConnect(opts ...ConnectOption) Option {
	return func(a *Auth) {
		a.Connect = &Connect{connectOption: &connectOption{}}
		for _, opt := range opts {
			opt(a.Connect.connectOption)
		}
	}
}
