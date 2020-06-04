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
		ras := ctx.ServerConf().GetRASConf()
		if ras.Disable {
			ctx.Next()
			return
		}

		//检查必须参数
		b, auth := ras.Contains(ctx.Request().Path().GetRouter().Path)
		if !b {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("ras")

		input, err := ctx.Request().GetData()
		if err != nil {
			ctx.Response().AbortWithError(http.StatusInternalServerError, err)
			return
		}

		input["__auth_"], err = auth.AuthString()
		if err != nil {
			ctx.Response().AbortWithError(http.StatusInternalServerError, err)
			return
		}

		respones, err := components.Def.RPC().GetRegularRPC().Request(ctx.Context(), auth.Service, input)
		if err != nil || !respones.Success() {
			ctx.Response().AbortWithError(types.GetMax(respones.Status, http.StatusForbidden), fmt.Errorf("远程认证失败:%s,err:%v(%d)", err, respones.Result, respones.Status))
			return
		}
		result, err := respones.GetResult()
		if err != nil {
			ctx.Response().AbortWithError(http.StatusForbidden, err)
			return
		}
		ctx.Meta().MergeMap(result)
		return
	}
}
