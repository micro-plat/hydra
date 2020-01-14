package rpcs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/rpc"
)

//IRPCInvoker RPC调用程序
type IRPCInvoker interface {
	PreInit(services ...string) (err error)
	RequestFailRetry(service string, method string, header map[string]string, input map[string]interface{}, times int) (status int, result string, params map[string]string, err error)
	Request(service string, method string, header map[string]string, input map[string]interface{}, failFast bool) (status int, result string, param map[string]string, err error)
	AsyncRequest(service string, method string, header map[string]string, input map[string]interface{}, failFast bool) rpc.IRPCResponse
	WaitWithFailFast(callback func(string, int, string, error), timeout time.Duration, rs ...rpc.IRPCResponse) error
}

//IComponentRPC Component rpc
type IComponentRPC interface {
	Request(service string, input map[string]interface{}, failFast ...bool) (status int, r string, param map[string]string, err error)
}

//StandardRPC rpc服务
type StandardRPC struct {
	c       components.IComponents
	invoker IRPCInvoker
}

//NewStandardRPC 创建RPC服务代理
func NewStandardRPC(c components.IComponents, platName string, systemName string, registryAddr string) *StandardQueue {
	return &StandardQueue{
		c:       c,
		invoker: rpc.NewInvoker(platName, systemName, registryAddr),
	}
}

//Request RPC请求
func (s *StandardRPC) Request(service string, input map[string]interface{}, opts ...Option) (res *RPCResponse, err error) {
	//处理可选参数
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	//发送远程请求
	status, r, param, err = s.invoker.Request(service, o.method, o.Params, input, o.failFast)
	if err != nil {
		return nil, err
	}
	return &RPCResponse{Status: status, Result: r, Params: param}, nil
}

//------------------------RPC响应---------------------------------------

//RPCResponse 请求结果
type RPCResponse struct {
	Status int
	Params map[string]interface{}
	Result string
}

//Success 请求是否成功
func (r *RPCResponse) Success() bool {
	return status == 200
}

//GetResult 获取请求结果
func (r *RPCResponse) GetResult() (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := json.Marshal([]byte(r.Result), &out)
	return out, err
}

//GetParam 根据KEY获取参数
func (r *RPCResponse) GetParam(key string) interface{} {
	return r.Params[key]
}

//-----------------RPC可选参数---------------------------------
type option struct {
	params   map[string]interface{}
	failFast bool
	method   string
}

func newOption() *option {
	return &Option{
		params:   map[string]string{"X-Request-Id": s.c.GetRequestID()},
		failFast: true,
		method:   "GET",
	}
}

//Option 配置选项
type Option func(*option)

//WithRPCParams RPC请求参数
func WithRPCParams(p map[string]interface{}) Option {
	return func(o *option) {
		o.params = p
	}
}

//WithRPCFailFast 快速失败
func WithRPCFailFast(b bool) Option {
	return func(o *option) {
		o.failFast = b
	}
}

//WithMethod 设置请求方法
func WithMethod(m string) Option {
	return func(o *option) {
		o.method = strings.ToUpper(m)
	}
}
