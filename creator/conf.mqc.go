package creator

import (
	"github.com/micro-plat/hydra/conf/server/mqc"
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
	b.customerBuilder["queue"] = queues
	return
}
