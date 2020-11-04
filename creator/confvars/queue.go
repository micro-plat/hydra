package confvars

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/vars/queue"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
)

type Varqueue struct {
	vars vars
}

func NewQueue(confVars map[string]map[string]interface{}) *Varqueue {
	return &Varqueue{
		vars: confVars,
	}
}

func (c *Varqueue) Redis(name string, q *queueredis.Redis) *Varqueue {
	if q.ConfigName != "" {
		redisCfg, ok := c.vars["redis"]
		if !ok {
			panic(fmt.Errorf("请确认已配置/var/redis"))
		}
		redisObj, ok := redisCfg[q.ConfigName]
		if !ok {
			panic(fmt.Errorf("请确认已配置/var/redis/%s", q.ConfigName))
		}
		q.Redis = redisObj.(*varredis.Redis)
	}
	return c.Custom(name, q)
}

func (c *Varqueue) MQTT(name string, q *queuemqtt.MQTT) *Varqueue {
	return c.Custom(name, q)
}

func (c *Varqueue) LMQ(name string, q *queuelmq.LMQ) *Varqueue {
	return c.Custom(name, q)
}

func (c *Varqueue) Custom(name string, q interface{}) *Varqueue {
	if _, ok := c.vars[queue.TypeNodeName]; !ok {
		c.vars[queue.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[queue.TypeNodeName][name] = q
	return c
}
