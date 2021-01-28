package queues

import (
	"github.com/micro-plat/hydra/components/pkgs"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
)

//IQueue 消息队列
type IQueue interface {
	Send(key string, value interface{}, requestID ...string) error
}

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

//Send 发送消息
func (q *queue) Send(key string, value interface{}, requestID ...string) error {
	hd := make([]string, 0, 2)
	if len(requestID) > 0 {
		hd = append(hd, context.XRequestID, requestID[0])
	} else {
		if ctx, ok := context.GetContext(); ok {
			hd = append(hd, context.XRequestID, ctx.User().GetTraceID())
		}
	}
	return q.q.Push(global.MQConf.GetQueueName(key), pkgs.GetStringByHeader(key, value, hd...))
}

func (q *queue) Close() error {
	return q.q.Close()
}
