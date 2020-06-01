package lmq

import (
	"errors"
	"time"

	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/queue/lmq"
)

//lmqProducer 基于本地channel的Producer
type lmqProducer struct {
	backupMsg chan *mq.ProcuderMessage
	client    *lmq.LMQClient
	closeCh   chan struct{}
	done      bool
}

//newlmqProducer 创建新的producer
func newlmqProducer(address string, opts ...mq.Option) (producer *lmqProducer, err error) {
	m, err := lmq.New([]string{address}, "")
	if err != nil {
		return nil, err
	}
	producer = &lmqProducer{client: m}
	producer.closeCh = make(chan struct{})
	return
}

//Connect  循环连接服务器
func (producer *lmqProducer) Connect() (err error) {
	return nil
}

//GetBackupMessage 获取备份数据
func (producer *lmqProducer) GetBackupMessage() chan *mq.ProcuderMessage {
	return producer.backupMsg
}

//Send 发送消息
func (producer *lmqProducer) Send(queue string, msg string, timeout time.Duration) (err error) {
	if producer.done {
		return errors.New("lmq producer 已关闭")
	}
	err = producer.client.Push(queue, msg)
	return
}

//Close 关闭当前连接
func (producer *lmqProducer) Close() {
	if !producer.done {
		producer.done = true
		close(producer.closeCh)
		close(producer.backupMsg)
		producer.client.Close()
	}
}

type lmqProducerResolver struct {
}

func (s *lmqProducerResolver) Resolve(address string, opts ...mq.Option) (mq.MQProducer, error) {
	return newlmqProducer(address, opts...)
}
func init() {
	mq.RegisterProducer("lmq", &lmqProducerResolver{})
}
