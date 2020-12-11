package creator

import (
	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
)

//Varcache 缓存配置对象
type Varcache struct {
	vars vars
}

//NewCache 构建缓存管理对象
func NewCache(internal map[string]map[string]interface{}) *Varcache {
	return &Varcache{
		vars: internal,
	}
}

//Redis 添加redis缓存
func (c *Varcache) Redis(nodeName string, address string, opts ...cacheredis.Option) vars {
	return c.Custom(nodeName, cacheredis.New(address, opts...))
}

//GoCache 添加本地内存作为缓存
func (c *Varcache) GoCache(nodeName string, opts ...gocache.Option) vars {
	return c.Custom(nodeName, gocache.New(opts...))
}

//Memcache 添加memcached作为缓存
func (c *Varcache) Memcache(nodeName string, addr string, opts ...memcached.Option) vars {
	return c.Custom(nodeName, memcached.New(addr, opts...))
}

//Custom 自定义缓存配置
func (c *Varcache) Custom(nodeName string, q interface{}) vars {
	if _, ok := c.vars[cache.TypeNodeName]; !ok {
		c.vars[cache.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[cache.TypeNodeName][nodeName] = q
	return c.vars
}
