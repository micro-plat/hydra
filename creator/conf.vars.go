package creator

import (
	"github.com/micro-plat/hydra/conf/vars/cache/gocache"
	"github.com/micro-plat/hydra/conf/vars/cache/memcached"
	cacheredis "github.com/micro-plat/hydra/conf/vars/cache/redis"
	"github.com/micro-plat/hydra/conf/vars/db"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) Redis(name string, redis *varredis.Redis) vars {
	if _, ok := v["redis"]; !ok {
		v["redis"] = make(map[string]interface{})
	}
	v["redis"][name] = redis
	return v
}

//DB 添加db配置
func (v vars) DB(name string, db *db.DB) vars {
	if _, ok := v["db"]; !ok {
		v["db"] = make(map[string]interface{})
	}
	v["db"][name] = db
	return v
}

func (v vars) Cache() *cache {
	return &cache{vars: v}
}

func (v vars) Queue() *varqueue {
	return &varqueue{vars: v}
}

func (v vars) RLog(service string, opts ...rlog.Option) vars {
	if _, ok := v[rlog.TypeNodeName]; !ok {
		v[rlog.TypeNodeName] = make(map[string]interface{})
	}
	v[rlog.TypeNodeName][rlog.LogName] = rlog.New(service, opts...)
	return v
}

type cache struct {
	vars vars
}

func (c *cache) Redis(name string, q *cacheredis.Redis) *cache {
	return c.Custom(name, q)
}

func (c *cache) GoCache(name string, q *gocache.GoCache) *cache {
	return c.Custom(name, q)
}

func (c *cache) Memcache(name string, q *memcached.Memcache) *cache {
	return c.Custom(name, q)
}

func (c *cache) Custom(name string, q interface{}) *cache {
	if _, ok := c.vars["cache"]; !ok {
		c.vars["cache"] = make(map[string]interface{})
	}
	c.vars["cache"][name] = q
	return c
}

type varqueue struct {
	vars vars
}

func (c *varqueue) Redis(name string, q *queueredis.Redis) *varqueue {
	return c.Custom(name, q)
}

func (c *varqueue) MQTT(name string, q *queuemqtt.MQTT) *varqueue {
	return c.Custom(name, q)
}

func (c *varqueue) LMQ(name string, q *queuelmq.LMQ) *varqueue {
	return c.Custom(name, q)
}

func (c *varqueue) Custom(name string, q interface{}) *varqueue {
	if _, ok := c.vars["queue"]; !ok {
		c.vars["queue"] = make(map[string]interface{})
	}
	c.vars["queue"][name] = q
	return c
}
