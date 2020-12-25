package creator

import (
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

type mqcBuilder struct {
	BaseBuilder
}

//newCronewMQCn 构建mqc生成器
func newMQC(addr string, opts ...mqc.Option) *mqcBuilder {
	b := &mqcBuilder{
		BaseBuilder: make(map[string]interface{}),
	}
	b.BaseBuilder[ServerMainNodeName] = mqc.New(addr, opts...)
	return b
}

//Load 加载配置信息
func (b *mqcBuilder) Load() {
	queues := services.MQC.GetQueues()
	if q, ok := b.BaseBuilder[queue.TypeNodeName].(*queue.Queues); ok {
		q.Append(queues.Queues...)
		return
	}
	b.BaseBuilder[queue.TypeNodeName] = queues
	return
}

//Queue 添加队列配置
func (b *mqcBuilder) Queue(mq ...*queue.Queue) *mqcBuilder {
	f := func() {
		oqueue, ok := b.BaseBuilder[queue.TypeNodeName].(*queue.Queues)
		if !ok {
			oqueue = queue.NewQueues()
			b.BaseBuilder[queue.TypeNodeName] = oqueue
		}
		for _, m := range mq {
			m.Queue = global.MQConf.GetQueueName(m.Queue)
			oqueue.Append(m)
		}
	}
	//队列名称需在进行转换，进行延迟处理
	global.OnReady(f)
	return b
}
