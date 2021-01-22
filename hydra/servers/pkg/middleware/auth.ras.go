package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components"
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
		input := ctx.Request().GetMap()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}
		input["__auth_"], err = auth.AuthString()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}
		respones, err := components.Def.RPC().GetRegularRPC().Request(auth.Service, input)
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, fmt.Errorf("rpc远程认证异常:err:%v", err))
			return
		}
		if !respones.IsSuccess() {
			ctx.Response().Abort(types.GetMax(respones.GetStatus(), http.StatusForbidden), fmt.Errorf("远程认证失败:%s,(%d)", respones.GetResult(), respones.GetStatus()))
			return
		}

		result := respones.GetStatus()
		if result != 200 {
			ctx.Response().Abort(http.StatusForbidden, fmt.Errorf("远程请求结果解析错误 %w", err))
			return
		}

		ctx.Request().GetMap().MergeMap(respones.GetMap())
		return
	}
}
