package middleware

import (
	"time"

	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Logging 记录日志
func Logging(name string, serverType string) swap.Handler {
	return func(ctx swap.IContext) {

		//1. 整个服务的开始,记录请求时间与日志
		start := time.Now()
		srvs := ctx.Request().GetService()
		log := r.GetLogger(name, srvs)
		log.Info(serverType+".request:", name, r.GetMethod(), srvs, "from", r.GetClientIP())

		//2. 处理业务
		ctx.Next()

		//3. 处理响应
		if ctx.Response().GetStatusCode() >= 200 && ctx.Response().GetStatusCode() < 400 {
			log.Info(serverType+".response:", name, ctx.Request().GetMethod(), srvs, ctx.Response().GetStatusCode(), r.GetExt(), time.Since(start))
		} else {
			log.Error(serverType+".response:", name, ctx.Request().GetMethod(), srvs, ctx.Response().GetStatusCode(), r.GetExt(), time.Since(start))
		}

		//4.释放资源
		ctx.Close()
	}
}
