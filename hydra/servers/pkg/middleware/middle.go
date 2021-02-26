package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type imiddle interface {
	Next()
	GetRouterPath() string
	Find(path string) bool
	Service(string)
	ClearAuth(c ...bool) bool
	ServicePrefix(string)
}

//IMiddleContext 中间件转换器，在context.IContext中扩展next函数
type IMiddleContext interface {
	imiddle
	context.IContext
	Trace(...interface{})
}

//MiddleContext 中间件转换器，在context.IContext中扩展next函数
type MiddleContext struct {
	context.IContext
	imiddle
}

//Trace 输出调试日志
func (m *MiddleContext) Trace(s ...interface{}) {
	if global.IsDebug {
		m.IContext.Log().Debug(s...)
	}
}

//NewMiddleContext 构建中间件处理handler
func NewMiddleContext(c context.IContext, n imiddle) IMiddleContext {
	return &MiddleContext{IContext: c, imiddle: n}
}

//Handler 通用的中间件处理服务
type Handler func(IMiddleContext)

//GinFunc 返回GIN对应的处理函数
func (h Handler) GinFunc(tps ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("__middle_context__")
		if !ok {
			if err := recover(); err != nil {
				global.Def.Log().Errorf("-----[Recovery] panic recovered:\n%s\n%s 构建context出现错误", err, global.GetStack())
				c.AbortWithError(http.StatusNotExtended, fmt.Errorf("%v", "Server Error"))
			}
			rawCtx := &ginCtx{Context: c}
			nctx := ctx.NewCtx(rawCtx, tps[0])
			nctx.Meta().SetValue("__context_", c)
			v = NewMiddleContext(nctx, rawCtx)
			c.Set("__middle_context__", v)
		}
		h(v.(IMiddleContext))
	}
}

//DispFunc 返回disp对应的处理函数
func (h Handler) DispFunc(tps ...string) dispatcher.HandlerFunc {
	return func(c *dispatcher.Context) {
		v, ok := c.Get("__middle_context__")
		if !ok {
			if err := recover(); err != nil {
				global.Def.Log().Errorf("-----[Recovery] panic recovered:\n%s\n%s 构建context出现错误", err, global.GetStack())
				c.AbortWithError(http.StatusNotExtended, fmt.Errorf("%v", "Server Error"))
			}
			rawCtx := &dispCtx{Context: c}
			nctx := ctx.NewCtx(rawCtx, tps[0])
			nctx.Meta().SetValue("__context_", c)
			v = NewMiddleContext(nctx, rawCtx)
			c.Set("__middle_context__", v)
		}
		h(v.(IMiddleContext))
	}
}
