package kafka

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/IBM/sarama"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/vars/queue/kafka"
	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

// Consumer Consumer
type Consumer struct {
	client   sarama.ConsumerGroup
	queues   cmap.ConcurrentMap
	closeCh  chan struct{}
	done     bool
	once     sync.Once
	log      logger.ILogger
	ConfOpts *kafka.Kafka
}

// NewConsumerByRaw 创建新的Consumer
func NewConsumerByRaw(cfg string) (consumer *Consumer, err error) {
	return NewConsumerByConfig(kafka.NewByRaw(cfg))
}

// NewConsumerByConfig 创建新的Consumer
func NewConsumerByConfig(cfg *kafka.Kafka) (consumer *Consumer, err error) {
	consumer = &Consumer{log: logger.GetSession("mq.kafka", logger.CreateSession())}
	consumer.ConfOpts = cfg

	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	return
}

// Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = consumer.ConfOpts.Offset
	if consumer.ConfOpts.Group == "" {
		return fmt.Errorf("kafka消费者必须配置group id节点")
	}

	// init consumer
	c, err := sarama.NewConsumerGroup(consumer.ConfOpts.Addrs,
		consumer.ConfOpts.Group, config)
	if err != nil {
		return err
	}
	consumer.client = c
	return
}

// Consume 注册消费信息
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
		nconcurrency := types.GetMax(concurrency, 1)
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			msgconsumer := NewConsumeHandler(consumer.log, callback, nconcurrency)
			defer cancel()
		START:
			for {
				select {
				case <-consumer.closeCh:
					break START
				case <-unconsumeCh:
					break START
				case <-time.After(time.Millisecond * (time.Duration((1000 / nconcurrency / 2)) + 1)):
					if consumer.client != nil && !consumer.done {
						if err := consumer.client.Consume(ctx, []string{queue}, msgconsumer); err != nil {
							if !consumer.done {
								consumer.log.Errorf("从kafka中获取消息失败:%v", err)
							}
							continue
						}
						if ctx.Err() != nil {
							consumer.log.Errorf("构建kafka consumer失败:%v", err)
							continue
						}
					}

				}
			}
		}()
		return unconsumeCh, nil
	}, queue)

	return
}

// UnConsume 取消注册消费
func (consumer *Consumer) UnConsume(queue string) {
	if consumer.client == nil {
		return
	}
	if c, ok := consumer.queues.Get(queue); ok {
		close(c.(chan struct{}))
	}
	consumer.queues.Remove(queue)
}

// Close 关闭当前连接
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
	return NewConsumerByRaw(confRaw)
}
func init() {
	mq.RegisterConsumer(global.ProtoKafka, &cresolver{})
}
