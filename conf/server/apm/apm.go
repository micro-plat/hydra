package apm

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

type APM struct {
	Name    string   `json:"name,required" toml:"name,required"`
	Disable bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	Cache   []string `json:"cache,omitempty" toml:"cache,omitempty"`
	DB      []string `json:"db,omitempty" toml:"db,omitempty"`
}

func (s *APM) GetEnable() bool {
	return !s.Disable
}

func (s *APM) GetCacheEnable(name string) bool {
	return s.getVal(s.Cache, name)
}

func (s *APM) GetDBEnable(name string) bool {
	return s.getVal(s.DB, name)
}

func (s *APM) GetName() string {
	return s.Name
}

func (s *APM) getVal(list []string, name string) (ok bool) {
	if len(list) == 0 {
		return false
	}
	for i := range list {
		if list[i] == name {
			return true
		}
	}
	return false
}

//New 创建固定密钥验证服务
func New(name string, opts ...Option) *APM {
	f := &APM{
		Name:    name,
		Disable: true,
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

//GetConf 获取APM
func GetConf(cnf conf.IMainConf) *APM {
	fsa := &APM{
		Disable: true,
	}
	_, err := cnf.GetSubObject("apm", fsa)
	if err == conf.ErrNoSetting {
		return fsa
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("APM配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(fsa); !b {
		panic(fmt.Errorf("APM配置有误:%v", err))
	}
	return fsa
}
