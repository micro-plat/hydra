package redis

import (
	"errors"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/redis"
)

//RedisProducer Producer
type RedisProducer struct {
	address   string
	client    *redis.Client
	backupMsg chan *mq.ProcuderMessage
	closeCh   chan struct{}
	done      bool
	*mq.OptionConf
}

//NewRedisProducer 创建新的producer
func NewRedisProducer(address string, opts ...mq.Option) (producer *RedisProducer, err error) {
	producer = &RedisProducer{address: address}
	producer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mq.redis", logger.CreateSession())}
	producer.closeCh = make(chan struct{})
	for _, opt := range opts {
		opt(producer.OptionConf)
	}
	return
}

//Connect  循环连接服务器
func (producer *RedisProducer) Connect() (err error) {
	producer.client, err = redis.NewClientByJSON(producer.Raw)
	return
}

//GetBackupMessage 获取备份数据
func (producer *RedisProducer) GetBackupMessage() chan *mq.ProcuderMessage {
	return producer.backupMsg
}

//Send 发送消息
func (producer *RedisProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	_, err = producer.client.RPush(queue, msg).Result()
	return
}

//Close 关闭当前连接
func (producer *RedisProducer) Close() {
	producer.done = true
	close(producer.closeCh)
	close(producer.backupMsg)
	producer.client.Close()
}

type redisProducerResolver struct {
}

func (s *redisProducerResolver) Resolve(address string, opts ...mq.Option) (mq.MQProducer, error) {
	return NewRedisProducer(address, opts...)
}
func init() {
	mq.RegisterProducer("redis", &redisProducerResolver{})
}
