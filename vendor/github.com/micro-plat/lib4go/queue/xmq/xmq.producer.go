package xmq

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/queue"
)

//XMQProducer Producer
type XMQProducer struct {
	conf        *Conf
	conn        net.Conn
	queues      cmap.ConcurrentMap
	connecting  bool
	isConnected bool
	closeCh     chan struct{}
	done        bool
	lk          sync.Mutex
	writeLock   sync.Mutex
	header      []string
	lastWrite   time.Time
	*mq.OptionConf
}

//New 创建新的producer
func New(address []string, c string) (producer *XMQProducer, err error) {
	producer = &XMQProducer{}
	producer.OptionConf = &mq.OptionConf{}
	producer.closeCh = make(chan struct{})
	producer.conf, err = NewConf(c)
	if err != nil {
		return nil, err
	}
	if producer.Logger == nil {
		producer.Logger = logger.GetSession("xmq.producer", logger.CreateSession())
	}
	producer.header = make([]string, 0, 4)
	return producer, producer.Connect()
}

//Connect  循环连接服务器
func (producer *XMQProducer) Connect() error {
	err := producer.connectOnce()
	if err != nil {
		producer.Logger.Error(err)
	}
	go func() {
	START:
		for {
			select {
			case <-producer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				if producer.isConnected {
					if time.Since(producer.lastWrite).Seconds() > 3 {
						message, err := NewXMQHeartBit().MakeMessage()
						if err != nil {
							producer.Logger.Error(err)
							continue
						}
						err = producer.writeMessage(message)
						if err == nil {
							continue
						}
						producer.Logger.Error(err)
						err = producer.connectOnce()
						if err != nil {
							producer.Logger.Error(err)
						}
					}
					continue
				}
				err = producer.connectOnce()
				if err != nil {
					producer.Logger.Error(err)
				}

			}
		}
	}()
	return nil
}

func (producer *XMQProducer) writeMessage(msg string) error {
	if !producer.isConnected {
		return fmt.Errorf("未连接到服务器")
	}
	producer.writeLock.Lock()
	producer.lastWrite = time.Now()
	result, err := encoding.ConvertBytes([]byte(msg), "gbk")
	if err != nil {
		return err
	}
	_, err = producer.conn.Write(result)
	producer.lastWrite = time.Now()
	producer.writeLock.Unlock()
	return err
}

func (producer *XMQProducer) disconnect() {
	producer.isConnected = false
	if producer.conn == nil {
		return
	}
	producer.conn.Close()
	return
}

//reconnect 自动重连
func (producer *XMQProducer) reconnect() {
	producer.conn.Close()
	producer.isConnected = false
	err := producer.Connect()
	if err != nil {
		producer.Logger.Errorf("连接到MQ服务器失败:%v", err)
	}
}

//ConnectOnce 连接到服务器
func (producer *XMQProducer) connectOnce() (err error) {
	if producer.connecting {
		return nil
	}
	producer.lk.Lock()
	defer producer.lk.Unlock()
	if producer.connecting {
		return nil
	}
	producer.connecting = true
	defer func() {
		producer.connecting = false
	}()
	producer.isConnected = false
	producer.conn, err = net.DialTimeout("tcp", producer.conf.Address, time.Second*2)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.isConnected = true
	producer.lastWrite = time.Now()
	return nil
}
func (producer *XMQProducer) Push(queue string, msg string) error {
	return producer.Send(queue, msg, time.Hour*1000)
}

func (producer *XMQProducer) Pop(key string) (string, error) {
	return "", fmt.Errorf("not support")
}
func (producer *XMQProducer) Count(key string) (int64, error) {
	return 0, nil
}

//Send 发送消息
func (producer *XMQProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	if !producer.connecting && producer.Retry {
		return fmt.Errorf("producer无法连接到MQ服务器:%s", producer.conf.Address)
	}
	message := NewXMQMessage(queue, msg, int(timeout/time.Second))
	if producer.OptionConf.Key != "" {
		message.signKey = producer.OptionConf.Key
	}
	smessage, err := message.MakeMessage()
	if err != nil {
		return
	}
	pmsg := &mq.ProcuderMessage{Queue: queue, Data: smessage, Timeout: timeout}
	return producer.writeMessage(pmsg.Data)
}

//Close 关闭当前连接
func (producer *XMQProducer) Close() error {
	producer.done = true
	close(producer.closeCh)
	return nil
}

type xmqResolver struct {
}

func (s *xmqResolver) Resolve(address []string, conf string) (queue.IQueue, error) {
	return New(address, conf)
}

func init() {
	queue.Register("xmq", &xmqResolver{})
}
