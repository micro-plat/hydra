package xmq

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/vars/queue/xmq"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/logger"
)

//Producer Producer
type Producer struct {
	conn           net.Conn
	connecting     bool
	isConnected    bool
	closeCh        chan struct{}
	done           bool
	lk             sync.Mutex
	writeLock      sync.Mutex
	lastWrite      time.Time
	firstConnected bool
	log            *logger.Logger
	confOpts       *xmq.XMQ
}

//NewProducer 创建新的producer
func NewProducer(confOpts *xmq.XMQ) (producer *Producer, err error) {
	producer = &Producer{
		log:      logger.GetSession("xmq", logger.CreateSession()),
		confOpts: confOpts,
	}
	producer.closeCh = make(chan struct{})

	return producer, producer.Connect()
}

//Connect  循环连接服务器
func (producer *Producer) Connect() error {
	err := producer.connectOnce()
	if err != nil {
		producer.log.Error(err)
	}
	go func() {
	START:
		for {
			select {
			case <-producer.closeCh:
				break START
			case <-time.After(time.Second * 3):
				if producer.done {
					return
				}
				if producer.isConnected {
					if time.Since(producer.lastWrite).Seconds() > 3 {
						message := newHeartBit()
						if producer.confOpts.SignKey != "" {
							message.signKey = producer.confOpts.SignKey
						}
						msgVal, err := message.Make()
						if err != nil {
							producer.log.Error(err)
							continue
						}
						err = producer.writeMessage(msgVal)
						if err == nil {
							continue
						}
						producer.log.Error(err)
						err = producer.connectOnce()
						if err != nil {
							producer.log.Error(err)
						}
					}
					continue
				}
				err = producer.connectOnce()
				if err != nil {
					producer.log.Error(err)
				}

			}
		}
	}()
	return nil
}

func (producer *Producer) writeMessage(msg string) error {
	if !producer.isConnected {
		return fmt.Errorf("未连接到服务器")
	}
	producer.writeLock.Lock()
	defer producer.writeLock.Unlock()
	producer.lastWrite = time.Now()
	result, err := encoding.Decode(msg, "gbk")
	if err != nil {
		return err
	}
	_, err = producer.conn.Write(result)
	producer.lastWrite = time.Now()
	if err != nil {
		producer.log.Warn("发送数据失败:", err)
		//producer.disconnect()
	}
	return err
}

func (producer *Producer) disconnect() {
	producer.isConnected = false
	if producer.conn == nil {
		return
	}
	producer.conn.Close()
	return
}

//reconnect 自动重连
func (producer *Producer) reconnect() {
	producer.conn.Close()
	producer.disconnect()
	err := producer.Connect()
	if err != nil {
		producer.log.Errorf("连接到MQ服务器失败:%v", err)
	}
}

//ConnectOnce 连接到服务器
func (producer *Producer) connectOnce() (err error) {
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
	producer.conn, err = net.DialTimeout("tcp", producer.confOpts.Address, time.Second*2)
	if err != nil {
		return fmt.Errorf("mq 无法连接到远程服务器:%v", err)
	}
	if !producer.firstConnected {
		producer.firstConnected = true
	} else {
		producer.log.Info("恢复连接:", producer.confOpts.Address)
	}
	producer.isConnected = true
	producer.lastWrite = time.Now()
	return nil
}

//Push 发送消息
func (producer *Producer) Push(queue string, msg string) error {
	return producer.Send(queue, msg, time.Hour*1000)
}

//Pop 拉取消息
func (producer *Producer) Pop(key string) (string, error) {
	return "", fmt.Errorf("not support")
}

//Count 获取消息条数
func (producer *Producer) Count(key string) (int64, error) {
	return 0, nil
}

//Send 发送消息
func (producer *Producer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("mq producer 已关闭")
	}
	if !producer.isConnected {
		return fmt.Errorf("producer无法连接到MQ服务器:%s", producer.confOpts.Address)
	}
	message := newMessage(queue, msg, int(timeout/time.Second))
	if producer.confOpts.SignKey != "" {
		message.signKey = producer.confOpts.SignKey
	}
	smessage, err := message.Make()
	if err != nil {
		return
	}
	return producer.writeMessage(smessage)
}

//Close 关闭当前连接
func (producer *Producer) Close() error {
	producer.done = true
	close(producer.closeCh)
	return nil
}

type presolver struct {
}

func (s *presolver) Resolve(confRaw string) (mq.IMQP, error) {
	return NewProducer(xmq.NewByRaw(confRaw))
}

func init() {
	mq.RegisterProducer("xmq", &presolver{})
}
