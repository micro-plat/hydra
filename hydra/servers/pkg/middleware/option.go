package middleware

import (
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/types"
)

//Options 请求处理
func Options() Handler {
	return func(ctx IMiddleContext) {

		//是否支持服务调用
		if doOption(ctx, false) {
			return
		}
		ctx.Next()
	}
}

func doOption(ctx IMiddleContext, staticCheck bool) (isOpt bool) {
	if strings.ToUpper(ctx.Request().Path().GetMethod()) != http.MethodOptions {
		return false
	}
	ctx.Response().AddSpecial(types.DecodeString(staticCheck, true, "sopt", "opt"))
	if staticCheck {
		ctx.Response().Abort(http.StatusOK, nil)
		return true
	}

	serverType := ctx.APPConf().GetServerConf().GetServerType()
	if services.Def.Has(serverType, ctx.Request().Path().GetService(), ctx.Request().Path().GetMethod()) {
		ctx.Response().Abort(http.StatusOK, nil)
		return true
	}
	ctx.Response().Abort(http.StatusNotFound, ctx.Request().Path().GetRequestPath())
	return true
}
