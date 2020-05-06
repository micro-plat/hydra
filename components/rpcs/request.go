package rpcs

import (
	"fmt"

	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/lib4go/types"
)

//rpcTypeNode rpc在var配置中的类型名称
const rpcTypeNode = "rpc"

//rpcNameNode rpc名称在var配置中的末节点名称
const rpcNameNode = "rpc"

//IRequest Component rpc
type IRequest = rpc.IRequest

//IComponentRPC Component Cache
type IComponentRPC interface {
	GetRegularRPC(names ...string) (c IRequest)
	GetRPC(names ...string) (c IRequest, err error)
}

//Request RPC Request
type Request struct {
	plat      string
	server    string
	node      string
	container container.IContainer
}

//NewRequest 构建请求
func NewRequest(plat string, server string, nameNode string, container container.IContainer) *Request {
	return &Request{
		plat:      plat,
		server:    server,
		node:      types.GetString(nameNode, rpcNameNode),
		container: container,
	}
}

//Request RPC请求
func (r *Request) Request(service string, form map[string]interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	isip, rservice, domain, server, err := rpc.ResolvePath(service, r.plat, r.server)
	if err != nil {
		return
	}

	//获取配置版本号
	version, err := r.container.Conf().GetConfVersion(rpcTypeNode, r.node)
	if err != nil {
		return nil, err
	}

	c, err := r.container.GetOrCreate(fmt.Sprintf("__rpc_service_%d_%s@%s.%s", version, rservice, server, domain), func(i ...interface{}) (interface{}, error) {
		if isip {
			tls, err := r.container.Conf().GetConf(rpcTypeNode, r.node)
			if err != conf.ErrNoSetting {
				return nil, err
			}
			if len(tls.GetStrings("tls")) == 2 {
				return rpc.NewClient(service, rpc.WithTLS(tls.GetStrings("tls")))
			}
		}
		return rpc.NewClient(service)
	})
	if err != nil {
		return nil, err
	}
	client := c.(*rpc.Client)
	return client.Request(service, form, opts...)
}
