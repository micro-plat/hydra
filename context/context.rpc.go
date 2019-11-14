package context

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/rpc"
)

type RPCInvoker interface {
	PreInit(services ...string) (err error)
	RequestFailRetry(service string, method string, header map[string]string, input map[string]interface{}, times int) (status int, result string, params map[string]string, err error)
	Request(service string, method string, header map[string]string, input map[string]interface{}, failFast bool) (status int, result string, param map[string]string, err error)
	AsyncRequest(service string, method string, header map[string]string, input map[string]interface{}, failFast bool) rpc.IRPCResponse
	WaitWithFailFast(callback func(string, int, string, error), timeout time.Duration, rs ...rpc.IRPCResponse) error
}

//IContextRPC rpc基础操作
type IContextRPC interface {
	PreInit(services ...string) error
	RequestFailRetry(service string, input map[string]interface{}, times int) (status int, r string, param map[string]string, err error)
	Request(service string, input map[string]interface{}, failFast bool) (status int, r string, param map[string]string, err error)
	RequestMap(service string, input map[string]interface{}, failFast bool) (status int, r map[string]interface{}, param map[string]string, err error)
}

//ContextRPC rpc操作实例
type ContextRPC struct {
	ctx *Context
	rpc RPCInvoker
}

func (cr *ContextRPC) clear() {
	cr.ctx = nil
	cr.rpc = nil
}

//Reset 重置context
func (cr *ContextRPC) reset(ctx *Context, rpc RPCInvoker) {
	cr.ctx = ctx
	cr.rpc = rpc
}

//PreInit 预加载服务
func (cr *ContextRPC) PreInit(services ...string) error {
	return cr.rpc.PreInit()
}

//AsyncRequest 异步请求
func (cr *ContextRPC) AsyncRequest(service string, header map[string]string, form map[string]interface{}, failFast bool) rpc.IRPCResponse {
	if header == nil {
		header = make(map[string]string)
	}
	if _, ok := header["X-Request-Id"]; !ok {
		header["X-Request-Id"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	method, ok := header["method"]
	if !ok {
		method = "get"
	}
	return cr.rpc.AsyncRequest(service, strings.ToUpper(method), header, form, failFast)
}

//RequestFailRetry RPC请求
func (cr *ContextRPC) RequestFailRetry(service string, header map[string]string, form map[string]interface{}, times int) (status int, r string, param map[string]string, err error) {
	if _, ok := header["X-Request-Id"]; !ok {
		header["X-Request-Id"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	method, ok := header["method"]
	if !ok {
		method = "get"
	}
	status, r, param, err = cr.rpc.RequestFailRetry(service, strings.ToUpper(method), header, form, times)
	if err != nil || status != 200 {
		return
	}
	return
}

//Request RPC请求
func (cr *ContextRPC) Request(service string, header map[string]string, form map[string]interface{}, failFast bool) (status int, r string, param map[string]string, err error) {
	if header == nil {
		header = map[string]string{}
	}
	if form == nil {
		form = map[string]interface{}{}
	}
	if _, ok := header["X-Request-Id"]; !ok {
		header["X-Request-Id"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	method, ok := header["method"]
	if !ok {
		method = "get"
	}
	status, r, param, err = cr.rpc.Request(service, strings.ToUpper(method), header, form, failFast)
	if err != nil || status != 200 {
		return
	}
	return
}

//RequestMap RPC请求返回结果转换为map
func (cr *ContextRPC) RequestMap(service string, header map[string]string, form map[string]interface{}, failFast bool) (status int, r map[string]interface{}, param map[string]string, err error) {
	if _, ok := header["X-Request-Id"]; !ok {
		header["X-Request-Id"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	status, result, param, err := cr.Request(service, header, form, failFast)
	if err != nil {
		return
	}

	r, err = jsons.Unmarshal([]byte(result))
	if err != nil {
		err = fmt.Errorf("rpc请求返结果不是有效的json串:%s,%v,%s,err:%v", service, form, result, err)
		return
	}
	return
}
