package creator

import (
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

type mqcBuilder struct {
	customerBuilder
}

//newCronewMQCn 构建mqc生成器
func newMQC(addr string, opts ...mqc.Option) *mqcBuilder {
	b := &mqcBuilder{
		customerBuilder: make(map[string]interface{}),
	}
	b.customerBuilder["main"] = mqc.New(addr, opts...)
	return b
}

//Load 加载配置信息
func (b *mqcBuilder) Load() {
	queues := services.MQC.GetQueues()
	if q, ok := b.customerBuilder["queue"].(*queue.Queues); ok {
		q.Append(queues.Queues...)
		return
	}
	b.customerBuilder["queue"] = queues
	return
}

//Queue 添加队列配置
func (b *mqcBuilder) Queue(mq ...*queue.Queue) *mqcBuilder {
	f := func() {
		oqueue, ok := b.customerBuilder["queue"].(*queue.Queues)
		if !ok {
			oqueue = queue.NewQueues()
			b.customerBuilder["queue"] = oqueue
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
