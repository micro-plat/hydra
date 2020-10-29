package redis

import (
	"github.com/micro-plat/hydra/conf/vars/cache"
)

//Options redis缓存配置
type Options struct {
	*cache.Cache
	Address      []string `json:"addrs,required"`
	DbIndex      int      `json:"db,omitempty"`
	DialTimeout  int      `json:"dial_timeout,omitempty"`
	ReadTimeout  int      `json:"read_timeout,omitempty"`
	WriteTimeout int      `json:"write_timeout,omitempty"`
	PoolSize     int      `json:"pool_size,omitempty"`
}

//NewOptions 构建redis消息队列配置
func NewOptions(address []string, opts ...Option) *Options {
	r := &Options{
		Address:      address,
		Cache:        &cache.Cache{Proto: "redis"},
		DbIndex:      1,
		DialTimeout:  10,
		ReadTimeout:  10,
		WriteTimeout: 10,
		PoolSize:     10,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//Option 配置选项
type Option func(*Options)

//WithDbIndex 设置数据库分片索引
func WithDbIndex(i int) Option {
	return func(a *Options) {
		a.DbIndex = i
	}
}

//WithTimeout 设置数据库连接超时，读写超时时间
func WithTimeout(dialTimeout int, readTimeout int, writeTimeout int) Option {
	return func(a *Options) {
		a.DialTimeout = dialTimeout
		a.ReadTimeout = readTimeout
		a.WriteTimeout = writeTimeout
	}
}

//WithPoolSize 设置数据库连接池大小
func WithPoolSize(i int) Option {
	return func(a *Options) {
		a.PoolSize = i
	}
}
