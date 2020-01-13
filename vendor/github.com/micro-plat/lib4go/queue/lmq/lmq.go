package lmq

import (
	"fmt"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/queue"
)

var queues cmap.ConcurrentMap

// LMQClient memcache配置文件
type LMQClient struct {
}

// New 根据配置文件创建一个redis连接
func New(addrs []string, conf string) (m *LMQClient, err error) {
	return &LMQClient{}, nil
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
func (c *LMQClient) Push(key string, value string) error {
	ch := GetOrAddQueue(key)
	select {
	case ch <- value:
		return nil
	default:
		return fmt.Errorf("消息队列(%s)已满", key)
	}
}

// Pop 移除并且返回 key 对应的 list 的第一个元素。
func (c *LMQClient) Pop(key string) (string, error) {
	ch := GetOrAddQueue(key)
	v, ok := <-ch
	if !ok {
		return "", queue.Nil
	}
	return v, nil
}

// Count 队列中元素个数
func (c *LMQClient) Count(key string) (int64, error) {
	return int64(len(GetOrAddQueue(key))), nil
}

// Close 释放资源
func (c *LMQClient) Close() error {
	queues.RemoveIterCb(func(key string, v interface{}) bool {
		close(v.(chan string))
		return true
	})
	return nil
}

type lmqResolver struct {
}

func (s *lmqResolver) Resolve(address []string, conf string) (queue.IQueue, error) {
	return New(address, conf)
}
func init() {
	queues = cmap.New(4)
	queue.Register("lmq", &lmqResolver{})
}
