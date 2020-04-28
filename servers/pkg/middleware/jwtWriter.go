package middleware

import (
	"fmt"
	"strings"
	"time"

	xjwt "github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/servers/pkg/swap"
	"github.com/qxnw/lib4go/security/jwt"
)

//JwtWriter 将jwt信息写入到请求中
func JwtWriter(f xjwt.IJWTAuth) swap.Handler {
	return func(ctx swap.IContext) {
		ctx.Next()
		conf, ok := f.GetConf()
		if !ok || conf.Disable {
			return
		}
		setJwtResponse(ctx, conf, ctx.GetResponseParam()["__jwt_"])
	}
}

func setJwtResponse(ctx swap.IContext, jwtAuth *xjwt.JWTAuth, data interface{}) {
	if data == nil {
		data, _ = ctx.Get("__jwt_")
	}
	if data == nil {
		return
	}
	jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
	if err != nil {
		ctx.Response().AbortWithError(500, fmt.Errorf("jwt配置出错：%v", err))
		return
	}
	setToken(r, jwtAuth, jwtToken)
}

//setToken 设置jwt到响应头或cookie中
func setToken(ctx swap.IContext, jwt *xjwt.JWTAuth, token string) {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		ctx.Header(jwt.Name, token)
	default:
		expireTime := time.Now().Add(time.Duration(time.Duration(jwt.ExpireAt)*time.Second - 8*60*60*time.Second))
		expireVal := expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")

		if jwt.Domain != "" {
			ctx.Response().Header("Set-Cookie", fmt.Sprintf("%s=%s;domain=%s;path=/;expires=%s;", jwt.Name, token, jwt.Domain, expireVal))
			return
		}
		ctx.Response().Header("Set-Cookie", fmt.Sprintf("%s=%s;path=/;expires=%s;", jwt.Name, token, expireVal))
	}
}
