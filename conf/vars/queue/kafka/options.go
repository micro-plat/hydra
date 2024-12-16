package kafka

import (
	"encoding/json"
)

// Option 配置选项
type Option func(*Kafka)

// WithAddrs 设置Addrs
func WithAddrs(addrs ...string) Option {
	return func(a *Kafka) {
		a.Addrs = addrs
	}
}

// WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *Kafka) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}

// WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *Kafka) {
		a.EnableEncryption = true
	}
}

// WithOffset 设置偏移量
func WithOffset(offsetNewest bool) Option {
	return func(a *Kafka) {
		if offsetNewest {
			a.Offset = -1
		} else {
			a.Offset = -2
		}
	}
}

// WithTimeout 设置数据库连接超时，读写超时时间
func WithTimeout(writeTimeout int) Option {
	return func(a *Kafka) {
		a.WriteTimeout = writeTimeout
	}
}

// WithGroup 设置消费者分组
func WithGroup(group string) Option {
	return func(a *Kafka) {
		a.Group = group
	}
}
