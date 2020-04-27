package swap

import (
	"strings"

	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/utility"
)

//PkgCtx   dispatcher.Context
type PkgCtx struct {
	*dispatcher.Context
}

//GetMethod 获取服务请求方式
func (c *PkgCtx) GetMethod() string {
	return c.Request.GetMethod()
}

//GetService 获取处理服务
func (c *PkgCtx) GetService() string {
	return c.Request.GetService()
}

//GetBody 获取body
func (c *PkgCtx) GetBody() (string, bool) {
	if body, ok := c.Request.GetForm()["__body_"]; ok {
		return body.(string), ok
	}
	return "", false
}

//GetXRequestID 获取请求编号
func (c *PkgCtx) GetXRequestID() string {
	if id := c.Context.GetHeader("X-Request-Id"); id != "" {
		return id
	}
	id := utility.GetGUID()[0:9]
	c.Context.Header("X-Request-Id", id)
	return id
}

//GetLogger 获取日志组件(不存在时根据名称创建)
func (c *PkgCtx) GetLogger(name ...string) logger.ILogger {
	if len(name) == 0 {
		l, ok := c.Get("__logger__")
		if ok {
			return l.(logger.ILogger)
		}
		panic("未获取到日志组件，请先创建")
	}
	l := logger.GetSession(name[0], c.GetXRequestID(),
		"biz", strings.Replace(strings.Trim(c.Context.Request.GetService(), "/"), "/", "_", -1))
	c.Set("__logger__", l)
	return l

}

//GetStatusCode 获取response状态码
func (c *PkgCtx) GetStatusCode() int {
	return c.Context.Writer.Status()
}

//GetExt 获取ext扩展信息
func (c *PkgCtx) GetExt() string {
	return c.Context.GetString("__ext__")
}

//Abort 根据错误码终止应用
func (c *PkgCtx) Abort(s int) {
	c.Context.AbortWithStatus(s)
}

//AbortWithError 根据错误码与错误消息终止应用
func (c *PkgCtx) AbortWithError(s int, err error) {
	c.Context.AbortWithError(s, err)
}

//Close 关闭并释放所有资源
func (c *PkgCtx) Close() {

}

//GetCookie 获取cookie信息
func (c *PkgCtx) GetCookie(name string) (string, bool) {
	return "", false
}
