package middleware

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//NoResponse 处理无响应的返回结果
func NoResponse(conf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		ctx.Next()
		context := getCTX(ctx)
		if context == nil {
			return
		}
		defer context.Close()
		if ctx.Writer.Written() {
			return
		}
		writeTrace(getTrace(conf), 1, ctx, fmt.Sprint(context.Response.GetContent()))
		ctx.Writer.WriteHeader(context.Response.GetStatus())
		ctx.Writer.WriteString(fmt.Sprint(context.Response.GetContent()))
	}
}
