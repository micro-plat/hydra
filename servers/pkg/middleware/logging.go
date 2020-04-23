package middleware

import (
	"time"

	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Logging 记录日志
func Logging(conf *conf.Metadata) swap.Handler {
	return func(r swap.IRequest) {

		start := time.Now()
		srvs := r.GetService()
		log := r.GetLogger(conf.Name, srvs)
		log.Info(conf.Type+".request:", conf.Name, r.GetMethod(), srvs, "from", r.GetClientIP())

		r.Next()

		if r.GetStatusCode() >= 200 && r.GetStatusCode() < 400 {
			log.Info(conf.Type+".response:", conf.Name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start), v)
		} else {
			log.Error(conf.Type+".response:", conf.Name, r.GetMethod(), srvs, r.GetStatusCode(), r.GetExt(), time.Since(start), v)
		}
	}
}
