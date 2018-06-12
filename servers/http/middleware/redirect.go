package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

func Redirect(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		context := getCTX(ctx)
		if context == nil {
			return
		}
		// //处理跳转3xx
		// if url, ok := context.Response.IsRedirect(); ok {
		// 	defer context.Close()
		// 	ctx.Redirect(context.Response.GetStatus(), url)
		// 	return
		// }
	}
}
