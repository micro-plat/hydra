package creator

import (
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

type mqcBuilder struct {
	CustomerBuilder
}

//newCronewMQCn 构建mqc生成器
func newMQC(addr string, opts ...mqc.Option) *mqcBuilder {
	b := &mqcBuilder{
		CustomerBuilder: make(map[string]interface{}),
	}
	b.CustomerBuilder[ServerMainNodeName] = mqc.New(addr, opts...)
	return b
}

//Load 加载配置信息
func (b *mqcBuilder) Load() {
	queues := services.MQC.GetQueues()
	if q, ok := b.CustomerBuilder[queue.TypeNodeName].(*queue.Queues); ok {
		q.Append(queues.Queues...)
		return
	}
	b.CustomerBuilder[queue.TypeNodeName] = queues
	return
}

//Queue 添加队列配置
func (b *mqcBuilder) Queue(mq ...*queue.Queue) *mqcBuilder {
	f := func() {
		oqueue, ok := b.CustomerBuilder[queue.TypeNodeName].(*queue.Queues)
		if !ok {
			oqueue = queue.NewQueues()
			b.CustomerBuilder[queue.TypeNodeName] = oqueue
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
