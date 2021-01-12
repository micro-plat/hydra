package ras

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

const (
	//ParNodeName auth-ras配置父节点名
	ParNodeName = "auth"
	//SubNodeName auth-ras配置子节点名
	SubNodeName = "ras"
)

//RASAuth 远程服务验证组
type RASAuth struct {
	Disable bool    `json:"disable,omitempty" toml:"disable,omitempty"`
	Auth    []*Auth `json:"auth" toml:"auth"`
}

//NewRASAuth 构建RASAuth认证
func NewRASAuth(opts ...Option) *RASAuth {
	r := &RASAuth{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//Match 检查指定的路径是否有对应的认证服务
func (a RASAuth) Match(p string) (bool, *Auth) {
	for _, auth := range a.Auth {
		if ok, _ := auth.Match(p); ok && !auth.Disable {
			return true, auth
		}
	}
	return false, nil
}

//GetConf 获取配置信息
func GetConf(cnf conf.IServerConf) (auths *RASAuth, err error) {
	auths = &RASAuth{}
	//设置Remote安全认证参数
	_, err = cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), auths)
	if errors.Is(err, conf.ErrNoSetting) {
		auths.Disable = true
		return auths, nil
	}
	if err != nil {
		return nil, fmt.Errorf("RASAuth配置格式有误:%v", err)
	}

	for _, auth := range auths.Auth {
		if b, err := govalidator.ValidateStruct(auth); !b {
			return nil, fmt.Errorf("RASAuth配置数据有误:%v", err)
		}
		auth.PathMatch = conf.NewPathMatch(auth.Requests...)
	}
	return auths, nil
}
