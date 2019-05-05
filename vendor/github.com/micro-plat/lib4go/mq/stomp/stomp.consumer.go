package stomp

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//StompConsumer Consumer
type StompConsumer struct {
	address    string
	conn       *stompngo.Connection
	cache      cmap.ConcurrentMap
	queues     cmap.ConcurrentMap
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	*mq.OptionConf
}

//NewStompConsumer 创建新的Consumer
func NewStompConsumer(address string, opts ...mq.Option) (consumer *StompConsumer, err error) {
	consumer = &StompConsumer{address: address}
	consumer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mq.consumer", logger.CreateSession())}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	consumer.cache = cmap.New(2)
	for _, opt := range opts {
		opt(consumer.OptionConf)
	}
	if strings.EqualFold(consumer.OptionConf.Version, "") {
		consumer.OptionConf.Version = "1.1"
	}
	if strings.EqualFold(consumer.OptionConf.Persistent, "") {
		consumer.OptionConf.Persistent = "true"
	}
	if strings.EqualFold(consumer.OptionConf.Ack, "") {
		consumer.OptionConf.Ack = "client-individual"
	}
	consumer.header = stompngo.Headers{"accept-version", consumer.OptionConf.Version}
	return
}

//Connect  循环连接服务器
func (consumer *StompConsumer) Connect() error {
	err := consumer.ConnectOnce()
	if err == nil {
		return nil
	}
	consumer.Logger.Error(err)
	go func() {
	START:
		for {
			select {
			case <-consumer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				err = consumer.ConnectOnce()
				if err == nil {
					break START
				}
				consumer.Logger.Error(err)
			}
		}
	}()
	return nil
}

//ConnectOnce 连接到服务器
func (consumer *StompConsumer) ConnectOnce() (err error) {
	if consumer.connecting {
		return nil
	}
	consumer.lk.Lock()
	defer consumer.lk.Unlock()
	if consumer.connecting {
		return nil
	}
	consumer.connecting = true
	defer func() {
		consumer.connecting = false
	}()
	con, err := net.Dial("tcp", consumer.address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	consumer.conn, err = stompngo.Connect(con, consumer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
	}

	//连接成功后开始订阅消息
	consumer.cache.IterCb(func(key string, value interface{}) bool {
		go func() {
			err = consumer.consume(key, value.(func(mq.IMessage)))
			if err != nil {
				consumer.Logger.Errorf("consume失败：%v", err)
			}
		}()
		return true
	})

	return nil
}

//Consume 订阅消息
func (consumer *StompConsumer) Consume(queue string, concurrency int, callback func(mq.IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}
	b, _ := consumer.cache.SetIfAbsent(queue, callback)
	if !b {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	return nil
}

//Consume 注册消费信息
func (consumer *StompConsumer) consume(queue string, callback func(mq.IMessage)) (err error) {
	success, ch, err := consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		header := stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.Ack}
		consumer.conn.SetSubChanCap(10)
		msgChan, err := consumer.conn.Subscribe(header)
		if err != nil {
			return
		}
		chans := &consumerChan{}
		chans.msgChan = msgChan
		chans.unconsumeCh = make(chan struct{})
		return chans, nil
	}, queue)
	if err != nil {
		return err
	}
	if !success {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	msgChan := ch.(*consumerChan)
START:
	for {
		select {
		case <-consumer.closeCh:
			break START
		case <-msgChan.unconsumeCh:
			break START
		case msg, ok := <-msgChan.msgChan:
			if !ok {
				break START
			}
			message := NewStompMessage(consumer, &msg.Message)
			if message.Has() {
				go callback(message)
			} else {
				consumer.reconnect(queue)
				break START
			}
		}
	}
	return
}
func (consumer *StompConsumer) reconnect(queue string) {
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.conn.Disconnect(stompngo.Headers{})
	consumer.Connect()
}

//UnConsume 取消注册消费
func (consumer *StompConsumer) UnConsume(queue string) {
	if consumer.conn == nil {
		return
	}
	header := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", "queue", queue), "ack", consumer.Ack}
	consumer.conn.Unsubscribe(header)
	if v, b := consumer.queues.Get(queue); b {
		ch := v.(*consumerChan)
		close(ch.unconsumeCh)
	}
	consumer.queues.Remove(queue)
	consumer.cache.Remove(queue)
}

//Close 关闭当前连接
func (consumer *StompConsumer) Close() {

	if consumer.conn == nil {
		return
	}
	consumer.once.Do(func() {
		close(consumer.closeCh)
	})

	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(*consumerChan)
		close(ch.unconsumeCh)
		return true
	})
	consumer.cache.Clear()
	go func() {
		defer recover()
		time.Sleep(time.Millisecond * 100)
		consumer.conn.Disconnect(stompngo.Headers{})
	}()

}

type stompConsumerResolver struct {
}

func (s *stompConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewStompConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("stomp", &stompConsumerResolver{})
}
