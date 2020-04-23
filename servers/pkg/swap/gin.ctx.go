package swap

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/utility"
)

//GinCtx gin.context
type GinCtx struct {
	*gin.Context
}

//GetBody 获取body
func (c *GinCtx) GetBody() (string, bool) {
	if body, err := ioutil.ReadAll(c.Request.Body); err == nil {
		return string(body), true
	}
	return "", false
}

//Header 设置头信息到response里
func (c *GinCtx) Header(k string, v string) {
	c.Context.Request.Response.Header.Set(k, v)
}

//GetXRequestID 获取请求编号
func (c *GinCtx) GetXRequestID() string {
	if id := c.Context.GetHeader("X-Request-Id"); id != "" {
		return id
	}
	id := utility.GetGUID()[0:9]
	c.Context.Header("X-Request-Id", id)
	c.Header("X-Request-Id", id)
	return id
}

//GetLogger 获取日志组件(不存在时根据名称创建)
func (c *GinCtx) GetLogger(name ...string) logger.ILogger {
	if len(name) == 0 {
		l, ok := c.Get("__logger__")
		if ok {
			return l.(logger.ILogger)
		}
		panic("未获取到日志组件，请先创建")
	}
	if len(name) != 2 {
		panic("创建日志组时必须传入组件名称和服务地址")
	}
	l := logger.GetSession(name[0], c.GetXRequestID(),
		"biz", strings.Replace(strings.Trim(name[1], "/"), "/", "_", -1))
	c.Set("__logger__", l)
	return l

}
