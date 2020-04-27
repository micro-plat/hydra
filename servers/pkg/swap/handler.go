package swap

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/logger"
)

//IRequest 用于中间件处理的上下文管理
type IRequest interface {
	GetMethod() string
	GetService() string
	GetBody() (string, bool)

	GetClientIP() string
	GetStatusCode() int
	GetExt() string
	GetHeader(string) string
	Header(string, string)
	GetCookie(string) (string, bool)
	GetLogger(name ...string) logger.ILogger
	conf.IMetadata
	Next()
	Abort(int)
	AbortWithError(int, error)
	Close()
}

//Handler 通用的中间件处理服务
type Handler func(IRequest)

//PkgFunc 返回当前包对应的处理函数
func (h Handler) PkgFunc() dispatcher.HandlerFunc {
	return func(c *dispatcher.Context) {
		var ctx = &PkgCtx{Context: c}
		h(ctx)
	}
}

//GinFunc 返回GIN对应的处理函数
func (h Handler) GinFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = &GinCtx{Context: c}
		h(ctx)
	}
}
