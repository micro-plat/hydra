package rpcs

import (
	"fmt"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf"
)

//IRequest Component rpc
type IRequest = rpc.IRequest

//IComponentRPC Component Cache
type IComponentRPC interface {
	GetRegularRPC(names ...string) (c IRequest)
	GetRPC(names ...string) (c IRequest, err error)
}

//Request RPC Request
type Request struct {
	plat     string
	server   string
	registry string
	conf     conf.IConf
	c        components.IComponents
}

//Request RPC请求
func (r *Request) Request(service string, form map[string]interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	isip, rservice, domain, server, err := rpc.ResolvePath(service, r.plat, r.server)
	if err != nil {
		return
	}
	c, err := r.c.GetOrCreate(fmt.Sprintf("%s@%s.%s", rservice, server, domain), func(i ...interface{}) (interface{}, error) {
		if isip {
			tls := r.conf.GetStrings("tls")
			if len(tls) == 2 {
				return rpc.NewClient(service, rpc.WithTLS(tls))
			}
			return rpc.NewClient(service)
		}
	})
	if err != nil {
		return nil, err
	}
	client := c.(*rpc.Client)
	return client.Request(service, form, opts...)
}
