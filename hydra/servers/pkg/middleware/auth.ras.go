package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/lib4go/types"
)

//RASAuth RAS远程认证
func RASAuth() Handler {
	return func(ctx IMiddleContext) {

		//获取FSA配置
		ras, err := ctx.APPConf().GetRASConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}

		if ras.Disable {
			ctx.Next()
			return
		}

		b, auth := ras.Match(ctx.Request().Path().GetRequestPath())
		if !b {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("ras")
		err = ctx.Request().GetError()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}
		authStr, err := auth.AuthString()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}
		ctx.Request().SetValue("__auth_", authStr)

		respones := ctx.Invoke(auth.Service)
		if !respones.IsSuccess() {
			ctx.Response().Abort(types.GetMax(respones.GetStatus(), http.StatusForbidden), fmt.Errorf("RAS认证失败:%s,(%d)", respones.GetResult(), respones.GetStatus()))
			return
		}
		result := respones.GetStatus()
		if result != 200 {
			ctx.Response().Abort(http.StatusForbidden, fmt.Errorf("RAS认证请求结果解析错误 %w", err))
			return
		}

		ctx.Request().GetMap().MergeMap(respones.GetMap())
	}
}
