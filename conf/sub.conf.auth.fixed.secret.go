package conf

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type FixedSecretAuth struct {
	Mode    string   `json:"mode" valid:"in(MD5|SHA1|SHA256),required"`
	Secret  string   `json:"secret" valid:"ascii,required"`
	Include []string `json:"include" valid:"required"`
	Disable bool     `json:"disable,omitempty"`
}

//WithFixedSecretSign 添加固定签名认证
func (a *Authes) WithFixedSecretSign(auth *FixedSecretAuth) *Authes {
	a.FixedScret = auth
	return a
}

//NewFixedSecretAuth 创建固定Secret签名认证
func NewFixedSecretAuth(secret string, mode ...string) *FixedSecretAuth {
	return &FixedSecretAuth{
		Secret:  secret,
		Include: []string{"*"},
		Mode:    strings.ToUpper(types.GetStringByIndex(mode, 0, "MD5")),
	}
}

//WithInclude 设置include的请求服务路径
func (a *FixedSecretAuth) WithInclude(path ...string) *FixedSecretAuth {
	if len(path) > 0 {
		a.Include = path
	}
	return a

}

//Contains 检查指定的路径是否允许签名
func (a *FixedSecretAuth) Contains(p string) bool {
	if len(a.Include) == 0 {
		return true
	}
	for _, i := range a.Include {
		if i == "*" || i == p {
			return true
		}
	}
	return false
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
