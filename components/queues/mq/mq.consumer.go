package mq

import (
	"fmt"
)

//IMQCMessage  队列消息
type IMQCMessage interface {
	Ack() error
	Nack() error
	GetMessage() string
}

//IMQC consumer接口
type IMQC interface {
	Connect() error
	Consume(queue string, concurrency int, callback func(IMQCMessage)) (err error)
	UnConsume(queue string)
	Close()
}

//mqcResover 定义消息消费解析器
type mqcResover interface {
	Resolve(confRaw string) (IMQC, error)
}

var mqcResolvers = make(map[string]mqcResover)

//RegisterConsumer 注册消息消费
func RegisterConsumer(adapter string, resolver mqcResover) {
	if _, ok := mqcResolvers[adapter]; ok {
		panic("mqc: 不能重复注册mqc " + adapter)
	}
	mqcResolvers[adapter] = resolver
}

//NewMQC 根据适配器名称及参数返回配置处理器
func NewMQC(proto string, confRaw string) (IMQC, error) {

	resolver, ok := mqcResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("mqc: 未知的协议类型 %s", proto)
	}
	return resolver.Resolve(confRaw)
}
