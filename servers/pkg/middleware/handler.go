package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//IRequest 用于中间件处理的上下文管理
type IRequest interface {
	GetBody() (string, bool)
	GetHeader(string) string
	Header(string, string)
	conf.IMetadata
	Next()
}

//Handler 通用的中间件处理服务
type Handler func(IRequest)

//PkgFunc 返回当前包对应的处理函数
func (h Handler) PkgFunc() dispatcher.HandlerFunc {
	return func(c *dispatcher.Context) {
		var ctx = &swap.PkgCtx{Context: c}
		h(ctx)
	}
}

//GinFunc 返回GIN对应的处理函数
func (h Handler) GinFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = &swap.GinCtx{Context: c}
		h(ctx)
	}
}
