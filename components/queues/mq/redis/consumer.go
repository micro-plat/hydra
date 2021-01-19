package redis

import (
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/micro-plat/hydra/components/pkgs/redis"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"

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
	ConfOpts   *varredis.Redis
}

//NewConsumerByRaw 创建新的Consumer
func NewConsumerByRaw(cfg string) (consumer *Consumer, err error) {
	return NewConsumerByConfig(varredis.NewByRaw(cfg))
}

//NewConsumerByConfig 创建新的Consumer
func NewConsumerByConfig(cfg *varredis.Redis) (consumer *Consumer, err error) {
	consumer = &Consumer{log: logger.GetSession("mq.redis", logger.CreateSession())}
	consumer.ConfOpts = cfg

	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	return
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	consumer.client, err = redis.NewByConfig(consumer.ConfOpts)
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
							//默认10个线程获取任务，开启新协程处理任务
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
						if err := cmd.Err(); err != nil {
							if !consumer.done && err.Error() != "redis: nil" {
								consumer.log.Error("从redis中获取消息失败:%w", err)
							}
							continue
						}
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
		consumer.done = true
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

func (s *cresolver) Resolve(confRaw string) (mq.IMQC, error) {
	return NewConsumerByRaw(queueredis.NewByRaw(confRaw).GetRaw())
}
func init() {
	mq.RegisterConsumer(Proto, &cresolver{})
}

//Proto redis
const Proto = "redis"
