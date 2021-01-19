package queueredis

import (
	"encoding/json"
	"fmt"
)

//Option 配置选项
type Option func(*Redis)

//WithConfigName 设置数据库分片索引
func WithConfigName(configName string) Option {
	return func(a *Redis) {
		a.ConfigName = configName
	}
}

//WithAddrs 设置Addrs
func WithAddrs(addrs ...string) Option {
	return func(a *Redis) {
		a.Addrs = append(a.Addrs, addrs...)
	}
}

//WithDbIndex 设置数据库分片索引
func WithDbIndex(i int) Option {
	return func(a *Redis) {
		a.DbIndex = i
	}
}

//WithTimeout 设置数据库连接超时，读写超时时间
func WithTimeout(dialTimeout int, readTimeout int, writeTimeout int) Option {
	return func(a *Redis) {
		a.DialTimeout = dialTimeout
		a.ReadTimeout = readTimeout
		a.WriteTimeout = writeTimeout
	}
}

//WithPoolSize 设置数据库连接池大小
func WithPoolSize(i int) Option {
	return func(a *Redis) {
		a.PoolSize = i
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *Redis) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(fmt.Errorf("queueredis.WithRaw:%w", err))
		}
	}
}
