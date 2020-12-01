package middleware

import (
	"time"
)

//xAddDelay  延时请求头名称
const xAddDelay = "X-Add-Delay"

//Delay 处理请求的延时时长
func Delay() Handler {
	return func(ctx IMiddleContext) {
		if delay := ctx.Request().Headers().GetString(xAddDelay); delay != "" {
			ctx.Response().AddSpecial("delay")
			delayDuration, err := time.ParseDuration(delay)
			if err != nil {
				ctx.Log().Errorf("%s 的值%s有误 %v", xAddDelay, delay, err)
			} else {
				time.Sleep(delayDuration)
			}
		}
		ctx.Next()
	}
}
