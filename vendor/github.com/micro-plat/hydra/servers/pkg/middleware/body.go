package middleware

import (
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Body 处理请求的body参数
func Body() dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		if body, ok := ctx.Request.GetForm()["__body_"]; ok {
			ctx.Set("__body_", body)
		}
		ctx.Next()
	}
}
