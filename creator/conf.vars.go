package creator

import (
	"github.com/micro-plat/hydra/conf/vars/cache"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
	cacheredis "github.com/micro-plat/hydra/conf/vars/cache/redis"

	"github.com/micro-plat/hydra/conf/vars/queue"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"

	"github.com/micro-plat/hydra/conf/vars/redis"

	"github.com/micro-plat/hydra/conf/vars/db"
	dbmysql "github.com/micro-plat/hydra/conf/vars/db/mysql"
	dboracle "github.com/micro-plat/hydra/conf/vars/db/oracle"

	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) Redis(name string, opts *redis.Redis) vars {
	if _, ok := v[redis.TypeNodeName]; !ok {
		v[redis.TypeNodeName] = make(map[string]interface{})
	}
	v[redis.TypeNodeName][name] = opts
	return v
}

//DB 添加db配置
func (v vars) DB() *vardb {
	return &vardb{vars: v}
}

func (v vars) Cache() *varcache {
	return &varcache{vars: v}
}

func (v vars) Queue() *varqueue {
	return &varqueue{vars: v}
}

type vardb struct {
	vars vars
}

func (c *vardb) Oracle(name string, q *dboracle.Oracle) *vardb {
	return c.Custom(name, q)
}

func (c *vardb) MySQL(name string, q *dbmysql.MySQL) *vardb {
	return c.Custom(name, q)
}

func (c *vardb) Custom(name string, q interface{}) *vardb {
	if _, ok := c.vars[db.TypeNodeName]; !ok {
		c.vars[db.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[db.TypeNodeName][name] = q
	return c
}

type varcache struct {
	vars vars
}

func (c *varcache) Redis(name string, q *cacheredis.Redis) *varcache {
	return c.Custom(name, q)
}

func (c *varcache) GoCache(name string, q *gocache.GoCache) *varcache {
	return c.Custom(name, q)
}

func (c *varcache) Memcache(name string, q *memcached.Memcache) *varcache {
	return c.Custom(name, q)
}

func (c *varcache) Custom(name string, q interface{}) *varcache {
	if _, ok := c.vars[cache.TypeNodeName]; !ok {
		c.vars[cache.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[cache.TypeNodeName][name] = q
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
	if _, ok := c.vars[queue.TypeNodeName]; !ok {
		c.vars[queue.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[queue.TypeNodeName][name] = q
	return c
}

func (v vars) RLog(service string, opts ...rlog.Option) vars {
	if _, ok := v[rlog.TypeNodeName]; !ok {
		v[rlog.TypeNodeName] = make(map[string]interface{})
	}
	v[rlog.TypeNodeName][rlog.LogName] = rlog.New(service, opts...)
	return v
}
