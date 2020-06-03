package rpcs

import (
	"context"
	"fmt"

	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf"
	r "github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//rpcTypeNode rpc在var配置中的类型名称
const rpcTypeNode = "rpc"

//rpcNameNode rpc名称在var配置中的末节点名称
const rpcNameNode = "rpc"

var requests = cmap.New(4)

//IRequest Component rpc
type IRequest interface {
	rpc.IRequest
	RequestByCtx(service string, ctx r.IContext) (res *rpc.Response, err error)
}

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

//RequestByCtx 将当前请求转化为RPC调用
func (r *Request) RequestByCtx(service string, ctx r.IContext) (res *rpc.Response, err error) {
	input, err := ctx.Request().GetData()
	if err != nil {
		return nil, err
	}
	headers := ctx.Request().Path().GetHeaders()
	return r.Request(ctx.Context(), service, input,
		rpc.WithHeaders(headers), rpc.WithXRequestID(ctx.User().GetRequestID()))
}

//Request RPC请求
func (r *Request) Request(ctx context.Context, service string, form map[string]interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	isip, rservice, domain, server, err := rpc.ResolvePath(service, global.Current().GetPlatName(), global.Current().GetSysName())
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
	nopts = append(nopts, rpc.WithXRequestID(ctx.Value("X-Request-Id").(string)))
	return client.Request(ctx, service, form, nopts...)
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
