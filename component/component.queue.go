package component

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/queue"
)

//QueueTypeNameInVar queue在var配置中的类型名称
const QueueTypeNameInVar = "queue"

//QueueNameInVar queue名称在var配置中的末节点名称
const QueueNameInVar = "queue"

//IComponentQueue Component Queue
type IComponentQueue interface {
	GetRegularQueue(names ...string) (c queue.IQueue)
	GetQueue(names ...string) (q queue.IQueue, err error)
	GetQueueBy(tpName string, name string) (c queue.IQueue, err error)
	SaveQueueObject(tpName string, name string, f func(c conf.IConf) (queue.IQueue, error)) (bool, queue.IQueue, error)
	Close() error
}

//StandardQueue queue
type StandardQueue struct {
	IContainer
	name       string
	queueCache cmap.ConcurrentMap
}

//NewStandardQueue 创建queue
func NewStandardQueue(c IContainer, name ...string) *StandardQueue {
	if len(name) > 0 {
		return &StandardQueue{IContainer: c, name: name[0], queueCache: cmap.New(2)}
	}
	return &StandardQueue{IContainer: c, name: QueueNameInVar, queueCache: cmap.New(2)}
}

//GetRegularQueue 获取正式的没有异常Queue实例
func (s *StandardQueue) GetRegularQueue(names ...string) (c queue.IQueue) {
	c, err := s.GetQueue(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetQueue GetQueue
func (s *StandardQueue) GetQueue(names ...string) (q queue.IQueue, err error) {
	name := s.name
	if len(names) > 0 {
		name = names[0]
	}
	return s.GetQueueBy(QueueTypeNameInVar, name)
}

//GetQueueBy 根据类型获取缓存数据
func (s *StandardQueue) GetQueueBy(tpName string, name string) (c queue.IQueue, err error) {
	_, c, err = s.SaveQueueObject(tpName, name, func(jConf conf.IConf) (queue.IQueue, error) {
		var qConf conf.QueueConf
		if err = jConf.Unmarshal(&qConf); err != nil {
			return nil, err
		}
		if b, err := govalidator.ValidateStruct(&qConf); !b {
			return nil, err
		}
		return queue.NewQueue(qConf.Proto, string(jConf.GetRaw()))
	})
	return c, err
}

//SaveQueueObject 缓存对象
func (s *StandardQueue) SaveQueueObject(tpName string, name string, f func(c conf.IConf) (queue.IQueue, error)) (bool, queue.IQueue, error) {
	cacheConf, err := s.IContainer.GetVarConf(tpName, name)
	if err != nil {
		return false, nil, fmt.Errorf("%s %v", registry.Join("/", s.GetPlatName(), "var", tpName, name), err)
	}
	key := fmt.Sprintf("%s/%s:%d", tpName, name, cacheConf.GetVersion())
	ok, ch, err := s.queueCache.SetIfAbsentCb(key, func(input ...interface{}) (c interface{}, err error) {
		return f(cacheConf)
	})
	if err != nil {
		err = fmt.Errorf("创建queue失败:%s,err:%v", string(cacheConf.GetRaw()), err)
		return ok, nil, err
	}
	return ok, ch.(queue.IQueue), err
}

//Close 释放所有缓存配置
func (s *StandardQueue) Close() error {
	s.queueCache.RemoveIterCb(func(k string, v interface{}) bool {
		v.(queue.IQueue).Close()
		return true
	})
	return nil
}
