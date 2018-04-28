package mqtt

import (
	"net"
	"strings"
	"sync"

	"errors"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	xnet "github.com/micro-plat/lib4go/net"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//Consumer Consumer
type Consumer struct {
	address    string
	client     *mqtt.ClientConn
	queues     cmap.ConcurrentMap
	subChan    chan string
	connecting bool
	closeCh    chan struct{}
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	*mq.OptionConf
	conf *Conf
}

//NewConsumer 创建新的Consumer
func NewConsumer(address string, opts ...mq.Option) (consumer *Consumer, err error) {
	consumer = &Consumer{address: address}
	consumer.OptionConf = &mq.OptionConf{Logger: logger.GetSession("mqtt", logger.CreateSession())}
	consumer.closeCh = make(chan struct{})
	consumer.queues = cmap.New(2)
	consumer.subChan = make(chan string, 3)
	for _, opt := range opts {
		opt(consumer.OptionConf)
	}
	consumer.conf, err = NewConf(consumer.Raw)
	return
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	conn, err := net.Dial("tcp", consumer.conf.Address)
	if err != nil {
		return err
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = consumer.conf.DumpData
	cc.ClientId = xnet.GetLocalIPAddress()
	if err = cc.Connect(consumer.conf.UserName, consumer.conf.Password); err != nil {
		return err
	}
	consumer.client = cc
	go consumer.recvMessage()
	return nil
}

//recvMessage 循环接收，并放入指定的队列
func (consumer *Consumer) recvMessage() {

START:
	for {
		select {
		case <-consumer.closeCh:
			break START
		case q := <-consumer.subChan:
			tq := make([]proto.TopicQos, 1)
			tq[0].Topic = q
			tq[0].Qos = proto.QosAtMostOnce
			consumer.client.Subscribe(tq)
		case msg := <-consumer.client.Incoming:
			nmsg := NewMessage()
			if err := msg.Payload.WritePayload(nmsg); err != nil {
				consumer.Logger.Error(err)
				continue
			}
			if nq, b := consumer.queues.Get(msg.TopicName); b {
				nQ := nq.(chan *Message)
				nQ <- nmsg
			}
		}
	}
}

//Consume 注册消费信息
func (consumer *Consumer) Consume(queue string, concurrency int, callback func(mq.IMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}

	_, _, err = consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
		if concurrency <= 0 {
			concurrency = 10
		}
		msgChan := make(chan *Message, concurrency)
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
		consumer.subChan <- queue
		return msgChan, nil
	}, queue)
	return
}

//UnConsume 取消注册消费
func (consumer *Consumer) UnConsume(queue string) {

}

//Close 关闭当前连接
func (consumer *Consumer) Close() {
	consumer.once.Do(func() {
		close(consumer.closeCh)
	})

	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(chan *Message)
		close(ch)
		return true
	})
	if consumer.client == nil {
		return
	}
	consumer.client.Disconnect()
}

type ConsumerResolver struct {
}

func (s *ConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("mqtt", &ConsumerResolver{})
}
