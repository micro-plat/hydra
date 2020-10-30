package lmq

import (
	"errors"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

//Consumer 基于本地channel的Consumer
type Consumer struct {
	queues  cmap.ConcurrentMap
	closeCh chan struct{}
	done    bool
	once    sync.Once
}

//newConsumer 创建新的Consumer
func newConsumer() (consumer *Consumer, err error) {
	consumer = &Consumer{
		queues:  cmap.New(4),
		closeCh: make(chan struct{})}
	return consumer, nil
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	return nil
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
		nconcurrency := types.GetMax(concurrency, 10)
		msgChan := make(chan *Message, nconcurrency)
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
			currQueue := GetOrAddQueue(queue)
		START:
			for {
				select {
				case <-consumer.closeCh:
					break START
				case <-unconsumeCh:
					break START
				case msg := <-currQueue:
					message := newMessage(msg)
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
func (consumer *Consumer) UnConsume(queue string) {
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
}

type consumerResolver struct {
}

func (s *consumerResolver) Resolve(confRaw string) (mq.IMQC, error) {
	return newConsumer()
}
func init() {
	mq.RegisterConsumer("lmq", &consumerResolver{})
}
