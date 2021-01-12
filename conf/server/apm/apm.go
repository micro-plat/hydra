package apm

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName APM配置节点名
const TypeNodeName = "apm"

type IAPM interface {
	GetConf() (*APM, bool)
}

//APM APM
type APM struct {
	Address string `json:"address,omitempty" valid:"required" toml:"address,omitempty" label:"应用程序性能监控地址"`
	Version int32  `json:"-"`
	Disable bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//New 构建api server配置信息
func New(address string, opts ...Option) *APM {
	m := &APM{
		Address: address,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

//GetConf 设置APM
func GetConf(cnf conf.IServerConf) (apm *APM, err error) {
	apm = &APM{}
	_, err = cnf.GetSubObject(TypeNodeName, apm)
	if errors.Is(err, conf.ErrNoSetting) {
		apm.Disable = true
		return apm, nil
	}
	if err != nil {
		return nil, err
	}

	if c, err := cnf.GetSubConf(TypeNodeName); err == nil {
		apm.Version = c.GetVersion()
	}
	if b, err := govalidator.ValidateStruct(apm); !b {
		return nil, fmt.Errorf("apm配置数据有误:%v", err)
	}
	return
}
