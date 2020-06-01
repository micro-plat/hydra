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
)

//XMQProducer Producer
type XMQProducer struct {
	conf        *Conf
	conn        net.Conn
	messages    chan *mq.ProcuderMessage
	backupMsg   chan *mq.ProcuderMessage
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

//NewXMQProducer 创建新的producer
func NewXMQProducer(address string, opts ...mq.Option) (producer *XMQProducer, err error) {
	producer = &XMQProducer{}
	producer.OptionConf = &mq.OptionConf{}
	producer.messages = make(chan *mq.ProcuderMessage, 10000)
	producer.backupMsg = make(chan *mq.ProcuderMessage, 100)
	producer.closeCh = make(chan struct{})
	for _, opt := range opts {
		opt(producer.OptionConf)
	}
	producer.conf, err = NewConf(producer.Raw)
	if err != nil {
		return nil, err
	}
	if producer.Logger == nil {
		producer.Logger = logger.GetSession("xmq.producer", logger.CreateSession())
	}
	producer.header = make([]string, 0, 4)
	return
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
	result, err := encoding.Decode(msg, "gbk")
	if err != nil {
		return err
	}
	_, err = producer.conn.Write(result)
	producer.lastWrite = time.Now()
	producer.writeLock.Unlock()
	return err
}

//sendLoop 循环发送消息
func (producer *XMQProducer) sendLoop() {
	if producer.done {
		producer.disconnect()
		return
	}
	if producer.Retry {
	Loop1:
		for {
			select {
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop1
				}
				msg.AddSendTimes()
				if msg.SendTimes > 3 {
					producer.Logger.Errorf("发送消息失败，丢弃消息(%s)(消息3次未发送成功)", msg.Queue)
				}
				//producer.Logger.Debug("发送消息...", msg.Data)
				err := producer.writeMessage(msg.Data)
				if err != nil {
					producer.Logger.Errorf("发送消息失败，稍后重新发送(%s)(err:%v)", msg.Queue, err)
					select {
					case producer.messages <- msg:
					default:
						producer.Logger.Errorf("发送失败，队列已满无法再次发送(%s):%s", msg.Queue, msg.Data)
					}
					break Loop1
				}
			}
		}
	} else {
	Loop2:
		for {
			select {
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop2
				}
				msg.AddSendTimes()
				//producer.Logger.Debug("发送消息...", msg.Data)
				err := producer.writeMessage(msg.Data)
				if err != nil {
					producer.Logger.Errorf("发送消息失败，放入备份队列(%s)(err:%v)", msg.Queue, err)
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("备份队列已满，无法放入队列(%s):%s", msg.Queue, msg.Data)
					}
					break Loop2
				}
			}
		}
	}
	if producer.done { //关闭连接
		producer.disconnect()
		return
	}
	producer.reconnect()
}
func (producer *XMQProducer) disconnect() {
	producer.isConnected = false
	if producer.conn == nil {
		return
	}
	producer.conn.Close()
	return
}

//GetBackupMessage 获取备份数据
func (producer *XMQProducer) GetBackupMessage() chan *mq.ProcuderMessage {
	return producer.backupMsg
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
	producer.Logger.Infof("连接到服务器:%s", producer.conf.Address)
	producer.conn, err = net.DialTimeout("tcp", producer.conf.Address, time.Second*2)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.isConnected = true
	producer.lastWrite = time.Now()
	go producer.sendLoop()
	return nil
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
	pm := &mq.ProcuderMessage{Queue: queue, Data: smessage, Timeout: timeout}
	select {
	case producer.messages <- pm:
		return nil
	default:
		return errors.New("producer无法连接到MQ服务器，消息队列已满无法发送")
	}
}

//Close 关闭当前连接
func (producer *XMQProducer) Close() {
	producer.done = true
	close(producer.closeCh)
	close(producer.messages)
	close(producer.backupMsg)
}

type xmqProducerResolver struct {
}

func (s *xmqProducerResolver) Resolve(address string, opts ...mq.Option) (mq.MQProducer, error) {
	return NewXMQProducer(address, opts...)
}
func init() {
	mq.RegisterProducer("xmq", &xmqProducerResolver{})
}
