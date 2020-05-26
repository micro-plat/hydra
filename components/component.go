package components

import (
	"github.com/micro-plat/hydra/components/caches"
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/dlock"
	"github.com/micro-plat/hydra/components/queues"
	"github.com/micro-plat/hydra/components/rpcs"
	"github.com/micro-plat/hydra/components/uuid"
	"github.com/micro-plat/hydra/global"
)

//IComponent 组件
type IComponent interface {
	RPC() rpcs.IComponentRPC
	Queue() queues.IComponentQueue
	Cache() caches.IComponentCache
}

//Def 默认组件
var Def IComponent = NewComponent()

//Component 组件
type Component struct {
	c     container.IContainer
	rpc   rpcs.IComponentRPC
	queue queues.IComponentQueue
	cache caches.IComponentCache
}

//NewComponent 创建组件
func NewComponent() *Component {
	c := &Component{
		c: container.NewContainer(),
	}
	c.rpc = rpcs.NewStandardRPC(c.c)
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

//DLock 获取分布式鍞
func (c *Component) DLock(name string) (dlock.ILock, error) {
	return dlock.NewLock(name, global.Def.RegistryAddr, global.Def.Context().Log())
}

//UUID 获取全局用户编号
func (c *Component) UUID() uuid.IUUID {
	return uuid.Get()
}
