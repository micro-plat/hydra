package lmq

import (
	"fmt"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var queues cmap.ConcurrentMap

// Producer 消息生产者
type Producer struct {
}

// New 创建消息生产者
func New() (m *Producer, err error) {
	return &Producer{}, nil
}

//GetQueue 获取队列
func GetQueue(key string) (chan string, bool) {
	v, ok := queues.Get(key)
	if !ok {
		return nil, ok
	}
	return v.(chan string), true
}

//GetOrAddQueue 获取队列
func GetOrAddQueue(key string) chan string {
	_, i := queues.SetIfAbsent(key, make(chan string, 10000))
	return i.(chan string)
}

// Push 向存于 key 的列表的尾部插入所有指定的值
func (c *Producer) Push(key string, value string) error {
	ch := GetOrAddQueue(key)
	select {
	case ch <- value:
		return nil
	default:
		return fmt.Errorf("消息队列(%s)已满", key)
	}
}

// Pop 移除并且返回 key 对应的 list 的第一个元素。
func (c *Producer) Pop(key string) (string, error) {
	ch := GetOrAddQueue(key)
	v, ok := <-ch
	if !ok {
		return "", mq.Nil
	}
	return v, nil
}

// Count 队列中元素个数
func (c *Producer) Count(key string) (int64, error) {
	return int64(len(GetOrAddQueue(key))), nil
}

// Close 释放资源
func (c *Producer) Close() error {
	queues.RemoveIterCb(func(key string, v interface{}) bool {
		close(v.(chan string))
		return true
	})
	return nil
}

type lmqResolver struct {
}

func (s *lmqResolver) Resolve(confRaw string) (mq.IMQP, error) {
	return New()
}
func init() {
	queues = cmap.New(4)
	mq.RegisterProducer("lmq", &lmqResolver{})
}
