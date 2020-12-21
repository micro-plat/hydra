package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/global"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
	"github.com/zkfy/stompngo"

	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
)

type consumerChan struct {
	msgChan     <-chan stompngo.MessageData
	unconsumeCh chan struct{}
}

//Consumer Consumer
type Consumer struct {
	client     mqtt.Client
	queues     cmap.ConcurrentMap
	subChan    chan string
	connecting bool
	closeCh    chan struct{}
	uid        string
	connCh     chan int
	done       bool
	lk         sync.Mutex
	header     []string
	once       sync.Once
	clientOnce sync.Once
	log        logger.ILogger
	confOpts   *queuemqtt.MQTT
}

//NewConsumer 创建新的Consumer
func NewConsumerByConfig(confOpts *queuemqtt.MQTT) (consumer *Consumer, err error) {
	consumer = &Consumer{uid: utility.GetGUID()[0:6], log: logger.GetSession("mqtt", logger.CreateSession())}
	consumer.confOpts = confOpts
	consumer.closeCh = make(chan struct{})
	consumer.connCh = make(chan int, 1)
	consumer.queues = cmap.New(2)

	consumer.subChan = make(chan string, 512)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

//Connect  连接服务器
func (consumer *Consumer) Connect() (err error) {
	cc, _, err := consumer.connect()
	if err != nil {
		consumer.log.Fatalf("创建eclipse.consumer连接失败，%v", err)
		return err
	}
	consumer.log.Info("创建eclipse.consumer连接成功")
	consumer.client = cc

	go consumer.reconnect()
	go consumer.subscribe()
	return nil
}
func (consumer *Consumer) reconnect() {
	for {
		select {
		case <-time.After(time.Second * 3): //延迟重连
			if !consumer.done && !consumer.client.IsConnected() {
				consumer.client.Disconnect(250)
				select {
				case consumer.connCh <- 1: //发送重连消息
				default:
				}
			}
			select {
			case <-consumer.connCh:
				consumer.log.Debugf("eclipse.consumer连接到服务器:%s", consumer.getAddr())
				client, b, err := consumer.connect()
				if err != nil {
					consumer.client.Disconnect(250)
					consumer.connCh <- 1
					consumer.log.Error("eclipse.consumer连接失败:", err)
				}
				if b {
					consumer.log.Info("eclipse.consumer成功连接到服务器")
					consumer.client = client

					//重新订阅所有消息
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

func (consumer *Consumer) connect() (mqtt.Client, bool, error) {
	consumer.lk.Lock()
	defer consumer.lk.Unlock()
	if consumer.done || consumer.client != nil && consumer.client.IsConnected() {
		return consumer.client, false, nil
	}
	cert, err := consumer.getCert(consumer.confOpts)
	if err != nil {
		return nil, false, err
	}

	opts := mqtt.NewClientOptions().AddBroker(consumer.getAddr())
	opts.SetUsername(consumer.confOpts.UserName)
	opts.SetPassword(consumer.confOpts.Password)
	opts.SetClientID(fmt.Sprintf("%s-%s", global.LocalIP(), consumer.uid))
	opts.SetTLSConfig(cert)
	opts.SetKeepAlive(10)
	opts.SetAutoReconnect(false)
	cc := mqtt.NewClient(opts)
	if token := cc.Connect(); token.Wait() && token.Error() != nil {
		return nil, false, token.Error()
	}
	return cc, true, nil
}

func (consumer *Consumer) getCert(conf *queuemqtt.MQTT) (*tls.Config, error) {
	if conf.Cert == "" {
		return nil, nil
	}
	b, err := ioutil.ReadFile(conf.Cert)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败:%s(%v)", conf.Cert, err)
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(b); !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}
	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: roots,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
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
			tk := consumer.client.Subscribe(q, 0, func(client mqtt.Client, msg mqtt.Message) {
				nmsg := NewMessage()
				_, err := nmsg.Write(msg.Payload())
				if err != nil {
					consumer.log.Error("写入消息失败:", string(msg.Payload()))
					return
				}
				if nq, b := consumer.queues.Get(string(msg.Topic())); b {
					nQ := nq.(chan *Message)
					nQ <- nmsg
				}
			})
			if tk.Wait() && tk.Error() != nil {
				consumer.log.Error("订阅消息出错", tk.Error())
			}
		}
	}
}

//Consume 注册消费信息
func (consumer *Consumer) Consume(queue string, concurrency int, callback func(mq.IMQCMessage)) (err error) {
	if strings.EqualFold(queue, "") {
		return errors.New("队列名字不能为空")
	}
	if callback == nil {
		return errors.New("回调函数不能为nil")
	}
	concurrency = types.GetMax(concurrency, 10)
	_, _, err = consumer.queues.SetIfAbsentCb(queue, func(input ...interface{}) (c interface{}, err error) {
		queue := input[0].(string)
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
						callback(message)
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
	consumer.queues.Remove(queue)
	err := consumer.client.Unsubscribe(queue)
	if err != nil {
		consumer.log.Errorf("取消订阅出错(queue:%s)err:%v", queue, err)
	}
}
func (consumer *Consumer) getAddr() string {
	if strings.Contains(consumer.confOpts.Address, "://") {
		return consumer.confOpts.Address
	}
	return fmt.Sprintf("ssl://%s", consumer.confOpts.Address)
}

//Close 关闭当前连接
func (consumer *Consumer) Close() {
	consumer.done = true
	consumer.once.Do(func() {
		close(consumer.closeCh)
	})
	consumer.queues.RemoveIterCb(func(key string, value interface{}) bool {
		ch := value.(chan *Message)
		close(ch)
		return true
	})
	if consumer.client != nil {
		consumer.client.Disconnect(250)
	}

}

type resolver struct {
}

func (s *resolver) Resolve(confRaw string) (mq.IMQC, error) {
	confOpts := queuemqtt.NewByRaw(confRaw)
	return NewConsumerByConfig(confOpts)
}
func init() {
	mq.RegisterConsumer("mqtt", &resolver{})
}
