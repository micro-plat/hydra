package rpcs

import (
	"context"
	"fmt"

	"github.com/micro-plat/hydra/components/pkgs"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	rpcconf "github.com/micro-plat/hydra/conf/vars/rpc"
	rc "github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	npkgs "github.com/micro-plat/hydra/pkgs"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

var requests = cmap.New(4)

//IRequest Component rpc
type IRequest interface {

	//Request request 请求
	Request(service string, input interface{}, opts ...rpc.RequestOption) (res *npkgs.Rspns, err error)

	//Swap 将当前请求参数作为RPC参数并发送RPC请求
	Swap(service string, ctx rc.IContext) (res *npkgs.Rspns, err error)

	//RequestByCtx RPC请求，可通过context撤销请求
	RequestByCtx(ctx context.Context, service string, input interface{}, opts ...rpc.RequestOption) (res *npkgs.Rspns, err error)
}

//Request RPC Request
type Request struct {
	conf    *rpcconf.RPCConf
	version int32
}

//NewRequest 构建请求
func NewRequest(version int32, conf *rpcconf.RPCConf) *Request {
	req := &Request{
		version: version,
		conf:    conf,
	}
	return req
}

//Request request 请求
func (r *Request) Request(service string, input interface{}, opts ...rpc.RequestOption) (res *npkgs.Rspns, err error) {

	//处理链路跟踪
	nopts := make([]rpc.RequestOption, 0, 2)
	nopts = append(nopts, opts...)
	if ctx, ok := rc.GetContext(); ok {
		nopts = append(opts, rpc.WithTraceID(ctx.User().GetTraceID()))
	}
	//发送请求
	return r.RequestByCtx(context.Background(), service, input, nopts...)
}

//Swap 将当前请求参数作为RPC参数并发送RPC请求
func (r *Request) Swap(service string, ctx rc.IContext) (res *npkgs.Rspns, err error) {

	//获取内容
	input := ctx.Request().GetMap()
	//处理链路跟踪
	opts := make([]rpc.RequestOption, 0, 2)
	opts = append(opts, rpc.WithTraceID(ctx.User().GetTraceID()))

	//复制请求头
	hd := make(map[string][]string)
	kv := ctx.Request().Headers()
	for k := range kv {
		hd[k] = []string{kv.GetString(k)}
	}
	opts = append(opts, rpc.WithHeaders(hd))

	// 发送请求
	return r.RequestByCtx(ctx.Context(), service, input, opts...)
}

//RequestByCtx RPC请求，可通过context撤销请求
func (r *Request) RequestByCtx(ctx context.Context, service string, input interface{}, opts ...rpc.RequestOption) (res *npkgs.Rspns, err error) {
	isip, rservice, platName, err := rpc.ResolvePath(service, global.Current().GetPlatName())
	if err != nil {
		return
	}
	//如果入参不是ip 通过注册中心去获取所请求平台的所有rpc服务子节点  再通过路由匹配获取真实的路由
	_, c, err := requests.SetIfAbsentCb(fmt.Sprintf("%s@%s.%d", rservice, platName, r.version), func(i ...interface{}) (interface{}, error) {

		if isip {
			return rpc.NewClientByConf(platName, "", rservice, r.conf)
		}
		return rpc.NewClientByConf(global.Def.RegistryAddr, platName, rservice, r.conf)
	})
	if err != nil {
		return nil, err
	}

	client := c.(*rpc.Client)
	nopts := make([]rpc.RequestOption, 0, len(opts)+1)
	nopts = append(nopts, opts...)
	if reqid := types.GetString(ctx.Value(rc.XRequestID)); reqid != "" {
		nopts = append(nopts, rpc.WithTraceID(reqid))
	} else {
		if ctx, ok := rc.GetContext(); ok {
			nopts = append(opts, rpc.WithTraceID(ctx.User().GetTraceID()))
		}
	}
	fm := pkgs.GetString(input)
	return client.RequestByString(ctx, rservice, fm, nopts...)
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
