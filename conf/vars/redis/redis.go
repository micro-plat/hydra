package redis

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/cache"
)

//TypeNodeName 分类节点名
const TypeNodeName = "redis"

//Redis redis缓存配置
type Redis struct {
	*cache.Cache
	Address      []string `json:"addrs,required" toml:"addrs,required"`
	Password     string   `json:"password,omitempty" toml:"password,omitempty"`
	DbIndex      int      `json:"db,omitempty" toml:"db,omitempty"`
	DialTimeout  int      `json:"dial_timeout,omitempty" toml:"dial_timeout,omitempty"`
	ReadTimeout  int      `json:"read_timeout,omitempty" toml:"read_timeout,omitempty"`
	WriteTimeout int      `json:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	PoolSize     int      `json:"pool_size,omitempty" toml:"pool_size,omitempty"`
}

//New 构建redis消息队列配置
func New(address []string, opts ...Option) *Redis {
	r := &Redis{
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

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *Redis {

	org := New(nil, WithRaw(raw))

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
