package middleware

import (
	"net/http"
	"strings"
)

//Options 请求处理
func Options() Handler {
	return func(ctx IMiddleContext) {
		//options请求则自动不再进行后续处理
		if strings.ToUpper(ctx.Request().Path().GetMethod()) == http.MethodOptions {
			ctx.Response().AddSpecial("opt")
			ctx.Response().Abort(http.StatusOK, nil)
			return
		}
		ctx.Next()

	}
}
