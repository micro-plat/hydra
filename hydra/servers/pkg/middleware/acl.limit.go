package middleware

import (
	"net/http"
	"time"
)

//Limit 服务器限流配置
func Limit() Handler {
	return func(ctx IMiddleContext) {

		//获取限流器
		limiter, err := ctx.APPConf().GetLimiterConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if limiter.Disable {
			ctx.Next()
			return
		}

		//判断请求是否指定限流规则
		enable, rule := limiter.GetLimiter(ctx.Request().Path().GetRequestPath())
		if !enable {
			ctx.Next()
			return
		}

		//获取执行令牌
		res := rule.GetLimiter().Reserve()

		//判断请求是否需要进行延迟处理
		delay := res.Delay()
		if delay <= 0 {
			ctx.Next()
			return
		}
		//当前请求被限流
		ctx.Response().AddSpecial("limit")
		wait := rule.GetDelay()
		if delay > wait { //当前请求将被限流，根据配置进行降级或结果输出处理
			res.Cancel()
			ctx.Request().Path().Limit(true, rule.Fallback)
			s, c := rule.GetResponse()
			ctx.Response().Write(s, c)
			ctx.Next()
			return
		}

		//等待一定时间后继续处理
		time.Sleep(delay)
		ctx.Next()
	}
}
