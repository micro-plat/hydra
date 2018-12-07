package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/security/jwt"
	"github.com/micro-plat/lib4go/types"
)

//JwtAuth jwt
func wsCheckJwt(ctx *gin.Context, curl string, token string) bool {
	cnf := getMetadataConf(ctx)
	jwtAuth, ok := cnf.GetMetadata("jwt").(*conf.Auth)
	if !ok || jwtAuth == nil || jwtAuth.Disable {
		return true
	}

	//检查jwt.token是否正确
	data, err := checkWSJWT(ctx, jwtAuth, token)
	if err == nil {
		setJWTRaw(ctx, data)
		return true
	}

	//不需要校验的URL自动跳过
	for _, u := range jwtAuth.Exclude {
		if u == curl {
			return true
		}
	}
	return false
}
func makeJwtToken(ctx *gin.Context, data interface{}) (string, bool) {
	if data == nil {
		data = getJWTRaw(ctx)
	}
	if data == nil {
		return "", false
	}
	cnf := getMetadataConf(ctx)
	jwtAuth, ok := cnf.GetMetadata("jwt").(*conf.Auth)
	if !ok || jwtAuth.Disable {
		return "", false
	}
	jwtToken, err := jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
	if err != nil {
		getLogger(ctx).Errorf("jwt配置出错：%v", err)
		return "", false
	}
	return jwtToken, true
}

// CheckJWT 检查jwk参数是否合法
func checkWSJWT(ctx *gin.Context, auth *conf.Auth, token string) (data interface{}, err context.IError) {
	if token == "" {
		return nil, context.NewError(types.GetInt(auth.FailedCode, 403), fmt.Errorf("获取%s失败或未传入该参数", auth.Name))
	}
	data, er := jwt.Decrypt(token, auth.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, context.NewError(types.GetInt(auth.FailedCode, 401), er)
		}
		return data, context.NewError(types.GetInt(auth.FailedCode, 403), er)
	}
	return data, nil
}
func getJWTToken(ctx *gin.Context) string {
	jwtAuth, ok := getMetadataConf(ctx).GetMetadata("jwt").(*conf.Auth)
	if !ok || jwtAuth == nil || jwtAuth.Disable {
		return ""
	}
	return getToken(ctx, jwtAuth)
	// cnf := getMetadataConf(ctx)
	// if cookie, err := ctx.Cookie(cnf.Name); err == nil {
	// 	return cookie
	// }
	// return ""
}
