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

//GetMethod 获取服务请求方式
func (c *GinCtx) GetMethod() string {
	return c.Request.Method
}

//GetService 获取处理服务
func (c *GinCtx) GetService() string {
	return ""
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

//GetClientIP 获取客户端IP地址
func (c *GinCtx) GetClientIP() string {
	return c.Context.ClientIP()
}

//GetStatusCode 获取response状态码
func (c *GinCtx) GetStatusCode() int {
	return c.Context.Writer.Status()
}

//GetExt 获取ext扩展信息
func (c *GinCtx) GetExt() string {
	return c.Context.GetString("__ext__")
}

//Abort 根据错误码终止应用
func (c *GinCtx) Abort(s int) {
	c.Context.AbortWithStatus(s)
}

//AbortWithError 根据错误码与错误消息终止应用
func (c *GinCtx) AbortWithError(s int, err error) {
	c.Context.AbortWithError(s, err)
}

//Close 关闭并释放所有资源
func (c *GinCtx) Close() {

}

//GetCookie 获取cookie信息
func (c *GinCtx) GetCookie(name string) (string, bool) {
	if cookie, err := c.Context.Request.Cookie(name); err == nil {
		return cookie.Value, true
	}
	return "", false
}

//File 输入文件
func (c *GinCtx) File(f string) {
	c.Context.File(f)
}

//GetRequestPath 获取请求路径
func (c *GinCtx) GetRequestPath() string {
	return c.Context.Request.URL.Path
}

//GetResponseParam 获取响应参数
func (c *GinCtx) GetResponseParam() map[string]interface{} {
	return nil
}
