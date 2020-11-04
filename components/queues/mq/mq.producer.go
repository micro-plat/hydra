package mq

import (
	"errors"
	"fmt"
)

var Nil = errors.New("nil")

//IMQP 消息生产
type IMQP interface {
	Push(key string, value string) error
	Pop(key string) (string, error)
	Count(key string) (int64, error)
	Close() error
}

//imqpResover 定义配置文件转换方法
type imqpResover interface {
	Resolve(confRaw string) (IMQP, error)
}

var mqpResolvers = make(map[string]imqpResover)

//RegisterProducer 注册配置文件适配器
func RegisterProducer(proto string, resolver imqpResover) {
	if _, ok := mqpResolvers[proto]; ok {
		panic("mqp: 不能重复注册producer " + proto)
	}
	mqpResolvers[proto] = resolver
}

//NewMQP 根据适配器名称及参数返回配置处理器
func NewMQP(proto string, confRaw string) (IMQP, error) {
	resolver, ok := mqpResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("mqp: 不支持的消息协议 %s", proto)
	}
	return resolver.Resolve(confRaw)
}
