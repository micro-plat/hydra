package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

//AjaxRequest ajax请求限制
func AjaxRequest(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		b, ok := cnf.GetMetadata("ajax-request").(bool)
		if ok && b && ctx.GetHeader("X-Requested-With") != "XMLHttpRequest" {
			ctx.AbortWithStatus(403)
			return
		}
		ctx.Next()
		return
	}
}
