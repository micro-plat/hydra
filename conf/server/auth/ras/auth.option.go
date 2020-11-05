package ras

import "github.com/micro-plat/lib4go/types"

//AuthOption 配置选项
type AuthOption func(*Auth)

//WithRequest 设置requests的请求服务路径
func WithRequest(path ...string) AuthOption {
	return func(a *Auth) {
		if len(path) > 0 {
			a.Requests = append(a.Requests, path...)
		}
	}
}

//WithRequired 设置必须字段
func WithRequired(fieldName ...string) AuthOption {
	return func(a *Auth) {
		if len(fieldName) > 0 {
			a.Required = append(a.Required, fieldName...)
		}
	}
}

//WithUIDAlias 设置用户euid的字段名
func WithUIDAlias(name string) AuthOption {
	return func(a *Auth) {
		a.Alias["euid"] = name
	}
}

//WithTimestampAlias 设置timestamp的字段名
func WithTimestampAlias(name string) AuthOption {
	return func(a *Auth) {
		a.Alias["timestamp"] = name
	}
}

//WithSignAlias 设置sign的字段名
func WithSignAlias(name string) AuthOption {
	return func(a *Auth) {
		a.Alias["sign"] = name
	}
}

//WithDecryptName 设置需要解密的字段名
func WithDecryptName(name ...string) AuthOption {
	return func(a *Auth) {
		a.Decrypt = append(a.Decrypt, name...)
	}
}

//WithCheckTimestamp 设置需要检查时间戳
func WithCheckTimestamp(e ...bool) AuthOption {
	return func(a *Auth) {
		a.CheckTS = types.GetBoolByIndex(e, 0, true)
	}
}

//WithParam 设置扩展参数
func WithParam(key string, value interface{}) AuthOption {
	return func(a *Auth) {
		a.Params[key] = value
	}
}

//WithAuthDisable 禁用配置
func WithAuthDisable() AuthOption {
	return func(a *Auth) {
		a.Disable = true
	}
}

//WithAuthEnable 启用配置
func WithAuthEnable() AuthOption {
	return func(a *Auth) {
		a.Disable = false
	}
}

//WithConnect 启用配置
func WithConnect(opts ...ConnectOption) AuthOption {
	return func(a *Auth) {
		a.Connect = &Connect{connectOption: &connectOption{}}
		for _, opt := range opts {
			opt(a.Connect.connectOption)
		}
	}
}
