package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/security/jwt"
	"github.com/micro-plat/lib4go/types"
)

//JwtAuth jwt
func JwtAuth(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToUpper(ctx.Request.Method) == "OPTIONS" {
			ctx.Next()
			return
		}
		jwtAuth, ok := cnf.GetMetadata("jwt").(*conf.Auth)
		if !ok || jwtAuth == nil || jwtAuth.Disable {
			ctx.Next()
			return
		}

		//检查jwt.token是否正确
		data, err := checkJWT(ctx, jwtAuth)
		if err == nil {
			setJWTRaw(ctx, data)
			ctx.Next()
			return
		}

		//不需要校验的URL自动跳过
		curl := ctx.Request.URL.Path
		for _, u := range jwtAuth.Exclude {
			if u == curl {
				ctx.Next()
				return
			}
		}
		if jwtAuth.Redirect != "" {
			l, errx := url.Parse(jwtAuth.Redirect)
			if errx != nil {
				getLogger(ctx).Error(errx)
				ctx.AbortWithStatus(err.GetCode())
				return
			}
			values := l.Query()
			values.Add("redirect", ctx.Request.RequestURI)
			if l.IsAbs() {
				ctx.Redirect(301, fmt.Sprintf("%s://%s%s?%s\n", l.Scheme, l.Host, l.Path, values.Encode()))
				return
			}
			ctx.Redirect(301, fmt.Sprintf("%s?%s\n", l.Path, values.Encode()))

			return
		}
		//jwt.token错误，返回错误码
		getLogger(ctx).Error(err.GetError())
		ctx.AbortWithStatus(err.GetCode())
		return

	}
}
func setJwtResponse(ctx *gin.Context, cnf *conf.MetadataConf, data interface{}) {
	if data == nil {
		data = getJWTRaw(ctx)
	}
	if data == nil {
		return
	}
	jwtAuth, ok := cnf.GetMetadata("jwt").(*conf.Auth)
	if !ok || jwtAuth.Disable {
		return
	}
	jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
	if err != nil {
		getLogger(ctx).Errorf("jwt配置出错：%v", err)
		ctx.AbortWithStatus(500)
		return
	}
	ctx.Header("Set-Cookie", fmt.Sprintf("%s=%s;path=/;", jwtAuth.Name, jwtToken))
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx *gin.Context, auth *conf.Auth) (data interface{}, err context.IError) {
	token := getToken(ctx, auth.Name)
	if token == "" {
		return nil, context.NewError(types.ToInt(auth.FailedCode, 403), fmt.Errorf("获取%s失败或未传入该参数", auth.Name))
	}
	data, er := jwt.Decrypt(token, auth.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, context.NewError(types.ToInt(auth.FailedCode, 401), er)
		}
		return data, context.NewError(types.ToInt(auth.FailedCode, 403), er)
	}
	return data, nil
}
func getToken(ctx *gin.Context, key string) string {
	if cookie, err := ctx.Cookie(key); err == nil {
		return cookie
	}
	return ""
}
func setToken(ctx *gin.Context, name string, token string) {
	ctx.Header(name, token)
}
