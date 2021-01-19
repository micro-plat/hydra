package redis

import (
	"errors"
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"

	"github.com/asaskevich/govalidator"
)

//TypeNodeName 分类节点名
const TypeNodeName = "redis"

//Redis redis缓存配置
type Redis struct {
	Addrs        []string `json:"addrs,omitempty" toml:"addrs,omitempty" valid:"required" label:"集群地址(|分割)"`
	Password     string   `json:"password,omitempty" toml:"password,omitempty"`
	DbIndex      int      `json:"db,omitempty" toml:"db,omitempty"`
	DialTimeout  int      `json:"dial_timeout,omitempty" toml:"dial_timeout,omitempty"`
	ReadTimeout  int      `json:"read_timeout,omitempty" toml:"read_timeout,omitempty"`
	WriteTimeout int      `json:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	PoolSize     int      `json:"pool_size,omitempty" toml:"pool_size,omitempty"`
}

//New 构建redis消息队列配置
func New(addrs string, opts ...Option) *Redis {
	r := &Redis{
		Addrs:        types.Split(addrs, ","),
		DbIndex:      0,
		DialTimeout:  3,
		ReadTimeout:  3,
		WriteTimeout: 3,
		PoolSize:     10,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *Redis {

	org := New("", WithRaw(raw))

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}

//GetConf GetConf
func GetConf(varConf conf.IVarConf, name string) (redis *Redis, err error) {
	js, err := varConf.GetConf("redis", name)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("未配置：/var/redis/%s", name)
	}
	if err != nil {
		return nil, err
	}
	return NewByRaw(string(js.GetRaw())), nil
}
