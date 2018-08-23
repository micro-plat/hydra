package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/logger"
)

//Logging 记录日志
func Logging(conf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		setStartTime(ctx)
		setMetadataConf(ctx, conf)
		p := ctx.Request.URL.Path
		if ctx.Request.URL.RawQuery != "" {
			p = p + "?" + ctx.Request.URL.RawQuery
		}
		uuid := getUUID(ctx)
		setUUID(ctx, uuid)
		log := logger.GetSession(conf.Name, uuid)
		log.Info(conf.Type+".request", ctx.Request.Method, p, "from", ctx.ClientIP())
		setLogger(ctx, log)
		ctx.Next()

		v, _ := getResponseRaw(ctx)

		statusCode := ctx.Writer.Status()
		if statusCode >= 200 && statusCode < 400 {
			log.Info(conf.Type+".response", ctx.Request.Method, p, statusCode, getExt(ctx), time.Since(start), v)
		} else {
			log.Error(conf.Type+".response", ctx.Request.Method, p, statusCode, getExt(ctx), time.Since(start), v)
		}

		context := getCTX(ctx)
		if context != nil {
			defer context.Close()
		}
	}
}

func wLogHead(ctx *gin.Context, p string) {
	conf := getMetadataConf(ctx)
	getLogger(ctx).Info(conf.Type+".request", ctx.Request.Method, p, "from", ctx.ClientIP())
}
func wLogTail(ctx *gin.Context, p string, start time.Time) {
	conf := getMetadataConf(ctx)
	statusCode := getCTX(ctx).Response.GetStatus()
	if statusCode >= 200 && statusCode < 400 {
		getLogger(ctx).Info(conf.Type+".response", ctx.Request.Method, p, statusCode, time.Since(start))
	} else {
		getLogger(ctx).Error(conf.Type+".response", ctx.Request.Method, p, statusCode, time.Since(start))
	}
}
