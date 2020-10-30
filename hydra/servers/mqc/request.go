package mqc

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
)

const DefMethod = "GET"

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
	r = &Request{
		IMQCMessage: m,
		queue:       queue,
		method:      DefMethod,
		form:        make(map[string]interface{}),
		header:      make(map[string]string),
	}
	if err = json.Unmarshal([]byte(m.GetMessage()), &r.form); err != nil {
		return nil, fmt.Errorf("队列%s中存放的数据不是有效的json:%s %w", queue.Queue, m.GetMessage(), err)
	}
	if v, ok := r.form["__header__"].(map[string]interface{}); ok {
		for n, m := range v {
			r.header[n] = fmt.Sprint(m)
		}
	}
	r.header["Client-IP"] = "127.0.0.1"
	r.form["__body_"] = m.GetMessage()
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
