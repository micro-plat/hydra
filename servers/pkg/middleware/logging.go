package middleware

import (
	"time"

	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Logging 记录日志
func Logging(name string, serverType string) swap.Handler {
	return func(r swap.IRequest) {

		//1. 整个服务的开始,记录请求时间与日志
		start := time.Now()
		srvs := r.GetService()
		log := r.GetLogger(name, srvs)
		log.Info(serverType+".request:", name, r.GetMethod(), srvs, "from", r.GetClientIP())

		//2. 处理业务
		r.Next()

		//3. 处理响应
		if r.GetStatusCode() >= 200 && r.GetStatusCode() < 400 {
			log.Info(serverType+".response:", name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start))
		} else {
			log.Error(serverType+".response:", name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start))
		}

		//4.释放资源
		r.Close()
	}
}
