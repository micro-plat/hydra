package middleware

import (
	"errors"
	"net/http"

	"github.com/micro-plat/hydra/context"
)

//BasicAuth  http basic认证
func BasicAuth() Handler {
	return BasicAuthForRealm()
}

//BasicAuthForRealm http basic认证
func BasicAuthForRealm() Handler {
	return func(ctx IMiddleContext) {

		basic, err := ctx.APPConf().GetBasicConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}

		if basic.Disable {
			ctx.Next()
			return
		}
		//检验当前请求是否被排除
		if ok, _ := basic.Match(ctx.Request().Path().GetRequestPath()); ok {
			ctx.Next()
			return
		}

		//验证当前请求的用户名密码是否有效
		ctx.Response().AddSpecial("basic")
		if user, ok := basic.Verify(ctx.Request().Headers().GetString("Authorization"), ctx.Invoke); ok {
			ctx.Meta().SetValue(context.UserName, user)
			ctx.User().Auth().Request(map[string]interface{}{
				context.UserName: user,
			})
			ctx.Next()
			return
		}

		ctx.Response().Header("WWW-Authenticate", basic.GetRealm())
		ctx.Response().Abort(http.StatusUnauthorized, errors.New("未提供验证信息"))
		return

	}
}
