package middleware

import (
	"errors"
	"strings"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

//JwtAuth jwt
func JwtAuth() Handler {

	return func(ctx IMiddleContext) {

		//1. 获取jwt配置
		jwtAuth, err := ctx.APPConf().GetJWTConf()
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfError, err)
			return
		}
		if jwtAuth.Disable {
			ctx.Next()
			return
		}

		//2.检查jwt是否有效
		ctx.Response().AddSpecial("jwt")

		//3.检查是否需要跳过请求
		if ok, _ := jwtAuth.Match(ctx.Request().Path().GetRequestPath()); ok {
			ctx.Next()
			return
		}

		//4. 验证jwt
		_, err = checkJWT(ctx, jwtAuth)
		if err == nil {
			ctx.Next()
			return
		}

		//5.jwt验证失败后返回错误
		ctx.Log().Error(err)
		if jwtAuth.AuthURL != "" {
			ctx.Response().Header("Location", ctx.Request().Headers().Translate(jwtAuth.AuthURL))
			ctx.Response().Abort(xjwt.JWTStatusTokenError)
			return
		}
		ctx.Response().Abort(errs.GetCode(err, xjwt.JWTStatusTokenError), errors.New("jwt验证串错误，禁止访问"))
		return

	}
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx context.IContext, j *xjwt.JWTAuth) (data interface{}, err error) {

	//1. 从请求中获取jwt信息
	token := getToken(ctx, j)
	data, err = j.CheckJWT(token)
	if err != nil {
		return nil, err
	}

	//保存到Context中
	ctx.User().Auth().Request(data)
	return data, nil
}

//getToken 从请求头或cookie中获取cookie
func getToken(ctx context.IContext, jwt *xjwt.JWTAuth) string {
	switch strings.ToUpper(jwt.Source) {
	case xjwt.SourceHeader, xjwt.SourceHeaderShort:
		return ctx.Request().Headers().GetString(xjwt.AuthorizationHeader)
	default:
		cookie := ctx.Request().Cookies().GetString(jwt.Name)
		return cookie
	}
}
