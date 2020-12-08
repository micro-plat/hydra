package ws

import (
	"encoding/json"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//Request 处理任务请求
type Request struct {
	method string
	form   types.XMap
	header map[string]string
}
type WSOption func(*Request)

func WithUUID(uuid string) WSOption {
	return func(req *Request) {
		req.header[context.XRequestID] = uuid
	}
}

func WithClientIP(clientIP string) WSOption {
	return func(req *Request) {
		req.header["Client-IP"] = clientIP
	}
}

//NewRequest 构建任务请求
func NewRequest(method string, content []byte, uuid string, clientip string, opts ...WSOption) (r *Request, err error) {
	r = &Request{
		method: method,
		form:   make(map[string]interface{}),
		header: map[string]string{
			context.XRequestID: uuid,
			"Client-IP":        clientip,
		},
	}
	for _, o := range opts {
		o(r)
	}

	json.Unmarshal(content, &r.form)

	//内容为json
	if _, ok := r.header["Content-Type"]; !ok {
		r.header["Content-Type"] = "application/json"
	}
	r.form["__body__"] = types.BytesToString(content)
	return r, nil
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.form.GetString("service", "/")
}

//GetHost 获取Client-IP
func (m *Request) GetHost() string {
	return m.header["Client-IP"]
}

//GetService 服务名
func (m *Request) GetService() string {
	return m.form.GetString("service", "/")
}

//GetMethod 方法名
func (m *Request) GetMethod() string {
	return m.method
}

//GetForm 输入参数
func (m *Request) GetForm() map[string]interface{} {
	return m.form
}

//GetHeader 头信息
func (m *Request) GetHeader() map[string]string {
	return m.header
}
