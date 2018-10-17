package redis

import (
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/redis"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//RedisConsumer Consumer
type RedisConsumer struct {
	address    string
	client     *redis.Client
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	*mq.OptionConf
}

//NewRedisConsumer 创建新的Consumer
func NewRedisConsumer(address string, opts ...mq.Option) (consumer *RedisConsumer, err error) {
	consumer = &RedisConsumer{address: address}
	consumer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mq.redis", logger.CreateSession())}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	for _, opt := range opts {
		opt(consumer.OptionConf)
	}
	return
}

//Connect  连接服务器
func (consumer *RedisConsumer) Connect() (err error) {
	consumer.client, err = redis.NewClientByJSON(consumer.Raw)
	return
}

//Consume 注册消费信息
func (consumer *RedisConsumer) Consume(queue string, concurrency int, callback func(mq.IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}

	_, _, err = consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		unconsumeCh := make(chan struct{}, 1)
		if concurrency <= 0 {
			concurrency = 10
		}
		msgChan := make(chan *RedisMessage, concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
			START:
				for {
					select {
					case message, ok := <-msgChan:
						if !ok {
							break START
						}
						go callback(message)
					}
				}
			}()
		}
		go func() {
		START:
			for {
				select {
				case <-consumer.closeCh:
					break START
				case <-unconsumeCh:
					break START
				case <-time.After(time.Millisecond * 50):
					if consumer.client != nil && !consumer.done {
						cmd := consumer.client.BLPop(time.Second, queue)
						message := NewRedisMessage(cmd)
						if message.Has() {
							msgChan <- message

						}
					}

				}
			}
			close(msgChan)
		}()
		return unconsumeCh, nil
	}, queue)
	return
}

//UnConsume 取消注册消费
func (consumer *RedisConsumer) UnConsume(queue string) {
	if consumer.client == nil {
		return
	}
	if c, ok := consumer.queues.Get(queue); ok {
		close(c.(chan struct{}))
	}
	consumer.queues.Remove(queue)
}

//Close 关闭当前连接
func (consumer *RedisConsumer) Close() {

	consumer.once.Do(func() {
		close(consumer.closeCh)
	})

	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(chan struct{})
		close(ch)
		return true
	})
	if consumer.client == nil {
		return
	}
	consumer.client.Close()
}

type redisConsumerResolver struct {
}

func (s *redisConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewRedisConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("redis", &redisConsumerResolver{})
}
