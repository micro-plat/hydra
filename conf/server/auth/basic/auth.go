package basic

import "github.com/micro-plat/lib4go/encoding/base64"

type member struct {
	UserName string `json:"userName,omitempty" toml:"userName,omitempty" valid:"required"`
	Password string `json:"password,omitempty" toml:"password,omitempty" valid:"required"`
}

type auth struct {
	userName string
	password string
	auth     string
}

//初始化认证对象列表
func newAuthorization(members []*member) []*auth {
	pairs := make([]*auth, 0, len(members))
	for _, m := range members {
		value := createAuth(m.UserName, m.Password)
		pairs = append(pairs, &auth{
			userName: m.UserName,
			password: m.Password,
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
