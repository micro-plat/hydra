package apm

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)


type APM struct {
	Disable  bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	Config  string    `json:"config,omitempty" toml:"config,omitempty"`
	Cache []string    `json:"cache,omitempty" toml:"cache,omitempty"`
	DB    []string    `json:"db,omitempty" toml:"db,omitempty"`
}

func (s *APM) GetEnable() bool {
	return !s.Disable
}

func (s *APM) GetCache(name string) bool { 
	return s.getVal( s.Cache,name)
}

func (s *APM) GetDB(name string) bool {
	return s.getVal(s.DB, name)
}

func (s *APM) GetConfig() string {
	return  s.Config	
}
 
func (s *APM) getVal(list []string ,name  string) ( ok bool) {
	if len(list) == 0 {
		return false
	}
	for i :=range list{
		if list[i]==name{
			return true
		}
	}
	return false	
}
 

//New 创建固定密钥验证服务
func New(opts ...Option) *APM {
	f := &APM{Disable:true}
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
	fsa := New(WithDisable(), WithConfig("apm"))
	_, err := cnf.GetSubObject("apm", fsa)
	if err == conf.ErrNoSetting {
		return fsa
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(fsa); !b {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	return fsa
}
