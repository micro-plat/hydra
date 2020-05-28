package ctx

import (
	"strings"

	"github.com/micro-plat/hydra/context"
)

var _ context.IPath = &rpath{}

//rpath 处理请求的路径信息
type rpath struct {
	ctx context.IInnerContext
}

//GetMethod 获取服务请求方式
func (c *rpath) GetMethod() string {
	return c.ctx.GetMethod()
}

//GetService 获取处理的服务名
func (c *rpath) GetService() string {
	return ""
}

//GetURL 获取请求路径
func (c *rpath) GetURL() string {
	return c.ctx.GetURL().String()
}

//GetPath 获取请求路径
func (c *rpath) GetPath() string {
	return c.ctx.GetURL().Path
}

//GetHeader 获取请求头信息
func (c *rpath) GetHeader(key string) string {
	return strings.Join(c.GetHeaders()[key], ",")
}

//GetHeaders 获取请求的header
func (c *rpath) GetHeaders() map[string][]string {
	return c.ctx.GetHeaders()
}

//GetHeaders 获取请求的header
func (c *rpath) GetCookies() map[string]string {
	out := make(map[string]string)
	cookies := c.ctx.GetCookies()
	for _, cookie := range cookies {
		out[cookie.Name] = cookie.Value
	}
	return out
}

//GetCookie 获取cookie信息
func (c *rpath) GetCookie(name string) (string, bool) {
	if cookie, ok := c.GetCookies()[name]; ok {
		return cookie, true
	}
	return "", false
}

//GetCookie 获取cookie信息
func (c *rpath) getCookie(name string) string {
	m, _ := c.GetCookie(name)
	return m
}
