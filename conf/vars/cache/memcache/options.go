package memcache

import (
	"github.com/micro-plat/hydra/conf/vars/cache"
	"encoding/json"
)

//Conf memcache客户端配置
type Options struct {
	*cache.Cache
	Address      []string `json:"addrs,required"`
	Timeout      int      `json:"timeout,omitempty"`
	MaxIdleConns int      `json:"max_idle_conns,omitempty"`
}

func NewOptions(opts ...Option) *Options {
	opt := Options{
		Timeout:      1,
		MaxIdleConns: 10,
		Cache:        &cache.Cache{Proto: "memcache"},
	}
	for i := range opts {
		opts[i](&opt)
	}
	return &opt
}

//Option 配置选项
type Option func(*Options)

//WithAddress 设置哨兵服务器
func WithAddress(address ...string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

//WithTimeout 设置连接超时时长
func WithTimeout(timeout int) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

//WithMaxIdleConns 设置连接空闲数
func WithMaxIdleConns(idleConns int) Option {
	return func(o *Options) {
		o.MaxIdleConns = idleConns
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
