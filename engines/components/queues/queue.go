package queues

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/queue"
	"github.com/micro-plat/lib4go/types"
)

//queueTypeNode queue在var配置中的类型名称
const queueTypeNode = "queue"

//queueNameNode queue名称在var配置中的末节点名称
const queueNameNode = "queue"

//StandardQueue queue
type StandardQueue struct {
	c components.IComponentContainer
}

//NewStandardQueue 创建queue
func NewStandardQueue(c components.IComponentContainer) *StandardQueue {
	return &StandardQueue{c: c}
}

//GetRegularQueue 获取正式的没有异常Queue实例
func (s *StandardQueue) GetRegularQueue(names ...string) (c IQueue) {
	c, err := s.GetQueue(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetQueue GetQueue
func (s *StandardQueue) GetQueue(names ...string) (q IQueue, err error) {
	name := types.GetStringByIndex(names, 0, queueNameNode)
	obj, err := s.c.GetOrCreateByConf(queueTypeNode, name, func(c conf.IConf) (interface{}, error) {
		return queue.NewQueue(c.GetString("proto"), string(c.GetRaw()))
	})
	if err != nil {
		return nil, err
	}
	return obj.(IQueue), nil
}
