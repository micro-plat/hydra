package context

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
)

//GetSignRaw 检查签名原串
func (r *Request) GetSignRaw(all bool, a string, b string, f ...string) (string, string) {
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
func (r *Request) CheckSign(secret string, f ...string) (bool, error) {
	return r.CheckSignAll(secret, false, "", "", f...)
}

//CheckSignAll 检查签名是否正确
func (r *Request) CheckSignAll(secret string, all bool, a string, b string, f ...string) (bool, error) {
	sign, raw := r.GetSignRaw(all, a, b, f...)
	expect := md5.Encrypt(raw + secret)
	if strings.EqualFold(expect, sign) {
		return true, nil
	}
	return false, fmt.Errorf("raw:%s,expect:%s,actual:%s", raw, expect, sign)
}
