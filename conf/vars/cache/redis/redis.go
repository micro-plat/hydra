package redis

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/cache"
)

//Redis redis缓存配置
type Redis struct {
	*cache.Cache
	ConfigName string `json:"config_name" toml:"config_name" valid:"required" `
}

//New 构建redis消息队列配置
func New(configName string, opts ...Option) *Redis {
	r := &Redis{
		ConfigName: configName,
		Cache:      &cache.Cache{Proto: "redis"},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *Redis {

	org := New("")
	if err := json.Unmarshal([]byte(raw), org); err != nil {
		panic(err)
	}

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
