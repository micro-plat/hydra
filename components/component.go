package components

import (
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/rpcs"
	"github.com/micro-plat/hydra/registry/conf/server"
)

//IRPC rpc请求对象
type IRPC interface {
	GetRegularRPC(names ...string) (c rpcs.IRequest)
	GetRPC(names ...string) (c rpcs.IRequest, err error)
}

//IComponent 组件
type IComponent interface {
	RPC() IRPC
}

//Component 组件
type Component struct {
	c   container.IContainer
	rpc IRPC
}

//NewComponent 创建组件
func NewComponent(conf server.IServerConf) *Component {
	c := &Component{
		c: container.NewContainer(conf.GetVarConf()),
	}
	c.rpc = rpcs.NewStandardRPC(c.c, conf.GetMainConf().GetPlatName(), conf.GetMainConf().GetSysName())
	return c
}

//RPC 获取rpc组件
func (c *Component) RPC() IRPC {
	return c.rpc
}
