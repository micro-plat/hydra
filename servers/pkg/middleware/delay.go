package middleware

import (
	"time"

	"github.com/micro-plat/lib4go/types"
)

//xAddDelay  延时请求头名称
const xAddDelay = "X-Add-Delay"

//Delay 处理请求的延时时长
func Delay() Handler {
	return func(ctx IMiddleContext) {
		if delay := types.GetInt64(ctx.Request().Path().GetHeader(xAddDelay), 0); delay > 0 {
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
		ctx.Next()
	}
}
