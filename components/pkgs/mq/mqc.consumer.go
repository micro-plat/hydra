package mq

import (
	"fmt"
	"strings"
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
	Resolve(address string, opts ...Option) (IMQC, error)
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
func NewMQC(address string, opts ...Option) (IMQC, error) {
	proto, addrs, err := getNames(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := mqcResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("mqc: 未知的协议类型 %s", proto)
	}
	return resolver.Resolve(addrs[0], opts...)
}
func getNames(address string) (proto string, raddr []string, err error) {
	addr := strings.Split(address, "://")
	if len(addr) > 2 {
		return "", nil, fmt.Errorf("mqc: 消息队列的地址配置错误 %s", addr)
	}
	if len(addr[0]) == 0 {
		return "", nil, fmt.Errorf("mqc: 消息队列的地址配置有误 %s", addr)
	}
	proto = addr[0]
	if len(addr) > 1 {
		raddr = strings.Split(addr[1], ",")
	} else {
		raddr = append(raddr, "")
	}
	return
}
