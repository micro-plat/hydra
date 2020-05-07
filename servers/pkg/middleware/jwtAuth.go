package middleware

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/context"
	xjwt "github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtAuth jwt
func JwtAuth() Handler {
	return func(ctx IMiddleContext) {

		//1. 获取jwt配置
		jwtAuth := ctx.ServerConf().GetJWTConf()
		if jwtAuth == nil || jwtAuth.Disable {
			ctx.Next()
			return
		}

		//2.检查jwt是否有效
		_, err := checkJWT(ctx, jwtAuth)
		if err == nil {
			ctx.Next()
			return
		}

		//3.检查是否需要跳过请求
		if jwtAuth.IsExcluded(ctx.Request().Path().GetService()) {
			ctx.Next()
			return
		}

		//4.jwt验证失败后返回错误
		ctx.Log().Error(err)
		ctx.Response().Abort(errs.GetCode(err, 401))
		return

	}
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx context.IContext, j *xjwt.JWTAuth) (data interface{}, err error) {

	//1. 从请求中获取jwt信息
	token := getToken(ctx, j)
	if token == "" {
		return nil, errs.NewError(403, fmt.Errorf("%s未传入jwt.token", j.Name))
	}

	//2. 解密jwt判断是否有效，是否过期
	data, er := jwt.Decrypt(token, j.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, errs.NewError(401, er)
		}
		return data, errs.NewError(403, er)
	}

	//保存到Context中
	ctx.User().SaveJwt(data)
	return data, nil
}

//getToken 从请求头或cookie中获取cookie
func getToken(ctx context.IContext, jwt *xjwt.JWTAuth) string {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		return ctx.Request().Path().GetHeader(jwt.Name)
	default:
		cookie, _ := ctx.Request().Path().GetCookie(jwt.Name)
		return cookie
	}
}
