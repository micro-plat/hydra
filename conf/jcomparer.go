package conf

import "reflect"

type JSONComparer struct {
	Oconf IConf
	Nconf IConf
}

func NewJSONComparer(Oconf IConf, Nconf IConf) *JSONComparer {
	return &JSONComparer{
		Oconf: Oconf,
		Nconf: Nconf,
	}
}
func (s *JSONComparer) IsChanged() bool {
	return s.Oconf == nil || reflect.ValueOf(s.Oconf).IsNil() || s.Oconf.GetVersion() != s.Nconf.GetVersion()
}

//IsValueChanged 检查值是否发生变化
func (s *JSONComparer) IsValueChanged(names ...string) (isChanged bool) {
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
