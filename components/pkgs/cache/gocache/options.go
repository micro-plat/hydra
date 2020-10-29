package gocache

import (
	"encoding/json"
	"time"
)

//Options gocache客户端配置
type Options struct {
	Expiration      time.Duration `json:"expiration"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

//NewOptions x
func NewOptions(opts ...Option) *Options {
	opt := Options{
		Expiration:      5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
	}
	for i := range opts {
		opts[i](&opt)
	}
	return &opt
}

//Option 配置选项
type Option func(*Options)

//WithExpiration 默认过期时间
func WithExpiration(expiration time.Duration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}

//WithCleanupInterval 清理周期
func WithCleanupInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.CleanupInterval = interval
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *Options) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
