package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

func JwtWriter(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		context := getCTX(ctx)
		if context == nil {
			return
		}
		fmt.Println("JwtWriter:1")
		setJwtResponse(ctx, cnf, context.Response.GetParams()["__jwt_"])
	}
}
