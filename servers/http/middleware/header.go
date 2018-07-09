package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

func setHeader(cnf *conf.MetadataConf, ctx *gin.Context) {
	headers := getCrossHeader(cnf, ctx)
	for k, v := range headers {
		ctx.Header(k, strings.Join(v, ";"))
	}
}

//Header 头设置
func Header(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToUpper(ctx.Request.Method) != "OPTIONS" {
			ctx.Next()
		}
		headers, ok := cnf.GetMetadata("headers").(conf.Headers)
		if ok {
			origin := ctx.Request.Header.Get("Origin")
			for k, v := range headers {
				if k != "Access-Control-Allow-Origin" { //非跨域设置
					ctx.Header(k, v)
					continue
				}
				if origin != "" && (v == "*" || strings.Contains(v, origin)) {
					ctx.Header(k, origin)
				}
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

func getCrossHeader(cnf *conf.MetadataConf, ctx *gin.Context) http.Header {
	h := make(map[string][]string)
	origin := ctx.Request.Header.Get("Origin")
	if origin == "" {
		return nil
	}
	headers, ok := cnf.GetMetadata("headers").(conf.Headers)
	if ok {
		for k, v := range headers {
			if strings.HasPrefix(k, "Access-Control-Allow") { //非跨域设置
				if k != "Access-Control-Allow-Origin" { //非跨域设置
					h[k] = []string{v}
					continue
				}
				if v == "*" || strings.Contains(v, origin) {
					h[k] = []string{origin}
				}
				continue
			}
		}
	}
	return h
}
