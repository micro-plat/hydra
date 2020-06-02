package redis

import "github.com/micro-plat/hydra/conf/server/queue"

//Redis redis缓存配置
type Redis struct {
	*queue.Queue
	Address      string `json:"addrs,omitempty"`
	DbIndex      int    `json:"db,omitempty"`
	DialTimeout  int    `json:"dial_timeout,omitempty"`
	ReadTimeout  int    `json:"read_timeout,omitempty"`
	WriteTimeout int    `json:"write_timeout,omitempty"`
	PoolSize     int    `json:"pool_size,omitempty"`
}

//New 构建redis消息队列配置
func New(address string, opts ...Option) *Redis {
	r := &Redis{
		Address:      address,
		Queue:        &queue.Queue{Proto: "redis"},
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
