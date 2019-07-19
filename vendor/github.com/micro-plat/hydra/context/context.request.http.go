package context

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//Request http.request
type httpRequest struct {
	ext map[string]interface{}
}

func (c *httpRequest) Clear() {
	c.ext = nil
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

func (c *httpRequest) GetResponse() (response gin.ResponseWriter, err error) {
	r := c.ext["__func_http_response_"]
	if r == nil {
		return nil, errors.New("未找到__func_http_response_")
	}
	if f, ok := r.(gin.ResponseWriter); ok {
		return f, nil
	}
	return nil, errors.New("未找到__func_http_response_传入类型错误")
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

//GetImageExt 根据content-type获取文件扩展名
func (w *httpRequest) GetImageExt() (string, error) {
	header, err := w.GetHeader()
	if err != nil {
		return "", err
	}
	ct := header["Content-Type"]
	if !strings.HasPrefix(ct, "image/") {
		return "", fmt.Errorf("Content-Type:%s不是图片格式", ct)
	}
	switch ct {
	case "image/x-icon":
		return ".ico", nil
	case "image/pnetvue":
		return ".net", nil
	case "vnd.rn-realpix":
		return ".rp", nil
	default:
		imgs := strings.Split(ct, "/")
		if len(imgs) < 2 {
			return "", fmt.Errorf("Content-Type:%s不是图片格式", ct)
		}
		g := strings.Split(imgs[1], ".")
		return fmt.Sprintf(".%s", g[len(g)-1]), nil
	}
}
