package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

//Header 头设置
func Header(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToUpper(ctx.Request.Method) != "OPTIONS" {
			ctx.Next()
		}
		headers, ok := cnf.GetMetadata("headers").(conf.Headers)
		if !ok {
			return
		}
		origin := ctx.Request.Header.Get("Origin")
		for k, v := range headers {
			if k != "Access-Control-Allow-Origin" { //非跨域设置
				ctx.Header(k, v)
				continue
			}
			if origin != "" && strings.Contains(v, origin) {
				ctx.Header(k, origin)
			}
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
