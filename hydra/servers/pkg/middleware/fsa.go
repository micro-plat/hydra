package middleware

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/security/sha1"
	"github.com/micro-plat/lib4go/security/sha256"
)

//FixedSecretAuth 静态密钥验证
func FixedSecretAuth() Handler {

	return func(ctx IMiddleContext) {

		//获取FSA配置
		auth := ctx.ServerConf().GetFSAConf()
		if auth.Disable {
			ctx.Next()
			return
		}
		if ok, _ := auth.In(ctx.Request().Path().GetPath()); ok {
			ctx.Next()
			return
		}

		//检查必须参数
		ctx.Response().AddSpecial("fsa--x")
		if err := ctx.Request().Check("sign", "timestamp"); err != nil {
			ctx.Response().AbortWithError(402, err)
			return
		}

		//验证签名
		_, err := checkSign(ctx.Request(), auth.Secret, auth.Mode)
		if err == nil {
			ctx.Next()
			return
		}
		ctx.Response().AbortWithError(401, err)
	}
}

//checkSign 检查签名是否正确
func checkSign(r context.IRequest, secret string, tp string) (bool, error) {
	sign, raw := getSignRaw(r, "", "")
	var expect string
	switch strings.ToUpper(tp) {
	case "MD5":
		expect = md5.Encrypt(raw + secret)
	case "SHA1":
		expect = sha1.Encrypt(raw + secret)
	case "SHA256":
		expect = sha256.Encrypt(raw + secret)
	default:
		return false, fmt.Errorf("不支持的签名验证方式:%v", tp)
	}
	if strings.EqualFold(expect, sign) {
		return true, nil
	}
	if global.IsDebug {
		return false, fmt.Errorf("签名错误:raw:%s,expect:%s,actual:%s", raw, expect, sign)
	}
	return false, fmt.Errorf("签名错误")

}

//getSignRaw 检查签名原串
func getSignRaw(r context.IRequest, a string, b string) (string, string) {
	keys := r.GetKeys()
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
