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

		statusCode := ctx.Writer.Status()
		if statusCode >= 200 && statusCode < 400 {
			log.Info(conf.Type+".response", ctx.Request.Method, p, statusCode, getExt(ctx), time.Since(start))
		} else {
			log.Error(conf.Type+".response", ctx.Request.Method, p, statusCode, getExt(ctx), time.Since(start))
		}
		context := getCTX(ctx)
		if context != nil {
			defer context.Close()
		}
	}

}
