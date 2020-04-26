package queues

import "github.com/micro-plat/lib4go/queue"

//IQueue 消息队列
type IQueue = queue.IQueue

//IComponentQueue Component Queue
type IComponentQueue interface {
	GetRegularQueue(names ...string) (c IQueue)
	GetQueue(names ...string) (q IQueue, err error)
}
