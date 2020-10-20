package middleware

import (
	"fmt"
	"net/http"
)

//BlackList 黑名单
func BlackList() Handler {
	return func(ctx IMiddleContext) {
		white := ctx.ServerConf().GetBlackListConf()
		if white.Disable {
			ctx.Next()
			return
		}
		ctx.Response().AddSpecial("black")
		if white.IsDeny(ctx.User().GetClientIP()) {
			err := fmt.Errorf("黑名单限制[%s]不允许访问", ctx.User().GetClientIP())
			ctx.Response().Abort(http.StatusForbidden, err)
			return
		}
		ctx.Next()
	}
}
