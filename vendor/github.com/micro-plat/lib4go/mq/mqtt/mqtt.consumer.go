package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	xnet "net"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/utility"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"github.com/zkfy/stompngo"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//Consumer Consumer
type Consumer struct {
	address    string
	client     *client.Client
	queues     cmap.ConcurrentMap
	subChan    chan string
	connecting bool
	closeCh    chan struct{}
	connCh     chan int
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
	consumer.OptionConf = &mq.OptionConf{QueueCount: 250, Logger: logger.GetSession("mqtt.consumer", logger.CreateSession())}
	consumer.closeCh = make(chan struct{})
	consumer.connCh = make(chan int, 1)
	consumer.queues = cmap.New(2)
	for _, opt := range opts {
		opt(consumer.OptionConf)
	}
	consumer.subChan = make(chan string, consumer.OptionConf.QueueCount)
	consumer.conf, err = NewConf(consumer.Raw)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	cc, _, err := consumer.connect()
	if err != nil {
		return err
	}
	consumer.client = cc
	go consumer.reconnect()
	go consumer.subscribe()
	return nil
}
func (consumer *Consumer) reconnect() {
	for {
		select {
		case <-time.After(time.Second * 3): //延迟重连
			select {
			case <-consumer.connCh:
				consumer.Logger.Debug("consumer与服务器断开连接，准备重连")
				func() {
					defer recover()
					consumer.client.Disconnect()
					consumer.client.Terminate()
				}()
				client, b, err := consumer.connect()
				if err != nil {
					consumer.Logger.Error("连接失败:", err)
				}
				if b {
					consumer.Logger.Info("consumer成功连接到服务器")
					consumer.client = client
					consumer.queues.IterCb(func(k string, v interface{}) bool {
						consumer.subChan <- k
						return true
					})
				}
			default:

			}
		}
	}
}

func (consumer *Consumer) connect() (*client.Client, bool, error) {
	consumer.lk.Lock()
	defer consumer.lk.Unlock()
	cert, err := consumer.getCert(consumer.conf)
	if err != nil {
		return nil, false, err
	}
	cc := client.New(&client.Options{
		ErrorHandler: func(err error) {
			select {
			case consumer.connCh <- 1: //发送重连消息
			default:
			}
		},
	})
	host, port, err := xnet.SplitHostPort(consumer.conf.Address)
	if err != nil {
		return nil, false, err
	}
	addrs, err := xnet.LookupHost(host)
	if err != nil {
		return nil, false, err
	}
	if err != nil {
		return nil, false, err
	}
	for _, addr := range addrs {
		if err := cc.Connect(&client.ConnectOptions{
			Network:   "tcp",
			Address:   addr + ":" + port,
			UserName:  []byte(consumer.conf.UserName),
			Password:  []byte(consumer.conf.Password),
			ClientID:  []byte(fmt.Sprintf("%s-%s", net.GetLocalIPAddress(), utility.GetGUID()[0:6])),
			TLSConfig: cert,
			KeepAlive: 3,
		}); err == nil {
			return cc, true, nil
		}
	}
	return nil, false, fmt.Errorf("连接失败:%v[%v](%s-%s/%s)", err, consumer.conf.Address, addrs, consumer.conf.UserName, consumer.conf.Password)
}

func (consumer *Consumer) getCert(conf *Conf) (*tls.Config, error) {
	if conf.CertPath == "" {
		return nil, nil
	}
	b, err := ioutil.ReadFile(conf.CertPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败:%s(%v)", conf.CertPath, err)
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(b); !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}
	return &tls.Config{
		RootCAs: roots,
	}, nil
}

//subscribe 循环接收，并放入指定的队列
// QoS0，最多一次送达。也就是发出去就fire掉，没有后面的事情了。
// QoS1，至少一次送达。发出去之后必须等待ack，没有ack，就要找时机重发
// QoS2，准确一次送达。消息id将拥有一个简单的生命周期。
func (consumer *Consumer) subscribe() {

START:
	for {
		select {
		case <-consumer.closeCh:
			break START
		case q := <-consumer.subChan:
			err := consumer.client.Subscribe(&client.SubscribeOptions{
				SubReqs: []*client.SubReq{
					&client.SubReq{
						TopicFilter: []byte(q),
						QoS:         mqtt.QoS0,
						Handler: func(topicName, message []byte) {
							nmsg := NewMessage()
							_, err := nmsg.Write(message)
							if err != nil {
								consumer.Logger.Error("写入消息失败:", string(message))
								return
							}
							if nq, b := consumer.queues.Get(string(topicName)); b {
								nQ := nq.(chan *Message)
								nQ <- nmsg
							}
						},
					},
				},
			})
			if err != nil {
				consumer.Logger.Error("消息订阅出错", err)
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
	err := consumer.client.Unsubscribe(&client.UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte(queue),
		},
	})
	if err != nil {
		consumer.Logger.Errorf("取消订单消息出错(queue:%s)err:%v", queue, err)
	}
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
	consumer.client.Terminate()
}

type ConsumerResolver struct {
}

func (s *ConsumerResolver) Resolve(address string, opts ...mq.Option) (mq.MQConsumer, error) {
	return NewConsumer(address, opts...)
}
func init() {
	mq.RegisterCosnumer("mqtt", &ConsumerResolver{})
}
