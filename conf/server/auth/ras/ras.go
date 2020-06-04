package ras

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//RASAuth 远程服务验证组
type RASAuth struct {
	Disable bool    `json:"disable,omitempty" toml:"disable,omitempty"`
	Auth    []*Auth `json:"auth"`
}

//NewRASAuth 构建RASAuth认证
func NewRASAuth(auth ...*Auth) *RASAuth {
	r := &RASAuth{}
	for _, a := range auth {
		a.PathMatch = conf.NewPathMatch(a.Requests...)
		r.Auth = append(r.Auth, a)
	}
	return r
}

//Match 检查指定的路径是否有对应的认证服务
func (a RASAuth) Match(p string) (bool, *Auth) {
	for _, auth := range a.Auth {
		if ok, _ := auth.Match(p); ok {
			return true, auth
		}
	}
	return false, nil
}

//GetConf 获取配置信息
func GetConf(cnf conf.IMainConf) (auths *RASAuth) {
	auths = &RASAuth{}
	//设置Remote安全认证参数
	_, err := cnf.GetSubObject(registry.Join("auth", "RASAuth"), auths)
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("RASAuth配置有误:%v", err))
	}
	if err == conf.ErrNoSetting {
		auths.Disable = true
		return auths
	}

	for _, auth := range auths.Auth {
		if b, err := govalidator.ValidateStruct(&auth); !b {
			panic(fmt.Errorf("RASAuth配置有误:%v", err))
		}
		auth.PathMatch = conf.NewPathMatch(auth.Requests...)
	}
	return auths
}
