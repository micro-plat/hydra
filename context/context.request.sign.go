package context

import (
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
func (r *Request) CheckSign(key string, f ...string) (bool, string) {
	return r.CheckSignAll(key, false, "", "", f...)
}

//CheckSignAll 检查签名是否正确
func (r *Request) CheckSignAll(key string, all bool, a string, b string, f ...string) (bool, string) {
	sign, raw := r.GetSignRaw(all, a, b, f...)
	return strings.EqualFold(md5.Encrypt(raw+key), sign), raw
}
