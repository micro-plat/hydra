package swap

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Handler 通用的中间件处理服务
type Handler func(IContext)

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
