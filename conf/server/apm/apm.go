package apm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

type m = map[string]interface{}

type APM struct {
	m
}

func (s *APM) GetEnable() bool {
	return s.GetStatus("enable")
}

func (s *APM) GetCache(subname ...string) bool {
	return s.GetStatus("cache", subname...)
}

func (s *APM) GetDB(subname ...string) bool {
	return s.GetStatus("db", subname...)
}

func (s *APM) GetConfig() string {
	v, ok := s.getVal("config")
	if !ok {
		return ""
	}
	r, ok := v.(string)
	if !ok {
		return ""
	}
	return r
}

func (s *APM) GetStatus(ttype string, subname ...string) bool {
	key := ttype
	if len(subname) > 0 && (strings.TrimSpace(subname[0]) != "") {
		key = fmt.Sprintf("%s.%s", ttype, strings.TrimSpace(subname[0]))
	}
	v, ok := s.getVal(key)
	if !ok {
		return false
	}
	r, ok := v.(bool)
	if !ok {
		return false
	}
	return r
}
func (s *APM) String() string {
	bytes, _ := json.Marshal(s.m)
	return string(bytes)
}

func (s *APM) getVal(key string) (v interface{}, ok bool) {
	if len(s.m) <= 0 {
		return
	}
	v, ok = s.m[key]
	return
}

// MarshalJSON implements the json.Marshaller interface.
func (s *APM) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.m)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (s *APM) UnmarshalJSON(data []byte) error {
	mapData := map[string]interface{}{}
	err := json.Unmarshal(data, &mapData)
	s.m = mapData
	return err
}

// //APM 调用链配置
// type APM struct {
// 	Enable bool   `json:"enable,omitempty" toml:"enable,omitempty"`
// 	Config string `json:"config,omitempty" toml:"config,omitempty"`
// 	DB     bool   `json:"db,omitempty" toml:"db,omitempty"`
// 	Cache  bool   `json:"cache,omitempty" toml:"cache,omitempty"`
// }

//New 创建固定密钥验证服务
func New(opts ...Option) *APM {
	f := &APM{
		m: map[string]interface{}{},
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
