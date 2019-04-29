package mq

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ProcuderMessage struct {
	Headers   []string
	Queue     string
	Data      string
	SendTimes int32
	Timeout   time.Duration
}

func (p *ProcuderMessage) AddSendTimes() {
	atomic.AddInt32(&p.SendTimes, 1)
}

type MQProducer interface {
	Connect() error
	GetBackupMessage() chan *ProcuderMessage
	Send(queue string, msg string, timeout time.Duration) (err error)
	Close()
}

//MQConsumerResover 定义配置文件转换方法
type MQProducerResover interface {
	Resolve(address string, opts ...Option) (MQProducer, error)
}

var mqProducerResolvers = make(map[string]MQProducerResover)

//RegisterProducer 注册配置文件适配器
func RegisterProducer(proto string, resolver MQProducerResover) {
	if resolver == nil {
		panic("mq: Register adapter is nil")
	}
	if _, ok := mqProducerResolvers[proto]; ok {
		panic("mq: Register called twice for adapter " + proto)
	}
	mqProducerResolvers[proto] = resolver
}

//NewMQProducer 根据适配器名称及参数返回配置处理器
func NewMQProducer(address string, opts ...Option) (MQProducer, error) {
	proto, addrs, err := getMQNames(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := mqProducerResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("mq.producer: unknown adapter name %q (forgotten import?)", proto)
	}
	return resolver.Resolve(addrs[0], opts...)
}
