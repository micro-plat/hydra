package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

//Header 头设置
func Header(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		headers, ok := cnf.GetMetadata("headers").(conf.Headers)
		if !ok {
			return
		}
		for k, v := range headers {
			if strings.Contains(v, ",") && k == "Access-Control-Allow-Origin" {
				if strings.Contains(v, ctx.Request.Host) {
					hosts := strings.Split(v, ",")
					for _, h := range hosts {
						if strings.Contains(h, ctx.Request.Host) {
							ctx.Header(k, h)
							continue
						}
					}
				}
				continue
			}
			ctx.Header(k, v)
		}
		context := getCTX(ctx)
		if context == nil {
			return
		}
		header := context.Response.GetHeaders()
		for k, v := range header {
			ctx.Header(k, v)
		}
	}
}
