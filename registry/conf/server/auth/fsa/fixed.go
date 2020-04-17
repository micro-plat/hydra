package fsa

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//FixedSecretAuth 创建固定密钥验证服务
type FixedSecretAuth struct {
	*fixedOption
}

//NewFixedSecret 创建固定密钥验证服务
func NewFixedSecret(secret string, opts ...FixedOption) *FixedSecretAuth {
	f := &FixedSecretAuth{fixedOption: &fixedOption{
		Secret:  secret,
		Include: []string{"*"},
		Mode:    "MD5",
	}}
	for _, opt := range opts {
		opt(f.fixedOption)
	}
	return f
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

//GetFixedSecret 获取FixedSecret
func GetFixedSecret(cnf conf.IMainConf) (auths *conf.Authes, err error) {
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("fixed-secret配置有误:%v", err)
	}
	if auths.FixedScret != nil {
		if b, err := govalidator.ValidateStruct(auths.FixedScret); !b {
			return nil, fmt.Errorf("fixed-secret配置有误:%v", err)
		}
	}
	return auths, nil
}
