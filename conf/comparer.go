package conf

import (
	"reflect"
)

//ICompare 配置比较器
type IComparer interface {
	Update(n IMainConf)
	IsChanged() bool
	IsValueChanged(names ...string) (isChanged bool)
	IsSubConfChanged(names ...string) (isChanged bool)
}

//Comparer 配置比较器
type Comparer struct {
	oconf      IMainConf
	nconf      IMainConf
	valueNames []string
	subNames   []string
}

//NewComparer 构建配置比较器
func NewComparer(oconf IMainConf, valueNames []string, subNames ...string) *Comparer {
	return &Comparer{
		oconf:      oconf,
		valueNames: valueNames,
		subNames:   subNames,
	}
}

//Update 更新配置
func (s *Comparer) Update(n IMainConf) {
	if s.nconf != nil {
		s.oconf = s.nconf
	}
	s.nconf = n
}

//IsChanged 主配置是否发生变更
func (s *Comparer) IsChanged() bool {
	if s.nconf == nil || reflect.ValueOf(s.nconf).IsNil() {
		return false
	}
	return s.oconf.GetVersion() != s.nconf.GetVersion()
}

//IsValueChanged 配置内容是否发生变化
func (s *Comparer) IsValueChanged(names ...string) (isChanged bool) {
	if !s.IsChanged() {
		return false
	}
	knames := append(s.valueNames, names...)
	for _, name := range knames {
		if s.nconf.GetMainConf().GetString(name) != s.oconf.GetMainConf().GetString(name) {
			return true
		}
	}
	return false
}

//IsSubConfChanged 子配置是否发生变化
func (s *Comparer) IsSubConfChanged(names ...string) (isChanged bool) {
	if !s.IsChanged() {
		return false
	}
	knames := append(s.subNames, names...)
	for _, name := range knames {
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
	}

	return false
}
