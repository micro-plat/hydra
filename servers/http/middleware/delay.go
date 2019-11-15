package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/lib4go/types"
)

//DELAY_HEADER_KEY  延时请求头名称
const DELAY_HEADER_KEY = "X-Add-Delay"

//Delay 处理请求的延时时长
func Delay() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if delay := types.GetInt64(ctx.Request.Header.Get(DELAY_HEADER_KEY), 0); delay > 0 {
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
		ctx.Next()
	}
}
