package rpc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/net"
	"golang.org/x/net/context"
)

//IRequest RPC请求
type IRequest interface {
	Request(ctx context.Context, service string, form map[string]interface{}, opts ...RequestOption) (res *Response, err error)
}

type requestOption struct {
	name     string
	service  string
	headers  map[string]string
	failFast bool
	method   string
}

func (r *requestOption) getData(v interface{}) ([]byte, error) {
	buff, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	if len(buff) == 0 {
		return []byte("{}"), nil
	}
	return buff, nil
}

//RequestOption 客户端配置选项
type RequestOption func(*requestOption)

func newOption() *requestOption {
	return &requestOption{
		headers:  map[string]string{},
		failFast: true,
		method:   "GET",
	}
}

//WithHeader 请求头信息
func WithHeader(k string, v string) RequestOption {
	return func(o *requestOption) {
		o.headers[k] = v
	}
}

//WithHeaders 设置请求头
func WithHeaders(p map[string][]string) RequestOption {
	return func(o *requestOption) {
		for k, v := range p {
			o.headers[k] = v[0]
		}
	}
}

//WithHost 设置当前机器IP
func WithHost(s ...string) RequestOption {
	return func(o *requestOption) {
		if len(s) > 0 {
			o.headers["Host"] = strings.Join(s, ",")
		} else {
			o.headers["Host"] = net.GetLocalIPAddress()
		}
	}
}

//WithXRequestID 设置请求编号
func WithXRequestID(s string) RequestOption {
	return func(o *requestOption) {
		o.headers["X-Request-Id"] = s
	}
}

//WithDelay 设置请求延迟时长
func WithDelay(s time.Duration) RequestOption {
	return func(o *requestOption) {
		o.headers["X-Add-Delay"] = fmt.Sprint(s)
	}
}

//WithFailFast 快速失败
func WithFailFast(b bool) RequestOption {
	return func(o *requestOption) {
		o.failFast = b
	}
}

//WithMethod 设置请求方法
func WithMethod(m string) RequestOption {
	return func(o *requestOption) {
		o.method = strings.ToUpper(m)
	}
}

//WithContentType 设置请求类型
func WithContentType(m string) RequestOption {
	return func(o *requestOption) {
		o.headers["Content-Type"] = m
	}
}

//WithOperationName 设置请求延迟时长
func WithOperationName(name string) RequestOption {
	return func(o *requestOption) {
		o.name = name
	}
}
