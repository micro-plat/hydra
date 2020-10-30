package xmq

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/queue"
)

//XMQ XMQ缓存配置
type XMQ struct {
	*queue.Queue
	SignKey string `json:"sign_key"  toml:"sign_key" valid:"required"`
	Address string `json:"address"  toml:"config_name" valid:"dialstring,required"`
}

//New 构建XMQ消息队列配置
func New(address string, opts ...Option) *XMQ {
	r := &XMQ{
		Address: address,
		Queue:   &queue.Queue{Proto: "xmq"},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *XMQ {
	org := New("", WithRaw(raw))
	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
