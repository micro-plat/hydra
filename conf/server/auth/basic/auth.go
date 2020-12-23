package basic

import "github.com/micro-plat/lib4go/encoding/base64"

type auth struct {
	userName string
	password string
	auth     string
}

//初始化认证对象列表
func newAuthorization(m map[string]string) []*auth {
	pairs := make([]*auth, 0, len(m))
	for user, password := range m {
		value := createAuth(user, password)
		pairs = append(pairs, &auth{
			userName: user,
			password: password,
			auth:     value,
		})
	}
	return pairs
}

//创建认证关键值
func createAuth(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.Encode(base)
}
