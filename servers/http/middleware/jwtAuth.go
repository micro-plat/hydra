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
		if jwtAuth.Redirect != "" && strings.ToUpper(ctx.Request.Method) == "GET" {
			l, errx := url.Parse(jwtAuth.Redirect)
			if errx != nil {
				getLogger(ctx).Error(errx)
				setHeader(cnf, ctx)
				ctx.AbortWithStatus(err.GetCode())
				return
			}
			values := l.Query()
			values.Add("redirect", ctx.Request.RequestURI)
			if l.IsAbs() {
				ctx.Redirect(302, fmt.Sprintf("%s://%s%s?%s\n", l.Scheme, l.Host, l.Path, values.Encode()))
				setHeader(cnf, ctx)
				ctx.Abort()
				return
			}
			ctx.Redirect(302, fmt.Sprintf("%s?%s\n", l.Path, values.Encode()))
			setHeader(cnf, ctx)
			ctx.Abort()
			return
		}
		//jwt.token错误，返回错误码
		getLogger(ctx).Error(err.GetError())
		setHeader(cnf, ctx)
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
	setToken(ctx, jwtAuth, jwtToken)
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx *gin.Context, auth *conf.Auth) (data interface{}, err context.IError) {
	token := getToken(ctx, auth)
	if token == "" {
		return nil, context.NewError(types.GetInt(auth.FailedCode, 403), fmt.Errorf("获取%s失败或未传入该参数", auth.Name))
	}
	data, er := jwt.Decrypt(token, auth.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, context.NewError(types.GetInt(auth.FailedCode, 403), er)
		}
		return data, context.NewError(types.GetInt(auth.FailedCode, 403), er)
	}
	return data, nil
}
func getToken(ctx *gin.Context, jwt *conf.Auth) string {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		return ctx.GetHeader(jwt.Name)
	default:
		cookie, _ := ctx.Cookie(jwt.Name)
		return cookie
	}
}

func setToken(ctx *gin.Context, jwt *conf.Auth, token string) {
	switch strings.ToUpper(jwt.Source) {
	case "HEADER", "H":
		ctx.Header(jwt.Name, token)
	default:
		if jwt.Domain != "" {
			ctx.Header("Set-Cookie", fmt.Sprintf("%s=%s;domain=%s;path=/;", jwt.Name, token, jwt.Domain))
			return
		}
		ctx.Header("Set-Cookie", fmt.Sprintf("%s=%s;path=/;", jwt.Name, token))
	}
}
