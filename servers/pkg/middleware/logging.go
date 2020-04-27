package middleware

import (
	"time"

	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Logging 记录日志
func Logging(name string, serverType string) swap.Handler {
	return func(r swap.IRequest) {

		//整个服务的开始
		start := time.Now()
		srvs := r.GetService()
		log := r.GetLogger(name, srvs)
		log.Info(serverType+".request:", name, r.GetMethod(), srvs, "from", r.GetClientIP())

		r.Next()

		if r.GetStatusCode() >= 200 && r.GetStatusCode() < 400 {
			log.Info(serverType+".response:", name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start))
		} else {
			log.Error(serverType+".response:", name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start))
		}

		//处理整个请求的关闭
		r.Close()
	}
}
