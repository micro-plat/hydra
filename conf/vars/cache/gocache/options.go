package gocache

import (
	"encoding/json"
	"time"
)

//Option 配置选项
type Option func(*GoCache)

//WithExpiration 默认过期时间
func WithExpiration(expiration time.Duration) Option {
	return func(o *GoCache) {
		o.Expiration = expiration
	}
}

//WithCleanupInterval 清理周期
func WithCleanupInterval(interval time.Duration) Option {
	return func(o *GoCache) {
		o.CleanupInterval = interval
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *GoCache) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
