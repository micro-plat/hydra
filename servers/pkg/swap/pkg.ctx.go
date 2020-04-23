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
