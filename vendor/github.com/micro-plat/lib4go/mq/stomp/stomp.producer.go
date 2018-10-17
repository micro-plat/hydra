package stomp

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/mq"
	"github.com/zkfy/stompngo"
)

//StompProducer Producer
type StompProducer struct {
	address     string
	conn        *stompngo.Connection
	messages    chan *mq.ProcuderMessage
	backupMsg   chan *mq.ProcuderMessage
	queues      cmap.ConcurrentMap
	connecting  bool
	isConnected bool
	closeCh     chan struct{}
	done        bool
	lk          sync.Mutex
	header      []string
	*mq.OptionConf
}

//NewStompProducer 创建新的producer
func NewStompProducer(address string, opts ...mq.Option) (producer *StompProducer, err error) {
	producer = &StompProducer{address: address}
	producer.OptionConf = &mq.OptionConf{}
	producer.messages = make(chan *mq.ProcuderMessage, 10000)
	producer.backupMsg = make(chan *mq.ProcuderMessage, 100)
	producer.closeCh = make(chan struct{})
	for _, opt := range opts {
		opt(producer.OptionConf)
	}
	if producer.Logger == nil {
		producer.Logger = logger.GetSession("mq.producer", logger.CreateSession())
	}
	if strings.EqualFold(producer.OptionConf.Version, "") {
		producer.OptionConf.Version = "1.1"
	}
	if strings.EqualFold(producer.OptionConf.Persistent, "") {
		producer.OptionConf.Persistent = "true"
	}
	if strings.EqualFold(producer.OptionConf.Ack, "") {
		producer.OptionConf.Ack = "client-individual"
	}
	producer.header = stompngo.Headers{"accept-version", producer.OptionConf.Version}
	return
}

//Connect  循环连接服务器
func (producer *StompProducer) Connect() error {
	err := producer.connectOnce()
	if err == nil {
		return nil
	}
	producer.Logger.Error(err)
	go func() {
	START:
		for {
			select {
			case <-producer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				err = producer.connectOnce()
				if err == nil {
					break START
				}
				producer.Logger.Error(err)
			}
		}
	}()
	return nil
}

//sendLoop 循环发送消息
func (producer *StompProducer) sendLoop() {
	if producer.done {
		producer.disconnect()
		return
	}
	if producer.Retry {
	Loop1:
		for {
			select {
			case msg, ok := <-producer.backupMsg:
				if !ok {
					break Loop1
				}
				err := producer.conn.Send(msg.Headers, msg.Data)
				if err != nil {
					producer.Logger.Errorf("发送消息失败，放入备份队列(%s)(err:%v)", msg.Queue, err)
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("重试发送失败，备份队列已满无法放入队列(%s):%s", msg.Queue, msg.Data)
					}
					break Loop1
				}
			case msg, ok := <-producer.messages:
				if !ok {
					break Loop1
				}
				err := producer.conn.Send(msg.Headers, msg.Data)
				if err != nil {
					producer.Logger.Errorf("发送消息失败，放入备份队列(%s)(err:%v)", msg.Queue, err)
					select {
					case producer.backupMsg <- msg:
					default:
						producer.Logger.Errorf("发送失败，备份队列已满无法放入队列(%s):%s", msg.Queue, msg.Data)
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
				err := producer.conn.Send(msg.Headers, msg.Data)
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
func (producer *StompProducer) disconnect() {
	if producer.conn == nil || !producer.conn.Connected() {
		return
	}
	producer.conn.Disconnect(stompngo.Headers{})
	return
}

//GetBackupMessage 获取备份数据
func (producer *StompProducer) GetBackupMessage() chan *mq.ProcuderMessage {
	return producer.backupMsg
}

//reconnect 自动重连
func (producer *StompProducer) reconnect() {
	producer.conn.Disconnect(stompngo.Headers{})
	err := producer.Connect()
	if err != nil {
		producer.Logger.Errorf("连接到MQ服务器失败:%v", err)
	}
}

//ConnectOnce 连接到服务器
func (producer *StompProducer) connectOnce() (err error) {
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
	producer.Logger.Infof("重新连接到服务器:%s", producer.address)
	con, err := net.Dial("tcp", producer.address)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	producer.conn, err = stompngo.Connect(con, producer.header)
	if err != nil {
		return fmt.Errorf("mq 无法连接到MQ:%v", err)
	}

	go producer.sendLoop()
	return nil
}

//Send 发送消息
func (producer *StompProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	if !producer.connecting && producer.Retry {
		return fmt.Errorf("producer无法连接到MQ服务器:%s", producer.address)
	}

	pm := &mq.ProcuderMessage{Queue: queue, Data: msg, Timeout: timeout}
	pm.Headers = make([]string, 0, len(producer.header)+4)
	copy(pm.Headers, producer.header)

	pm.Headers = append(pm.Headers, "destination", "/queue/"+queue)
	if timeout > 0 && timeout < time.Second*10 {
		return fmt.Errorf("超时时长不能小于10秒:%s,%v", queue, timeout)
	}
	if timeout > 0 {
		pm.Headers = append(pm.Headers, "expires",
			fmt.Sprintf("%d000", time.Now().Add(timeout).Unix()))
	}
	select {
	case producer.messages <- pm:
		return nil
	default:
		return errors.New("producer无法连接到MQ服务器，消息队列已满无法发送")
	}
}

//Close 关闭当前连接
func (producer *StompProducer) Close() {
	producer.done = true
	close(producer.closeCh)
	close(producer.messages)
	close(producer.backupMsg)
}

type stompProducerResolver struct {
}

func (s *stompProducerResolver) Resolve(address string, opts ...mq.Option) (mq.MQProducer, error) {
	return NewStompProducer(address, opts...)
}
func init() {
	mq.RegisterProducer("stomp", &stompProducerResolver{})
}
