package conf

type Comparer struct {
	Oconf IServerConf
	Nconf IServerConf
	*JSONComparer
}

func NewComparer(Oconf IServerConf, Nconf IServerConf) *Comparer {
	return &Comparer{
		Oconf:        Oconf,
		Nconf:        Nconf,
		JSONComparer: NewJSONComparer(Oconf, Nconf),
	}
}

//IsVarChanged var节点是否发生变化
func (s *Comparer) IsVarChanged() bool {
	return s.Oconf.GetVarVersion() != s.Nconf.GetVarVersion()
}

//IsSubConfChanged 检查节点是否发生变化
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

//IsRequiredSubConfChanged 检查必须节点是否发生变化
func (s *Comparer) IsRequiredSubConfChanged(name string) (isChanged bool, err error) {
	oldConf, _ := s.Oconf.GetSubConf(name)
	newConf, _ := s.Nconf.GetSubConf(name)
	if oldConf == nil {
		oldConf = &JSONConf{version: 0}
	}
	if newConf == nil {
		newConf = &JSONConf{version: 0}
	}
	return oldConf.GetVersion() != newConf.GetVersion(), nil
}
