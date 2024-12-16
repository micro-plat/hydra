package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/vars/queue/kafka"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

// Producer memcache配置文件
type Producer struct {
	client   sarama.SyncProducer
	confOpts *kafka.Kafka
}

// NewProducerByRaw 根据配置文件创建一个redis连接
func NewProducerByRaw(cfg string) (m *Producer, err error) {
	return NewProducerByConfig(kafka.NewByRaw(cfg))
}

// NewProducerByConfig 根据配置文件创建一个redis连接
func NewProducerByConfig(confOpts *kafka.Kafka) (m *Producer, err error) {
	m = &Producer{confOpts: confOpts}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //ACK,发送完数据需要leader和follow都确认
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = time.Duration(types.GetInt(confOpts.WriteTimeout, 5)) * time.Second
	p, err := sarama.NewSyncProducer(confOpts.Addrs, config)
	if err != nil {
		return nil, err
	}
	m.client = p
	return m, nil
}

// Push 向存于 key 的列表的尾部插入所有指定的值
func (c *Producer) Push(key string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: key,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := c.client.SendMessage(msg)
	return err
}

// Pop 移除并且返回 key 对应的 list 的第一个元素。
func (c *Producer) Pop(key string) (string, error) {
	return "", fmt.Errorf("kafka不支持pop方法")
}

// Count 获取列表中的元素个数
func (c *Producer) Count(key string) (int64, error) {
	return 0, fmt.Errorf("kafka不支持Count方法")
}

// Close 释放资源
func (c *Producer) Close() error {
	return c.client.Close()
}

type producerResolver struct {
}

func (s *producerResolver) Resolve(confRaw string) (mq.IMQP, error) {
	return NewProducerByRaw(confRaw)
}

func init() {
	mq.RegisterProducer(global.ProtoKafka, &producerResolver{})
}
