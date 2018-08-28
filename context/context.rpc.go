package context

import (
	"fmt"
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
func (cr *ContextRPC) AsyncRequest(service string, method string, header map[string]string, form map[string]interface{}, failFast bool) rpc.IRPCResponse {
	if header == nil {
		header = make(map[string]string)
	}
	if _, ok := header["__hydra_sid_"]; !ok {
		header["__hydra_sid_"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	return cr.rpc.AsyncRequest(service, method, header, form, failFast)
}

//RequestFailRetry RPC请求
func (cr *ContextRPC) RequestFailRetry(service string, method string, header map[string]string, form map[string]interface{}, times int) (status int, r string, param map[string]string, err error) {
	if _, ok := header["__hydra_sid_"]; !ok {
		header["__hydra_sid_"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	status, r, param, err = cr.rpc.RequestFailRetry(service, method, header, form, times)
	if err != nil || status != 200 {
		err = fmt.Errorf("rpc请求(%s)失败:%d,err:%v", service, status, err)
		return
	}
	return
}

//Request RPC请求
func (cr *ContextRPC) Request(service string, method string, header map[string]string, form map[string]interface{}, failFast bool) (status int, r string, param map[string]string, err error) {
	if _, ok := header["__hydra_sid_"]; !ok {
		header["__hydra_sid_"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	status, r, param, err = cr.rpc.Request(service, method, header, form, failFast)
	if err != nil || status != 200 {
		err = fmt.Errorf("rpc请求(%s)失败:%d,err:%v", service, status, err)
		return
	}
	return
}

//RequestMap RPC请求返回结果转换为map
func (cr *ContextRPC) RequestMap(service string, method string, header map[string]string, form map[string]interface{}, failFast bool) (status int, r map[string]interface{}, param map[string]string, err error) {
	if _, ok := header["__hydra_sid_"]; !ok {
		header["__hydra_sid_"] = cr.ctx.Request.GetUUID()
	}
	if _, ok := header["__body"]; !ok {
		header["__body"], _ = cr.ctx.Request.GetBody()
	}
	status, result, param, err := cr.Request(service, method, header, form, failFast)
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
