package components

import (
	"time"

	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/hydra/rpc"
)

//queueTypeNode queue在var配置中的类型名称
const rpcTypeNode = "queue"

//QueueNameInVar queue名称在var配置中的末节点名称
const rpcNameInVar = "queue"

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
	c       IComponents
	invoker IRPCInvoker
}

//NewStandardRPC 创建queue
func NewStandardRPC(c IComponents,platName string, systemName string, registryAddr string) *StandardQueue {
	return &StandardQueue{
		c:       c,
		invoker: rpc.NewInvoker(platName, systemName, registryAddr),
	}
}

//Request RPC请求
func (s *StandardRPC) Request(service string, input map[string]interface{}, failFast ...bool) (status int, r string, param map[string]string, err error) {
	header := map[string]string{}
	if input == nil {
		input = map[string]interface{}{}
	}
	if _, ok := header["X-Request-Id"]; !ok {
		header["X-Request-Id"] = s.c.GetRequestID()
	}
	status, r, param, err = s.invoker.Request(service, "GET", header, input, types.GetBoolByIndex(failFast, 0, true))
	if err != nil || status != 200 {
		return
	}
	return
}
