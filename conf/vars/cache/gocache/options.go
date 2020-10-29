package gocache

import (
	"encoding/json"
	"time" 
	"github.com/micro-plat/hydra/conf/vars/cache"
	
)

//Options gocache客户端配置
type Options struct {
	*cache.Cache
	Expiration      time.Duration `json:"expiration,omitempty"`
	CleanupInterval time.Duration `json:"cleanup_interval,omitempty"`
}

//NewOptions x
func NewOptions(opts ...Option) *Options {
	opt := Options{
		Expiration:      5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
		Cache:           &cache.Cache{Proto: "gocache"},
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
