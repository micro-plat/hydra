package middleware

import (
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//NoResponse 处理无响应的返回结果
func NoResponse() swap.Handler {
	return func(r swap.IContext) {
		r.Next()

		// if ctx.Writer.Written() {
		// 	return
		// }
		// writeTrace(getTrace(conf), 1, ctx, fmt.Sprint(context.Response.GetContent()))
		// ctx.Writer.WriteHeader(context.Response.GetStatus())
		// ctx.Writer.WriteString(fmt.Sprint(context.Response.GetContent()))
	}
}
