package queues

import (
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/global"
)

//IQueue 消息队列
type IQueue = mq.IMQP

//IComponentQueue Component Queue
type IComponentQueue interface {
	GetRegularQueue(names ...string) (c IQueue)
	GetQueue(names ...string) (q IQueue, err error)
}

//queue 对输入KEY进行封装处理
type queue struct {
	q mq.IMQP
}

func newQueue(proto string, confRaw string) (q *queue, err error) {
	q = &queue{}
	q.q, err = mq.NewMQP(proto, confRaw)
	return q, err
}
func (q *queue) Push(key string, value string) error {
	return q.q.Push(global.MQConf.GetQueueName(key), value)
}
func (q *queue) Pop(key string) (string, error) {
	return q.q.Pop((global.MQConf.GetQueueName(key)))
}
func (q *queue) Count(key string) (int64, error) {
	return q.q.Count(global.MQConf.GetQueueName(key))
}
func (q *queue) Close() error {
	return q.q.Close()
}
