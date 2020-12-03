package internal

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
)

type Varcache struct {
	vars vars
}

func NewCache(internal map[string]map[string]interface{}) *Varcache {
	return &Varcache{
		vars: internal,
	}
}

func (c *Varcache) Redis(name string, q *cacheredis.Redis) *Varcache {
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

func (c *Varcache) GoCache(name string, q *gocache.GoCache) *Varcache {
	return c.Custom(name, q)
}

func (c *Varcache) Memcache(name string, q *memcached.Memcache) *Varcache {
	return c.Custom(name, q)
}

func (c *Varcache) Custom(name string, q interface{}) *Varcache {
	if _, ok := c.vars[cache.TypeNodeName]; !ok {
		c.vars[cache.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[cache.TypeNodeName][name] = q
	return c
}
