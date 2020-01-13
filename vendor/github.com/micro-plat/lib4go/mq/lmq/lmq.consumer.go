package lmq

import (
	"errors"
	"strings"
	"sync"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/queue/lmq"
	"github.com/micro-plat/lib4go/types"
)

//lmqConsumer 基于本地channel的Consumer
type lmqConsumer struct {
	queues  cmap.ConcurrentMap
	closeCh chan struct{}
	done    bool
	once    sync.Once
}

//newlmqConsumer 创建新的Consumer
func newlmqConsumer(address string, opts ...mq.Option) (consumer *lmqConsumer, err error) {
	consumer = &lmqConsumer{
		queues:  cmap.New(4),
		closeCh: make(chan struct{})}
	return consumer, nil
}

//Connect  连接服务器
func (consumer *lmqConsumer) Connect() (err error) {
	return nil
}

//Consume 注册消费信息
func (consumer *lmqConsumer) Consume(queue string, concurrency int, callback func(mq.IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}
	_, _, err = consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		unconsumeCh := make(chan struct{}, 1)
		nconcurrency := types.GetMax(concurrency, 10)
		msgChan := make(chan *LMQMessage, nconcurrency)
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
			currQueue := lmq.GetOrAddQueue(queue)
		START:
			for {
				select {
				case <-consumer.closeCh:
					break START
				case <-unconsumeCh:
					break START
				case msg := <-currQueue:
					message := NewLMQMessage(msg)
					if message.Has() {
						msgChan <- message
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
func (consumer *lmqConsumer) UnConsume(queue string) {
	if c, ok := consumer.queues.Get(queue); ok {
		close(c.(chan struct{}))
	}
	consumer.queues.Remove(queue)
}

//Close 关闭当前连接
func (consumer *lmqConsumer) Close() {
	consumer.once.Do(func() {
		close(consumer.closeCh)
	})

	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(chan struct{})
		close(ch)
		return true
	})
}

type lmqConsumerResolver struct {
}

func (s *lmqConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return newlmqConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("lmq", &lmqConsumerResolver{})
}
