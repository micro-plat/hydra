package creator

import (
	"github.com/micro-plat/hydra/conf/vars/queue"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
)

//Varqueue 消息队列配置
type Varqueue struct {
	vars vars
}

//NewQueue 构建消息队列配置
func NewQueue(internal map[string]map[string]interface{}) *Varqueue {
	return &Varqueue{
		vars: internal,
	}
}

//Redis 添加Redis
func (c *Varqueue) Redis(nodeName string, address string, opts ...queueredis.Option) vars {
	return c.Custom(nodeName, queueredis.New(address, opts...))
}

//MQTT 添加MQTT
func (c *Varqueue) MQTT(nodeName string, address string, opts ...queuemqtt.Option) vars {
	return c.Custom(nodeName, queuemqtt.New(address, opts...))
}

//LMQ 添加本地内存作为消息队列
func (c *Varqueue) LMQ(nodeName string) vars {
	return c.Custom(nodeName, queuelmq.New())
}

//Custom 用户自定义消息队列
func (c *Varqueue) Custom(nodeName string, q interface{}) vars {
	if _, ok := c.vars[queue.TypeNodeName]; !ok {
		c.vars[queue.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[queue.TypeNodeName][nodeName] = q
	return c.vars
}
