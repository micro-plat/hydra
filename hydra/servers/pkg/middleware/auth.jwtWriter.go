package middleware

import (
	"fmt"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtWriter 将jwt信息写入到请求中
func JwtWriter() Handler {
	return func(ctx IMiddleContext) {
		ctx.Next()
		conf, err := ctx.APPConf().GetJWTConf()
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfError, err)
			return
		}

		if conf.Disable {
			return
		}
		setJwtResponse(ctx, conf, ctx.User().Auth().Response())
	}
}

func setJwtResponse(ctx IMiddleContext, jwtAuth *xjwt.JWTAuth, data interface{}) {

	//清除jwt认证信息
	if ctx.ClearAuth() {
		if k, v, ok := jwtAuth.GetJWTForRspns("", true); ok {
			ctx.Response().Header(k, v)
		}
		return
	}

	//写入响应
	if data != nil {
		jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfDataError, fmt.Errorf("jwt配置出错：%v", err))
			return
		}
		if k, v, ok := jwtAuth.GetJWTForRspns(jwtToken); ok {
			ctx.Response().Header(k, v)
		}
	}
	return

}
