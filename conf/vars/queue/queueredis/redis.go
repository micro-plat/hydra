package queueredis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/lib4go/types"
)

//Redis redis缓存配置
type Redis struct {
	*queue.Queue

	Addrs        []string `json:"addrs,omitempty" toml:"addrs,omitempty" label:"集群地址(|分割)"`
	Password     string   `json:"password,omitempty" toml:"password,omitempty"`
	DbIndex      int      `json:"db,omitempty" toml:"db,omitempty"`
	DialTimeout  int      `json:"dial_timeout,omitempty" toml:"dial_timeout,omitempty"`
	ReadTimeout  int      `json:"read_timeout,omitempty" toml:"read_timeout,omitempty"`
	WriteTimeout int      `json:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	PoolSize     int      `json:"pool_size,omitempty" toml:"pool_size,omitempty"`

	ConfigName string `json:"config_name,omitempty"  toml:"config_name,omitempty" valid:"ascii"`
}

//New 构建redis消息队列配置
func New(addrs string, opts ...Option) (org *Redis) {
	org = &Redis{
		Queue: &queue.Queue{Proto: "redis"},
		Addrs: types.Split(addrs, ","),
	}
	for _, opt := range opts {
		opt(org)
	}
	b, err := govalidator.ValidateStruct(org)
	if !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}
	if org.ConfigName == "" && len(org.Addrs) == 0 {
		panic(fmt.Errorf("redis配置数据有误:至少存在Addrs或ConfigName一种,%+v", org))
	}
	return org
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) (org *Redis) {
	return New("", WithRaw(raw))
}

//GetRaw 获取配置数据（json)
func (org *Redis) GetRaw() string {
	if org.ConfigName == "" {
		bytes, _ := json.Marshal(org)
		return string(bytes)
	}

	varConf, err := app.Cache.GetVarConf()
	if err != nil {
		panic(fmt.Errorf("app.Cache.GetVarConf:%w", err))
	}
	js, err := varConf.GetConf("redis", org.ConfigName)
	if errors.Is(err, conf.ErrNoSetting) {
		panic(fmt.Errorf("未配置：/var/redis/%s", org.ConfigName))
	}
	tmp := &Redis{}
	if err := json.Unmarshal(js.GetRaw(), tmp); err != nil {
		panic(fmt.Errorf("cacheredis.WithRaw:%w", err))
	}
	bytes, _ := json.Marshal(tmp)
	return string(bytes)
}
