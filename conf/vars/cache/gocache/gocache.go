package gocache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf/vars/cache"
)

//GoCache gocache客户端配置
type GoCache struct {
	*cache.Cache
	Expiration      time.Duration `json:"expiration,omitempty"  toml:"expiration,omitempty"`
	CleanupInterval time.Duration `json:"cleanup_interval,omitempty" toml:"cleanup_interval,omitempty"`
}

//New x
func New(opts ...Option) *GoCache {
	opt := GoCache{
		Expiration:      5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
		Cache:           &cache.Cache{Proto: "gocache"},
	}
	for i := range opts {
		opts[i](&opt)
	}
	return &opt
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *GoCache {

	org := New()
	if err := json.Unmarshal([]byte(raw), org); err != nil {
		panic(err)
	}

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
