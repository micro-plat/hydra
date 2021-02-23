package middleware

import (
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/services"
)

//Options 请求处理
func Options() Handler {
	return func(ctx IMiddleContext) {

		//是否支持服务调用
		if doOption(ctx, services.Def.Has(ctx.APPConf().GetServerConf().GetServerType(),
			ctx.FullPath(),
			ctx.Request().Path().GetMethod())) {
			return
		}
		ctx.Next()
		return
	}
}
func doOption(ctx IMiddleContext, v bool) bool {
	if strings.ToUpper(ctx.Request().Path().GetMethod()) != http.MethodOptions {
		return false
	}
	if v {
		ctx.Response().AddSpecial("opt")
		ctx.Response().Abort(http.StatusOK, nil)
		return true
	}
	return false
}
