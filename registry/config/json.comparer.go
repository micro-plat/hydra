package config

import (
	"reflect"
)

//ConfComparer 配置比较器
type ConfComparer struct {
	Oconf IConf
	Nconf IConf
}


//NewConfComparer 构建配置比较器
func NewConfComparer(Oconf IMainConf, Nconf IMainConf) *ConfComparer {
	return &ConfComparer{
		Oconf: Oconf,
		Nconf: Nconf,
	}
}

//IsChanged 主配置是否发生变更
func (s *ConfComparer) IsChanged() bool {
	return s.Oconf == nil || reflect.ValueOf(s.Oconf).IsNil() || s.Oconf.GetVersion() != s.Nconf.GetVersion()
}

//IsValueChanged 配置内容是否发生变化
func (s *ConfComparer) IsValueChanged(names ...string) (isChanged bool) {
	if reflect.ValueOf(s.Oconf).IsNil() {
		return true
	}
	for _, name := range names {
		if s.Nconf.GetString(name) != s.Oconf.GetString(name) {
			return true
		}
	}
	return false
}

//IsSubConfChanged 子配置是否发生变化
func (s *Comparer) IsSubConfChanged(name string) (isChanged bool) {
	oldConf, _ := s.Oconf.GetSubConf(name)
	newConf, _ := s.Nconf.GetSubConf(name)
	if oldConf == nil {
		oldConf = &JSONConf{version: 0}
	}
	if newConf == nil {
		newConf = &JSONConf{version: 0}
	}
	return oldConf.version != newConf.version
}