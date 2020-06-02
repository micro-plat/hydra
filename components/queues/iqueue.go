package queues

import "github.com/micro-plat/hydra/components/pkgs/mq"

//IQueue 消息队列
type IQueue = mq.IMQP

//IComponentQueue Component Queue
type IComponentQueue interface {
	GetRegularQueue(names ...string) (c IQueue)
	GetQueue(names ...string) (q IQueue, err error)
}
