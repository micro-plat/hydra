package context

import (
	"errors"
	"net/http"
	"strings"
)

//Request http.request
type httpRequest struct {
	ext map[string]interface{}
}

func (c *httpRequest) GetHeader() (map[string]string, error) {
	header := make(map[string]string)
	request, err := c.Get()
	if err != nil {
		return header, err
	}

	for k, v := range request.Header {
		header[k] = strings.Join(v, ",")
	}
	return header, nil
}
func (c *httpRequest) GetCookies() (map[string]string, error) {
	request, err := c.Get()
	if err != nil {
		return nil, err
	}
	cookies := make(map[string]string)
	for _, ck := range request.Cookies() {
		cookies[ck.Name] = ck.Value
	}
	return cookies, nil
}

//Get 获和取http.request对象
func (c *httpRequest) Get() (request *http.Request, err error) {
	r := c.ext["__func_http_request_"]
	if r == nil {
		return nil, errors.New("未找到__func_http_request_")
	}
	if f, ok := r.(*http.Request); ok {
		return f, nil
	}
	return nil, errors.New("未找到__func_http_request_传入类型错误")
}

//GetHost 获取host
func (c *httpRequest) GetHost() (string, error) {
	request, err := c.Get()
	if err != nil {
		return "", err
	}
	return request.Host, nil
}

//GetClientIP 获取客户端IP地址
func (c *httpRequest) GetClientIP() (string, error) {
	request, err := c.Get()
	if err != nil {
		return "", err
	}
	proxy := []string{}
	if ips := request.Header.Get("X-Forwarded-For"); ips != "" {
		proxy = strings.Split(ips, ",")
	}
	if len(proxy) > 0 && proxy[0] != "" {
		return proxy[0], nil
	}
	ip := strings.Split(request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0], nil
		}
	}
	return "127.0.0.1", nil
}

//GetCookie 从http.request中获取cookie
func (c *httpRequest) GetCookie(name string) (string, error) {
	request, err := c.Get()
	if err != nil {
		return "", err
	}
	cookie, err := request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
