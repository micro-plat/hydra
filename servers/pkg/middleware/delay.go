package middleware

import (
	"time"

	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/types"
)

//DELAY_HEADER_KEY  延时请求头名称
const DELAY_HEADER_KEY = "X-Add-Delay"

//Delay 处理请求的延时时长
func Delay() dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		if delay := types.GetInt64(ctx.Request.GetHeader()[DELAY_HEADER_KEY], 0); delay > 0 {
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
		ctx.Next()
	}
}
