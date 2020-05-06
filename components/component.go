package components

import (
	"github.com/micro-plat/hydra/components/caches"
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/queues"
	"github.com/micro-plat/hydra/components/rpcs"
	"github.com/micro-plat/hydra/registry/conf/server"
)

//IComponent 组件
type IComponent interface {
	RPC() rpcs.IComponentRPC
	Queue() queues.IComponentQueue
	Cache() caches.IComponentCache
}

//Component 组件
type Component struct {
	c     container.IContainer
	rpc   rpcs.IComponentRPC
	queue queues.IComponentQueue
	cache caches.IComponentCache
}

//NewComponent 创建组件
func NewComponent(conf server.IServerConf) *Component {
	c := &Component{
		c: container.NewContainer(conf.GetVarConf()),
	}
	c.rpc = rpcs.NewStandardRPC(c.c, conf.GetMainConf().GetPlatName(), conf.GetMainConf().GetSysName())
	c.queue = queues.NewStandardQueue(c.c)
	c.cache = caches.NewStandardCache(c.c)
	return c
}

//RPC 获取rpc组件
func (c *Component) RPC() rpcs.IComponentRPC {
	return c.rpc
}

//Queue 获取Queue组件
func (c *Component) Queue() queues.IComponentQueue {
	return c.queue
}

//Cache 获取Queue组件
func (c *Component) Cache() caches.IComponentCache {
	return c.cache
}
