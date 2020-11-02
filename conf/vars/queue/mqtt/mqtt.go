package mqtt

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/hydra/conf/vars/queue"
)

//MQTT MQTT对列配置
type MQTT struct {
	*queue.Queue
	Address     string `json:"address"  toml:"address" valid:"required"`
	DialTimeout int64  `json:"dial_timeout" toml:"dial_timeout" valid:"required"`
	UserName    string `json:"userName,omitempty"  toml:"userName,omitempty" `
	Password    string `json:"password,omitempty"  toml:"password,omitempty" `
	Cert        string `json:"cert,omitempty"  toml:"cert,omitempty"`
}

//New 构建mqtt配置
func New(address string, opts ...Option) *MQTT {
	r := &MQTT{
		Address: address,
		Queue:   &queue.Queue{Proto: "mqtt"},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//NewByRaw 通过json原串初始化
func NewByRaw(raw string) *MQTT {
	org := New("", WithRaw(raw))
	if b, err := govalidator.ValidateStruct(org); !b {
		panic(fmt.Errorf("redis配置数据有误:%v %+v", err, org))
	}

	return org
}
