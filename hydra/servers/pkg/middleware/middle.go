package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
)

type imiddle interface {
	Next()
}

//IMiddleContext 中间件转换器，在context.IContext中扩展next函数
type IMiddleContext interface {
	imiddle
	context.IContext
}

//MiddleContext 中间件转换器，在context.IContext中扩展next函数
type MiddleContext struct {
	context.IContext
	imiddle
}

//newMiddleContext 构建中间件处理handler
func newMiddleContext(c context.IContext, n imiddle) IMiddleContext {
	return &MiddleContext{IContext: c, imiddle: n}
}

//Handler 通用的中间件处理服务
type Handler func(IMiddleContext)

//GinFunc 返回GIN对应的处理函数
func (h Handler) GinFunc(tps ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("__gin_middle_context__")
		if !ok {
			v = newMiddleContext(ctx.NewCtx(&ginCtx{Context: c}, tps[0]), c)
			c.Set("__gin_middle_context__", v)
		}
		h(v.(IMiddleContext))
	}
}
