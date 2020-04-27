package middleware

import (
	"fmt"
	"strings"

	xjwt "github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/servers/pkg/swap"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtAuth jwt
func JwtAuth(f xjwt.IJWTAuth) swap.Handler {
	return func(r swap.IRequest) {

		//1. 获取jwt配置
		jwtAuth, ok := f.GetConf()
		if !ok || jwtAuth == nil || jwtAuth.Disable {
			r.Next()
			return
		}

		//2.检查jwt是否有效
		_, err := checkJWT(r, jwtAuth)
		if err == nil {
			r.Next()
			return
		}

		//3.检查是否需要跳过请求
		if jwtAuth.IsExcluded(r.GetService()) {
			r.Next()
			return
		}

		//4.jwt验证失败后返回错误
		r.GetLogger().Error(err)
		r.Abort(errs.GetCode(err, 401))
		return

	}
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(r swap.IRequest, j *xjwt.JWTAuth) (data interface{}, err error) {

	//1. 从请求中获取jwt信息
	token := getToken(r, j)
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

	//保存到请求头中
	r.Set("__jwt_", data)
	return data, nil
}

//getToken 从请求头或cookie中获取cookie
func getToken(r swap.IRequest, jwt *xjwt.JWTAuth) string {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		return r.GetHeader(jwt.Name)
	default:
		cookie, _ := r.GetCookie(jwt.Name)
		return cookie
	}
}
