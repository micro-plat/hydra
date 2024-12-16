package kafka

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

// TypeNodeName 分类节点名
const TypeNodeName = "kafka"

// Kafka Kafka缓存配置
type Kafka struct {
	security.ConfEncrypt
	*queue.Queue
	Addrs        []string `json:"addrs,omitempty" toml:"addrs,omitempty" valid:"required" label:"集群地址(|分割)"`
	WriteTimeout int      `json:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	Offset       int64    `json:"offset,omitempty" toml:"offset,omitempty"`
	Group        string   `json:"group,omitempty" toml:"group,omitempty"`
}

// New 构建Kafka消息队列配置
func New(addrs string, opts ...Option) *Kafka {
	r := &Kafka{
		Queue:  &queue.Queue{Proto: global.ProtoKafka},
		Addrs:  types.Split(addrs, ","),
		Offset: -1,
	}
	for _, opt := range opts {
		opt(r)
	}
	if r.Offset >= -1 {
		r.Offset = -1
	}
	if r.Offset <= -2 {
		r.Offset = -2
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = 30
	}
	return r
}

// NewByRaw 通过json原串初始化
func NewByRaw(raw string) *Kafka {

	org := New("", WithRaw(raw))

	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("kafka配置数据有误:%v %+v", err, org))
	}

	return org
}

// GetConf GetConf
func GetConf(varConf conf.IVarConf, name string) (Kafka *Kafka, err error) {
	js, err := varConf.GetConf("kafka", name)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("未配置：/var/kafka/%s", name)
	}
	if err != nil {
		return nil, err
	}
	return NewByRaw(string(js.GetRaw())), nil
}
