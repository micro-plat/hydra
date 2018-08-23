package middleware

import (
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/logger"
)

//Logging 记录日志
func Logging(conf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		start := time.Now()
		setStartTime(ctx)
		p := ctx.Request.GetService()
		uuid := getUUID(ctx)
		setUUID(ctx, uuid)
		log := logger.GetSession(conf.Name, uuid)
		log.Info(conf.Type+".request:", conf.Name, ctx.Request.GetMethod(), p, "from", ctx.ClientIP())
		setLogger(ctx, log)
		ctx.Next()

		v, _ := getResponseRaw(ctx)
		statusCode := ctx.Writer.Status()
		if statusCode >= 200 && statusCode < 400 {
			log.Info(conf.Type+".response:", conf.Name, ctx.Request.GetMethod(), p, statusCode, time.Since(start), v)
		} else {
			log.Error(conf.Type+".response:", conf.Name, ctx.Request.GetMethod(), p, statusCode, time.Since(start), v)
		}
	}

}
