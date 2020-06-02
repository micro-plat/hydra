package redis

import (
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/micro-plat/hydra/components/pkg/mq"
	"github.com/micro-plat/hydra/components/pkg/redis"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//Consumer Consumer
type Consumer struct {
	address    string
	client     *redis.Client
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	log        logger.ILogger
	*mq.ConfOpt
}

//NewConsumer 创建新的Consumer
func NewConsumer(address string, opts ...mq.Option) (consumer *Consumer, err error) {
	consumer = &Consumer{address: address, log: logger.GetSession("mq.redis", logger.CreateSession())}
	consumer.ConfOpt = &mq.ConfOpt{}
	for _, opt := range opts {
		opt(consumer.ConfOpt)
	}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	return
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	consumer.client, err = redis.New(redis.WithRaw(consumer.ConfOpt.Raw))
	return
}

//Consume 注册消费信息
func (consumer *Consumer) Consume(queue string, concurrency int, callback func(mq.IMQCMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}

	_, _, err = consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		unconsumeCh := make(chan struct{}, 1)
		nconcurrency := concurrency
		if concurrency <= 0 {
			nconcurrency = 10
		}
		msgChan := make(chan *RedisMessage, nconcurrency)
		for i := 0; i < nconcurrency; i++ {
			go func() {
			START:
				for {
					select {
					case message, ok := <-msgChan:
						if !ok {
							break START
						}
						if concurrency == 0 {
							go callback(message)
						} else {
							callback(message)
						}
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
				case <-time.After(time.Millisecond * (time.Duration((1000 / nconcurrency / 2)) + 1)):
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
func (consumer *Consumer) UnConsume(queue string) {
	if consumer.client == nil {
		return
	}
	if c, ok := consumer.queues.Get(queue); ok {
		close(c.(chan struct{}))
	}
	consumer.queues.Remove(queue)
}

//Close 关闭当前连接
func (consumer *Consumer) Close() {

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

type cresolver struct {
}

func (s *cresolver) Resolve(address string, opts ...mq.Option) (mq.IMQC, error) {
	return NewConsumer(address, opts...)
}
func init() {
	mq.RegisterConsumer("redis", &cresolver{})
}