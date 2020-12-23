package cacheredis

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/redis"
)

//Redis redis缓存配置
type Redis struct {
	*cache.Cache
	*redis.Redis
	ConfigName string `json:"config_name,omitempty" toml:"config_name,omitempty" valid:"ascii"`
}

//New 构建redis消息队列配置
func New(address string, opts ...Option) (org *Redis) {
	org = &Redis{
		Cache: &cache.Cache{Proto: "redis"},
		Redis: redis.New(address),
	}

	for _, opt := range opts {
		opt(org)
	}
	b, err := govalidator.ValidateStruct(org)
	if !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	if org.ConfigName == "" && (reflect.ValueOf(org.Redis).IsNil() || len(org.Addrs) == 0) {
		panic(fmt.Errorf("redis配置数据有误:至少存在Addrs或ConfigName,%+v", org))
	}
	return org
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) (org *Redis) {
	return New("", WithRaw(raw))
}
