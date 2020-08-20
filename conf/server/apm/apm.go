package apm

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//APM 调用链配置
type APM struct {
	Disable bool `json:"disable,omitempty" toml:"disable,omitempty"`
}

//New 创建固定密钥验证服务
func New(opts ...Option) *APM {
	f := &APM{
		Disable: false,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type ConfHandler func(cnf conf.IMainConf) *APM

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 获取APIKeyAuth
func GetConf(cnf conf.IMainConf) *APM {
	fsa := APM{}
	_, err := cnf.GetSubObject("apm", &fsa)
	if err == conf.ErrNoSetting {
		return &APM{Disable: true}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&fsa); !b {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	return &fsa
}
