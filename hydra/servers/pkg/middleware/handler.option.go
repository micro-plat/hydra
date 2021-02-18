package middleware

import (
	"net/http"
	"strings"
)

//checkOption 请求处理
func checkOption(ctx IMiddleContext) bool {
	if strings.ToUpper(ctx.Request().Path().GetMethod()) != http.MethodOptions {
		return false
	}
	ctx.Response().AddSpecial("opt")
	ctx.Response().Abort(http.StatusOK, nil)
	return true
}
