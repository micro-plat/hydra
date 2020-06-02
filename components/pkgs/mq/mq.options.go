package mq

import "encoding/json"

//ConfOpt 配置选项
type ConfOpt struct {
	Address     string `json:"address"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	CertPath    string `json:"cert_path"`
	DialTimeout int64  `json:"dial_timeout"`
	Version     string `json:"version"`
	Persistent  string `json:"persistent"`
	Ack         string `json:"ack"`
	Retry       bool   `json:"retry"`
	Key         string `json:"key"`
	Raw         string `json:"raw"`
}

//Option 配置选项
type Option func(*ConfOpt)

//WithConf 根据配置文件初始化
func WithConf(conf *ConfOpt) Option {
	return func(o *ConfOpt) {
		o = conf
	}
}

//WithVersion 设置版本号
func WithVersion(version string) Option {
	return func(o *ConfOpt) {
		o.Version = version
	}
}

//WithPersistent 设置数据模式
func WithPersistent(persistent string) Option {
	return func(o *ConfOpt) {
		o.Persistent = persistent
	}
}

//WithAck 设置客户端确认模式
func WithAck(ack string) Option {
	return func(o *ConfOpt) {
		o.Ack = ack
	}
}

//WithRetrySend 发送失败后重试
func WithRetrySend(b bool) Option {
	return func(o *ConfOpt) {
		o.Retry = b
	}
}

//WithKey 设置数据加密key
func WithKey(key string) Option {
	return func(o *ConfOpt) {
		o.Key = key
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw []byte) Option {
	return func(o *ConfOpt) {
		if err := json.Unmarshal(raw, o); err != nil {
			panic(err)
		}
		o.Raw = string(raw)
	}
}
