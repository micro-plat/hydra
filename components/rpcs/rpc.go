package rpcs

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

//rpcTypeNode rpc在var配置中的类型名称
const rpcTypeNode = "rpc"

//rpcNameNode rpc名称在var配置中的末节点名称
const rpcNameNode = "rpc"

//StandardRPC rpc服务
type StandardRPC struct {
	c      components.IComponents
	client *rpc.Client
}

//NewStandardRPC 创建RPC服务代理
func NewStandardRPC(c components.IComponents, platName string, systemName string, registryAddr string) *StandardRPC {
	return &StandardRPC{
		c: c,
	}
}

//GetRegularRPC 获取正式的没有异常缓存实例
func (s *StandardRPC) GetRegularRPC(names ...string) (c IRequest) {
	c, err := s.GetRPC(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetRPC 获取缓存操作对象
func (s *StandardRPC) GetRPC(names ...string) (c IRequest, err error) {
	name := types.GetStringByIndex(names, 0, rpcNameNode)
	obj, err := s.c.GetOrCreateByConf(rpcTypeNode, name, func(c conf.IConf) (interface{}, error) {
		return &Request{}, nil
	})
	if err != nil {
		return nil, err
	}
	return obj.(IRequest), nil
}
