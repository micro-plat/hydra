package middleware

import (
	"net/http"
)

const authUserKey = "user"

//BasicAuth  http basic认证
func BasicAuth() Handler {
	return BasicAuthForRealm("")
}

//BasicAuthForRealm http basic认证
func BasicAuthForRealm(realm string) Handler {
	return func(ctx IMiddleContext) {

		basic := ctx.ServerConf().GetBasicConf()
		if basic.Disable {
			ctx.Next()
			return
		}

		//检验当前请求是否被排除
		if ok, _ := basic.Match(ctx.Request().Path().GetRouter().Path); ok {
			ctx.Next()
			return
		}

		//验证当前请求的用户名密码是否有效
		ctx.Response().AddSpecial("basic")
		if user, ok := basic.Verify(ctx.Request().Path().GetHeader("Authorization")); ok {
			ctx.User().Auth().Request(map[string]interface{}{
				authUserKey: user,
			})
			ctx.Next()
			return
		}

		ctx.Response().SetHeader("WWW-Authenticate", basic.GetRealm(realm))
		ctx.Response().Abort(http.StatusUnauthorized)
		return

	}
}
