package rpcs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	rpcconf "github.com/micro-plat/hydra/conf/vars/rpc"
	r "github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var requests = cmap.New(4)

//IRequest Component rpc
type IRequest interface {
	Request(ctx context.Context, service string, input interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error)
	RequestByCtx(service string, ctx r.IContext) (res *rpc.Response, err error)
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

//RequestByCtx 将当前请求转化为RPC调用
func (r *Request) RequestByCtx(service string, ctx r.IContext) (res *rpc.Response, err error) {
	input, err := ctx.Request().GetMap()
	if err != nil {
		return nil, err
	}
	headers := ctx.Request().Path().GetHeaders()
	return r.Request(ctx.Context(), service, input,
		rpc.WithHeaders(headers), rpc.WithXRequestID(ctx.User().GetRequestID()))
}

//Request RPC请求
func (r *Request) Request(ctx context.Context, service string, input interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	isip, rservice, platName, err := rpc.ResolvePath(service, global.Current().GetPlatName())
	if err != nil {
		return
	}

	//获取rpc所有的注册服务路由
	// matchPath := conf.PathMatch
	rpcConf, _ := app.Cache.GetAPPConf("rpc")
	routerObj, _ := rpcConf.GetRouterConf()
	paths := []string{}
	for _, r := range routerObj.GetRouters() {
		paths = append(paths, r.Path)
	}

	matchPath := conf.NewPathMatch(paths...)
	matchPath.Match(rservice)

	_, c, err := requests.SetIfAbsentCb(fmt.Sprintf("%s@%s.%d", rservice, platName, r.version), func(i ...interface{}) (interface{}, error) {
		if isip {
			return rpc.NewClientByConf(platName, "", "", r.conf)
		}
		//return rpc.NewClient(global.Def.RegistryAddr, rpc.WithLocalFirstBalancer(platName, rservice, pkgs.LocalIP()))

		return rpc.NewClientByConf(global.Def.RegistryAddr, platName, rservice, r.conf)
	})
	if err != nil {
		return nil, err
	}

	client := c.(*rpc.Client)
	nopts := make([]rpc.RequestOption, 0, len(opts)+1)
	nopts = append(nopts, opts...)
	nopts = append(nopts, rpc.WithXRequestID(fmt.Sprint(ctx.Value("X-Request-Id"))))
	fm := getRequestForm(input)
	return client.Request(ctx, rservice, fm, nopts...)
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
func getRequestForm(content interface{}) map[string]interface{} {
	mp := make(map[string]interface{})
	if content == nil {
		return mp
	}
	switch v := content.(type) {
	case []byte:
		return map[string]interface{}{
			"__body_": v,
		}
	case string:
		return map[string]interface{}{
			"__body_": v,
		}
	}

	//反射检查
	tp := reflect.TypeOf(content).Kind()
	value := reflect.ValueOf(content)
	if tp == reflect.Ptr {
		value = value.Elem()
	}
	switch tp {
	case reflect.String:
		return map[string]interface{}{
			"__body_": content,
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return map[string]interface{}{
			"__body_": fmt.Sprint(content),
		}
	default:
		buff, err := json.Marshal(content)
		if err != nil {
			panic(fmt.Errorf("将请求转换为json串时错误%w", err))
		}
		return map[string]interface{}{
			"__body_": string(buff),
		}
	}
}
