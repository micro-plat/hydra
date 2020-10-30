package memcached

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/cache"
)

//Memcache memcache客户端配置
type Memcache struct {
	*cache.Cache
	Address      []string `json:"addrs"   toml:"addrs" valid:"required" `
	Timeout      int      `json:"timeout,omitempty"  toml:"timeout,omitempty"`
	MaxIdleConns int      `json:"max_idle_conns,omitempty"  toml:"max_idle_conns,omitempty"`
}

//New x
func New(opts ...Option) *Memcache {
	opt := Memcache{
		Timeout:      1,
		MaxIdleConns: 10,
		Cache:        &cache.Cache{Proto: "memcache"},
	}
	for i := range opts {
		opts[i](&opt)
	}
	return &opt
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *Memcache {

	org := New()
	if err := json.Unmarshal([]byte(raw), org); err != nil {
		panic(err)
	}

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
