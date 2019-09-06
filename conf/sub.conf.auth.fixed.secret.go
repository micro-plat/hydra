package conf

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type FixedSecretAuth struct {
	Mode    string `json:"mode" valid:"in(MD5|SHA1|SHA256),required"`
	Secret  string `json:"secret" valid:"ascii,required"`
	Disable bool   `json:"disable,omitempty"`
}

//WithFixedSecretSign 添加固定签名认证
func (a Authes) WithFixedSecretSign(auth *FixedSecretAuth) Authes {
	a["fixed-secret"] = auth
	return a
}

//NewFixedSecretAuth 创建固定Secret签名认证
func NewFixedSecretAuth(secret string, mode ...string) *FixedSecretAuth {
	return &FixedSecretAuth{
		Secret: secret,
		Mode:   strings.ToUpper(types.GetStringByIndex(mode, 0, "MD5")),
	}
}

//WithDisable 禁用配置
func (a *FixedSecretAuth) WithDisable() *FixedSecretAuth {
	a.Disable = true
	return a
}

//WithEnable 启用配置
func (a *FixedSecretAuth) WithEnable() *FixedSecretAuth {
	a.Disable = false
	return a
}
