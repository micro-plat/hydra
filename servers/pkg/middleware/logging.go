package middleware

import (
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/lib4go/logger"
)

//Logging 记录日志
func Logging(conf conf.IMetadata) Handler {
	return func(ctx IRequest) {
		start := time.Now()
		setStartTime(ctx)
		p := ctx.Request.GetService()
		uuid := getUUID(ctx)
		setUUID(ctx, uuid)
		log := logger.GetSession(conf.Name, uuid, "biz", strings.Replace(strings.Trim(ctx.Request.GetService(), "/"), "/", "_", -1))
		log.Info(conf.Type+".request:", conf.Name, ctx.Request.GetMethod(), p, "from", ctx.ClientIP())
		setLogger(ctx, log)
		ctx.Next()

		v, _ := getResponseRaw(ctx)
		statusCode := ctx.Writer.Status()
		if statusCode >= 200 && statusCode < 400 {
			log.Info(conf.Type+".response:", conf.Name, ctx.Request.GetMethod(), p, statusCode, getExt(ctx), time.Since(start), v)
		} else {
			log.Error(conf.Type+".response:", conf.Name, ctx.Request.GetMethod(), p, statusCode, getExt(ctx), time.Since(start), v)
		}
	}
}
