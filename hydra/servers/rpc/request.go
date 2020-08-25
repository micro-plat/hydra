package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
)

//Request 处理任务请求
type Request struct {
	request *pb.RequestContext
	form    map[string]interface{}
	header  map[string]string
}

//NewRequest 构建任务请求
func NewRequest(request *pb.RequestContext) (r *Request, err error) {
	r = &Request{
		request: request,
		form:    make(map[string]interface{}),
		header:  make(map[string]string),
	}
	//处理请求头
	if err = json.Unmarshal([]byte(request.Header), &r.header); err != nil {
		return nil, fmt.Errorf("rpc请求头转换失败 %s %w", request.Header, err)
	}

	//处理正常请求参数
	if err = json.Unmarshal([]byte(request.Input), &r.form); err != nil {
		return nil, fmt.Errorf("rpc请求参数转换失败:%s %w", request.Input, err)
	}

	return r, nil
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.request.Service
}

//GetHost 远程Host
func (m *Request) GetHost() string {
	return m.header["Host"]
}

//GetService 服务名
func (m *Request) GetService() string {
	return m.request.Service
}

//GetMethod 方法名
func (m *Request) GetMethod() string {
	return m.request.Method
}

//GetForm 输入参数
func (m *Request) GetForm() map[string]interface{} {
	return m.form
}

//GetHeader 头信息
func (m *Request) GetHeader() map[string]string {
	return m.header
}

func (m *Request) getHeader(key string) string {
	return m.header[key]
}
