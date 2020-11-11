package middleware

import (
	"fmt"
	"net/http"
)

//WhiteList 白名单
func WhiteList() Handler {
	return func(ctx IMiddleContext) {
		//获取FSA配置
		white, err := ctx.APPConf().GetWhiteListConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}

		if white.Disable {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("white")
		if !white.IsAllow(ctx.Request().Path().GetRequestPath(), ctx.User().GetClientIP()) {
			err := fmt.Errorf("白名单限制[%s]不允许访问服务[%s]", ctx.User().GetClientIP(), ctx.Request().Path().GetRequestPath())
			ctx.Response().Abort(http.StatusForbidden, err)
			return
		}
		ctx.Next()
	}
}
