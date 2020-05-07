package rpcs

import (
	"fmt"

	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//rpcTypeNode rpc在var配置中的类型名称
const rpcTypeNode = "rpc"

//rpcNameNode rpc名称在var配置中的末节点名称
const rpcNameNode = "rpc"

var requests = cmap.New(4)

//Request RPC Request
type Request struct {
	j *conf.JSONConf
}

//NewRequest 构建请求
func NewRequest(j *conf.JSONConf) *Request {
	return &Request{
		j: j,
	}
}

//Request RPC请求
func (r *Request) Request(service string, form map[string]interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	isip, rservice, domain, server, err := rpc.ResolvePath(service, application.Current().GetPlatName(), application.Current().GetSysName())
	if err != nil {
		return
	}
	_, c, err := requests.SetIfAbsentCb(fmt.Sprintf("%s@%s.%s_%d", rservice, server, domain, r.j.GetVersion()), func(i ...interface{}) (interface{}, error) {
		if isip {
			if len(r.j.GetStrings("tls")) == 2 {
				return rpc.NewClient(service, rpc.WithTLS(r.j.GetStrings("tls")))
			}
		}
		return rpc.NewClient(service)
	})
	if err != nil {
		return nil, err
	}
	client := c.(*rpc.Client)
	nopts := make([]rpc.RequestOption, 0, len(opts)+1)
	nopts = append(nopts, rpc.WithXRequestID(application.Current().CurrentContext().User().GetRequestID()))
	return client.Request(service, form, nopts...)
}

//Close 关闭RPC连接
func (r *Request) Close() error {
	requests.RemoveIterCb(func(key string, v interface{}) bool {
		client := v.(*rpc.Client)
		client.Close()
		return true
	})
	return nil
}
