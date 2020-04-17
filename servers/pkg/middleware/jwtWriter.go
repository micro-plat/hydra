package middleware

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

func JwtWriter(cnf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		ctx.Next()
		context := getCTX(ctx)
		if context == nil {
			return
		}
		setJwtResponse(ctx, cnf, context.Response.GetParams()["__jwt_"])

	}
}
