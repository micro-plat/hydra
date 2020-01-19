package rpc

import (
	"github.com/micro-plat/lib4go/jsons"
)

//------------------------RPC响应---------------------------------------

//Response 请求结果
type Response struct {
	Status int
	Header map[string]interface{}
	Result string
}

//NewResponse 请求响应
func NewResponse(status int, p string, result string) (*Response, error) {
	res := &Response{
		Status: status,
		Result: result,
	}
	if p != "" {
		mh, err := jsons.Unmarshal([]byte(p))
		if err != nil {
			return nil, err
		}
		res.Header = mh
	}
	return res, nil
}

//NewResponseByStatus 根据状态构建响应
func NewResponseByStatus(status int, err error) (*Response, error) {
	r, _ := NewResponse(500, "", "{}")
	return r, err
}

//Success 请求是否成功
func (r *Response) Success() bool {
	return r.Status == 200
}

//GetResult 获取请求结果
func (r *Response) GetResult() (map[string]interface{}, error) {
	out, err := jsons.Unmarshal([]byte(r.Result))
	return out, err
}

//GetHeader 根据KEY获取参数
func (r *Response) GetHeader(key string) interface{} {
	return r.Header[key]
}
