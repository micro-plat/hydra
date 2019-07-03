package http

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/micro-plat/lib4go/encoding"
)

// NewRequest 创建新请求
func (c *HTTPClient) NewRequest(method string, url string, args ...string) *HTTPClientRequest {
	request := &HTTPClientRequest{}
	request.client = c.client
	request.headers = make(map[string]string)
	request.method = strings.ToUpper(method)
	request.params = ""
	request.url = url
	request.encoding = getEncoding(args...)
	return request
}

// SetData 设置请求参数
func (c *HTTPClientRequest) SetData(params string) {
	c.params = params
}

// SetHeader 设置http header
func (c *HTTPClientRequest) SetHeader(key string, value string) {
	c.headers[key] = value
}

// Request 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
// header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClientRequest) Request() (content string, status int, err error) {
	req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.params))
	if err != nil {
		return
	}
	req.Close = true
	for i, v := range c.headers {
		req.Header.Set(i, v)
	}
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	status = resp.StatusCode
	rc, err := encoding.DecodeBytes(body, c.encoding)
	content = string(rc)
	return
}
