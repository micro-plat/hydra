package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type imiddle interface {
	Next()
}

//IMiddleContext 中间件转换器，在context.IContext中扩展next函数
type IMiddleContext interface {
	imiddle
	context.IContext
	Trace(...interface{})
	GetHttpReqResp() (*http.Request, http.ResponseWriter)
}

//MiddleContext 中间件转换器，在context.IContext中扩展next函数

type MiddleContext struct {
	context.IContext
	imiddle
	req  *http.Request
	resp http.ResponseWriter
}

//Trace 输出调试日志
func (m *MiddleContext) Trace(s ...interface{}) {
	if global.IsDebug {
		m.IContext.Log().Debug(s...)
	}
}

//GetHttpReqResp 获取http请求与响应对象
func (m *MiddleContext) GetHttpReqResp() (*http.Request, http.ResponseWriter) {
	return m.req, m.resp
}

//newMiddleContext 构建中间件处理handler
func newMiddleContext(c context.IContext, n imiddle, req *http.Request, resp http.ResponseWriter) IMiddleContext {
	return &MiddleContext{IContext: c, imiddle: n, req: req, resp: resp}
}

//Handler 通用的中间件处理服务
type Handler func(IMiddleContext)

//GinFunc 返回GIN对应的处理函数
func (h Handler) GinFunc(tps ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("__middle_context__")
		if !ok {
			nctx := ctx.NewCtx(&ginCtx{Context: c}, tps[0])
			nctx.Meta().SetValue("__context_", c)
			v = newMiddleContext(nctx, c, c.Request, c.Writer)
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
			nctx := ctx.NewCtx(&dispCtx{Context: c}, tps[0])
			nctx.Meta().SetValue("__context_", c)
			v = newMiddleContext(nctx, c, nil, nil)
			c.Set("__middle_context__", v)
		}
		h(v.(IMiddleContext))
	}
}

func AddMiddlewareHook(handlers []Handler, callback func(Handler)) {
	for i := range handlers {
		callback(handlers[i])
	}
}
