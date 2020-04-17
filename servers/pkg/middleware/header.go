package middleware

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Header 头设置
func Header(cnf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		ctx.Next()
		headers, ok := cnf.GetMetadata("headers").(conf.Headers)
		if !ok {
			return
		}
		for k, v := range headers {
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
