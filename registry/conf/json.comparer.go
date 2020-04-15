package config

import (
	"reflect"
)

//ConfComparer 配置比较器
type ConfComparer struct {
	oconf IMainConf
	nconf IMainConf
}

//NewConfComparer 构建配置比较器
func NewConfComparer(oconf IMainConf, nconf IMainConf) *ConfComparer {
	return &ConfComparer{
		oconf: oconf,
		nconf: nconf,
	}
}

//IsChanged 主配置是否发生变更
func (s *ConfComparer) IsChanged() bool {
	return s.oconf == nil || reflect.ValueOf(s.oconf).IsNil() || s.oconf.GetVersion() != s.nconf.GetVersion()
}

//IsValueChanged 配置内容是否发生变化
func (s *ConfComparer) IsValueChanged(names ...string) (isChanged bool) {
	if reflect.ValueOf(s.oconf).IsNil() {
		return true
	}
	for _, name := range names {
		if s.oconf == nil || reflect.ValueOf(s.oconf).IsNil() {
			return true
		}
		if s.nconf.GetMainConf().GetString(name) != s.oconf.GetMainConf().GetString(name) {
			return true
		}
	}
	return false
}

//IsSubConfChanged 子配置是否发生变化
func (s *ConfComparer) IsSubConfChanged(name string) (isChanged bool) {
	if s.oconf == nil || reflect.ValueOf(s.oconf).IsNil() {
		return true
	}
	o, err := s.oconf.GetSubConf(name)
	if err != nil {
		return true
	}

	n, err := s.nconf.GetSubConf(name)
	if err != nil {
		return true
	}

	if o.version != n.version {
		return true
	}
	return false
}
