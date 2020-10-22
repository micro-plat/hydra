package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtAuth jwt
func JwtAuth() Handler {

	return func(ctx IMiddleContext) {

		//1. 获取jwt配置
		jwtAuth, err := ctx.ServerConf().GetJWTConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if jwtAuth.Disable {
			ctx.Next()
			return
		}

		//2.检查jwt是否有效
		ctx.Response().AddSpecial("jwt")

		//3.检查是否需要跳过请求
		if ok, _ := jwtAuth.Match(ctx.Request().Path().GetRequestPath(), "/"); ok {
			ctx.Next()
			return
		}

		routerObj, err := ctx.Request().Path().GetRouter()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if ok, _ := jwtAuth.Match(routerObj.Service, "/"); ok {
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
		ctx.Response().Abort(errs.GetCode(err, http.StatusForbidden), errors.New("jwt验证串错误，禁止访问"))
		return

	}
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx context.IContext, j *xjwt.JWTAuth) (data interface{}, err error) {

	//1. 从请求中获取jwt信息
	token := getToken(ctx, j)
	if token == "" {
		return nil, errs.NewError(http.StatusUnauthorized, fmt.Errorf("未传入jwt.token(%s %s值为空)", j.Source, j.Name))
	}
	//2. 解密jwt判断是否有效，是否过期
	data, er := jwt.Decrypt(token, j.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, errs.NewError(http.StatusForbidden, er)
		}
		return data, errs.NewError(http.StatusForbidden, fmt.Errorf("jwt.token值(%s)有误 %w", token, er))
	}

	//保存到Context中
	ctx.User().Auth().Request(data)
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
