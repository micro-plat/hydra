package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/net"
)

//APIKeyAuth 静态密钥验证
func APIKeyAuth() Handler {

	return func(ctx IMiddleContext) {

		//获取apikey配置
		auth, err := ctx.APPConf().GetAPIKeyConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}

		if auth.Disable {
			ctx.Next()
			return
		}
		if ok, _ := auth.Match(ctx.Request().Path().GetRequestPath()); ok {
			ctx.Next()
			return
		}

		//检查必须参数
		ctx.Response().AddSpecial("apikey")
		if err := ctx.Request().Check("sign", "timestamp"); err != nil {
			ctx.Response().Abort(http.StatusUnauthorized, err)
			return
		}

		//验证签名
		sign, raw := getSignRaw(ctx.Request(), "", "")
		if err := auth.Verify(raw, sign, ctx.Invoke); err != nil {
			ctx.Response().Abort(http.StatusForbidden, err)
			return
		}
		ctx.Next()
	}
}

//getSecret 获取密钥
func getSecret(ctx context.IContext, auth *apikey.APIKeyAuth) (string, error) {
	var secret = auth.Secret
	proto, addr, err := global.ParseProto(secret)
	if err != nil {
		return secret, nil
	}
	if proto == global.ProtoRPC {
		response, err := components.Def.RPC().GetRegularRPC().Swap(addr, ctx)
		if err != nil {
			return "", err
		}
		secret = response.Result
	}
	return "", fmt.Errorf("apikey不支持协议%s", proto)
}

//getSignRaw 构建原串
func getSignRaw(r context.IRequest, a string, b string) (string, string) {
	keys := r.Keys()
	values := net.NewValues()
	var sign string
	for _, key := range keys {
		switch key {
		case "sign":
			sign = r.GetString(key)
		default:
			values.Set(key, r.GetString(key))
		}
	}
	values.Sort()
	return sign, values.Join(a, b)
}
