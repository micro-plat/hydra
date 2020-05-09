package ras

import "github.com/micro-plat/lib4go/types"

type remotingOption struct {

	//指定需要远程验证的请求列表
	Requests []string `json:"requests,omitempty" valid:"required" toml:"requests,omitempty"`

	//必须传入的字段列表
	Required []string `json:"required,omitempty" toml:"required,omitempty"`

	//远程验证的字段连接方式
	Connect *Connect `json:"connect,omitempty" toml:"connect,omitempty"`

	//字段别名表
	Alias map[string]string `json:"alias,omitempty" toml:"alias,omitempty"`

	//扩展参数表
	Params map[string]interface{} `json:"params,omitempty" toml:"params,omitempty"`

	//需要解密的字段
	Decrypt []string `json:"decrypt,omitempty" toml:"decrypt,omitempty"`

	//是否需要检查时间戳
	CheckTS bool `json:"check-timestamp,omitempty" toml:"check-timestamp,omitempty"`

	//配置是否禁用
	Disable bool `json:"disable,omitempty" toml:"disable,omitempty"`
}

//RemotingOption 配置选项
type RemotingOption func(*remotingOption)

//WithRequest 设置requests的请求服务路径
func WithRequest(path ...string) RemotingOption {
	return func(a *remotingOption) {
		if len(path) > 0 {
			a.Requests = append(a.Requests, path...)
		}
	}
}

//WithRequired 设置必须字段
func WithRequired(fieldName ...string) RemotingOption {
	return func(a *remotingOption) {
		if len(fieldName) > 0 {
			a.Required = append(a.Required, fieldName...)
		}
	}
}

//WithUIDAlias 设置用户id的字段名
func WithUIDAlias(name string) RemotingOption {
	return func(a *remotingOption) {
		a.Alias["euid"] = name
	}
}

//WithTimestampAlias 设置timestamp的字段名
func WithTimestampAlias(name string) RemotingOption {
	return func(a *remotingOption) {
		a.Alias["timestamp"] = name
	}
}

//WithSignAlias 设置sign的字段名
func WithSignAlias(name string) RemotingOption {
	return func(a *remotingOption) {
		a.Alias["sign"] = name
	}
}

//WithDecryptName 设置需要解密的字段名
func WithDecryptName(name ...string) RemotingOption {
	return func(a *remotingOption) {
		a.Decrypt = append(a.Decrypt, name...)
	}
}

//WithCheckTimestamp 设置需要检查时间戳
func WithCheckTimestamp(e ...bool) RemotingOption {
	return func(a *remotingOption) {
		a.CheckTS = types.GetBoolByIndex(e, 0, true)
	}
}

//WithParam 设置扩展参数
func WithParam(key string, value interface{}) RemotingOption {
	return func(a *remotingOption) {
		a.Params[key] = value
	}
}

//WithDisable 禁用配置
func WithDisable() RemotingOption {
	return func(a *remotingOption) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() RemotingOption {
	return func(a *remotingOption) {
		a.Disable = false
	}
}

//WithConnect 启用配置
func WithConnect(opts ...ConnectOption) RemotingOption {
	return func(a *remotingOption) {
		a.Connect = &Connect{connectOption: &connectOption{}}
		for _, opt := range opts {
			opt(a.Connect.connectOption)
		}
	}
}
