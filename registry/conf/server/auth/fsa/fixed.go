package fsa

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/registry/conf"
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

//GetConf 获取FixedSecretAuth
func GetConf(cnf conf.IMainConf) (fsa *FixedSecretAuth) {
	_, err := cnf.GetSubObject("fixed-secret", &fsa)
	if err == conf.ErrNoSetting {
		return &FixedSecretAuth{fixedOption: &fixedOption{Disable: true}}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("fixed-secret配置有误:%v", err))
	}
	if fsa != nil {
		if b, err := govalidator.ValidateStruct(&fsa); !b {
			panic(fmt.Errorf("fixed-secret配置有误:%v", err))
		}
	}
	return fsa
}
