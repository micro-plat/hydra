package rpc

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/rpc"
)

//Response 异步请求的响应内容
type Response struct {
	Service string
	Result  chan rpc.IRPCResult
}

//Result 请求执行结果
type Result struct {
	Service string
	Status  int
	Result  string
	Params  map[string]string
	Err     error
}

//GetService 获取服务名称
func (r *Result) GetService() string {
	return r.Service
}

//GetStatus 获取状态码
func (r *Result) GetStatus() int {
	return r.Status
}

//GetResult 获取执行结果
func (r *Result) GetResult() string {
	return r.Result
}

//GetParams 获取执行结果
func (r *Result) GetParams() map[string]string {
	return r.Params
}

//GetErr 获取执行错误信息
func (r *Result) GetErr() error {
	return r.Err
}

//NewResponse 构建异步请求响应
func NewResponse(service string) *Response {
	return &Response{Service: service, Result: make(chan rpc.IRPCResult, 1)}
}

//Wait 等待请求返回
func (r *Response) Wait(timeout time.Duration) (int, string, map[string]string, error) {
	select {
	case <-time.After(timeout):
		return 504, "", nil, fmt.Errorf("%s请求超时(%v)", r.Service, timeout)
	case value := <-r.Result:
		return value.GetStatus(), value.GetResult(), value.GetParams(), value.GetErr()
	}
}

//GetResult 获取响应的近观回结果
func (r *Response) GetResult() chan rpc.IRPCResult {
	return r.Result
}

//GetService 获取服务
func (r *Response) GetService() string {
	return r.Service
}
