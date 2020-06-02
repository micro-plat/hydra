package mqtt

import "github.com/micro-plat/hydra/conf/vars/queue"

//MQTT MQTT对列配置
type MQTT struct {
	*queue.Queue
	Address  string `json:"address,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Cert     string `json:"cert,omitempty"`
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
