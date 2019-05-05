package mq

import (
	"fmt"
	"strings"
)

type MQConsumer interface {
	Connect() error
	Consume(queue string, concurrency int, callback func(IMessage)) (err error)
	UnConsume(queue string)
	Close()
}

//MQConsumerResover 定义配置文件转换方法
type MQConsumerResover interface {
	Resolve(address string, opts ...Option) (MQConsumer, error)
}

var mqConsumerResolvers = make(map[string]MQConsumerResover)

//RegisterCosnumer 注册配置文件适配器
func RegisterCosnumer(adapter string, resolver MQConsumerResover) {
	if resolver == nil {
		panic("mq: Register adapter is nil")
	}
	if _, ok := mqConsumerResolvers[adapter]; ok {
		panic("mq: Register called twice for adapter " + adapter)
	}
	mqConsumerResolvers[adapter] = resolver
}

//NewMQConsumer 根据适配器名称及参数返回配置处理器
func NewMQConsumer(address string, opts ...Option) (MQConsumer, error) {
	proto, addrs, err := getMQNames(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := mqConsumerResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("mq.consumer: unknown adapter name %q (forgotten import?)", proto)
	}
	return resolver.Resolve(addrs[0], opts...)
}
func getMQNames(address string) (proto string, raddr []string, err error) {
	addr := strings.Split(address, "://")
	if len(addr) > 2 {
		return "", nil, fmt.Errorf("MQ地址配置错误%s，格式:stomp://192.168.0.1:61613", addr)
	}
	if len(addr[0]) == 0 {
		return "", nil, fmt.Errorf("MQ地址配置错误%s，格式:stomp://192.168.0.1:61613", addr)
	}
	proto = addr[0]
	if len(addr) > 1 {
		raddr = strings.Split(addr[1], ",")
	} else {
		raddr = append(raddr, "")
	}
	return
}
