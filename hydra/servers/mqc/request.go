package mqc

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/lib4go/encoding/base64"

	"github.com/micro-plat/hydra/components/pkgs"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/lib4go/types"
)

//DefMethod DefMethod
var DefMethod = "GET"

//Request 处理任务请求
type Request struct {
	queue *queue.Queue
	mq.IMQCMessage
	method string
	form   map[string]interface{}
	header map[string]string
}

//NewRequest 构建任务请求
func NewRequest(queue *queue.Queue, m mq.IMQCMessage) (r *Request, err error) {
	if pkgs.IsOriginalQueue(queue.Queue) {
		return newOldRequest(queue, m)
	}
	return newRequest(queue, m)
}

func newOldRequest(queue *queue.Queue, m mq.IMQCMessage) (r *Request, err error) {
	r = &Request{
		IMQCMessage: m,
		queue:       queue,
		method:      DefMethod,
		form:        make(map[string]interface{}),
		header:      make(map[string]string),
	}
	//将消息原串转换为map
	input := make(map[string]interface{})
	message := m.GetMessage()
	json.Unmarshal(types.StringToBytes(message), &input)

	r.form = input
	//检查是否包含头信息
	r.form["__body__"] = message
	r.header["Content-Type"] = "application/json"

	//处理头信息
	r.header["__all__"] = message
	return
}

func newRequest(queue *queue.Queue, m mq.IMQCMessage) (r *Request, err error) {
	r = &Request{
		IMQCMessage: m,
		queue:       queue,
		method:      DefMethod,
		form:        make(map[string]interface{}),
		header:      make(map[string]string),
	}

	//将消息原串转换为map
	input := make(map[string]interface{})
	message := m.GetMessage()
	json.Unmarshal(types.StringToBytes(message), &input)

	//检查是否包含头信息
	r.form["__body__"] = message
	if v, ok := input["__header__"].(map[string]interface{}); ok {
		for n, m := range v {
			r.header[n] = fmt.Sprint(m)
		}
	}

	//处理头信息
	r.header["__all__"] = message
	if _, ok := r.header["Content-Type"]; !ok {
		r.header["Content-Type"] = "application/json"
	}

	if v, ok := input["__data__"].(string); ok {
		buff, _ := base64.DecodeBytes(v)
		r.form["__body__"] = string(buff)
	}
	return r, nil
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.queue.Queue
}

//GetService 服务名
func (m *Request) GetService() string {
	return m.queue.Service
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
