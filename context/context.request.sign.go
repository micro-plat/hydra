package context

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/security/sha1"
	"github.com/micro-plat/lib4go/security/sha256"
	"github.com/micro-plat/lib4go/types"
)

//GetSignRaw 检查签名原串
func (r *Request) GetSignRaw(all bool, a string, b string) (string, string) {
	keys := r.GetKeys()
	values := net.NewValues()
	var sign string
	for _, key := range keys {
		switch key {
		case "sign", "signature":
			sign = r.GetString(key)
		default:
			values.Set(key, r.GetString(key))
		}
	}
	values.Sort()
	if all {
		return sign, values.JoinAll(a, b)
	}
	return sign, values.Join(a, b)

}

//CheckSign 检查签名是否正确(只值为空不参与签名，键值及每个串之前使用空进行连接)
func (r *Request) CheckSign(secret string, tp ...string) (bool, error) {
	return r.CheckSignAll(secret, false, "", "", tp...)
}

//CheckSignAll 检查签名是否正确
func (r *Request) CheckSignAll(secret string, all bool, a string, b string, tp ...string) (bool, error) {
	sign, raw := r.GetSignRaw(all, a, b)
	var expect string
	switch strings.ToUpper(types.GetStringByIndex(tp, 0, "md5")) {
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
	return false, fmt.Errorf("raw:%s,expect:%s,actual:%s", raw, expect, sign)
}
