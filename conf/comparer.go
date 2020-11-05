package conf

import (
	"reflect"
)

//ICompare 配置比较器
type IComparer interface {
	Update(n IServerConf)
	IsChanged() bool
	IsValueChanged(names ...string) (isChanged bool)
	IsSubConfChanged(names ...string) (isChanged bool)
}

//Comparer 配置比较器
type Comparer struct {
	oconf      IServerConf
	nconf      IServerConf
	valueNames []string
	subNames   []string
}

//NewComparer 构建配置比较器
//初始配置不能为空，新配置为空时认为没有变更
func NewComparer(oconf IServerConf, valueNames []string, subNames ...string) *Comparer {
	if oconf == nil {
		panic("配置不能为空")
	}
	return &Comparer{
		oconf:      oconf,
		valueNames: valueNames,
		subNames:   subNames,
	}
}

//Update 更新配置
func (s *Comparer) Update(n IServerConf) {
	if s.nconf != nil {
		s.oconf = s.nconf
	}
	s.nconf = n
}

//IsChanged 检查版本号是否发生变化
//当新配置nconf为空时，系统认为未发生变化
func (s *Comparer) IsChanged() bool {
	if s.nconf == nil || reflect.ValueOf(s.nconf).IsNil() {
		return false
	}
	return s.oconf.GetVersion() < s.nconf.GetVersion()
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

//IsSubConfChanged 子配置项目是否发生变化
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
