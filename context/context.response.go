package context

import (
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/types"
)

//IResponse 响应
type IResponse interface {
	GetContent() interface{}
	GetCode(interface{}) int
	GetStatus() int
	SetStatus(int)
	GetView() string
	GetParams() map[string]interface{}
	SetParams(v map[string]interface{})
	SetParam(key string, v interface{})
	GetHTMLRenderContent() (int, interface{}, error)
	GetJSONRenderContent() (int, interface{}, error)
	ShouldContent(content interface{})
	MustContent(status int, content interface{})
	IsRedirect() (string, bool)
	GetHeaders() map[string]string
	SetHeader(name string, value string)
	SetHeaders(map[string]string)
	SetJWT(data interface{})
	ShouldStatus(status int)
	HasError(v interface{}) bool
	GetError() error
	IsSuccess() bool
	JSON(content interface{})
	XML(content interface{})
	Text(content interface{})
	HTML(content interface{})
}

var _ IResponse = &Response{}

type Response struct {
	Status     int
	err        error
	Content    interface{}
	Params     map[string]interface{}
	SkipHandle bool
}

func NewResponse() *Response {
	return &Response{
		Status: 0,
		Params: make(map[string]interface{}),
	}
}
func (r *Response) HasError(v interface{}) bool {
	switch v.(type) {
	case IError, error:
		return true
	}
	return false
}
func (r *Response) GetError() error {
	return r.err
}
func (r *Response) clear() {
	r.Content = nil
	r.Params = make(map[string]interface{})
	r.Status = 0
	r.err = nil
}
func (r *Response) GetStatus() int {
	return r.Status
}
func (r *Response) SetStatus(status int) {
	r.Status = types.DecodeInt(status, 0, 200, status)
}
func (r *Response) ShouldStatus(status int) {
	r.Status = types.DecodeInt(r.Status, 0, status, r.Status)
}

func (r *Response) IsSuccess() bool {
	return r.Status >= 200 && r.Status < 400
}

func (r *Response) SetJWT(data interface{}) {
	r.Params["__jwt_"] = data
}

//MakeSignAll 生成签名串，排序，并将key放到原串最后
func (r *Response) MakeSignAll(input map[string]interface{}, key string, a string, b string) (string, string) {
	values := net.NewValues()
	values.SetMap(input)
	values.Sort()
	raw := values.Join(a, b)
	return md5.Encrypt(raw + key), raw
}

//MakeSign 生成签名串
func (r *Response) MakeSign(input map[string]interface{}, key string) (string, string) {
	return r.MakeSignAll(input, key, "", "")
}
