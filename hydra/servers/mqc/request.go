package mqc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/lib4go/types"
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

	//将消息原串转换为map
	data := make(map[string]interface{})
	input := make(map[string]interface{})
	message := m.GetMessage()
	if err = json.Unmarshal(types.StringToBytes(messge)), &input); err != nil {
		return nil, fmt.Errorf("队列%s中存放的数据不是有效的json:%s %w", queue.Queue, m.GetMessage(), err)
	}

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
		r.header["Content-Type"] = "/json"
	}

	if v, ok := input["__data__"].([]byte); ok {
		// json.Unmarshal(v, &data)
		r.form["__body__"] = string(v)
	}

	

	// //将所有非"__""参数加到form列表
	// for k, v := range input {
	// 	if !strings.HasPrefix(k, "__") {
	// 		r.form[k] = v
	// 	}
	// }

	// //处理data数据
	// for k, v := range data {	
	// 	r.form[k] = v
	// }
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
