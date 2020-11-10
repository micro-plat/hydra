package mocks

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

var _ conf.IVarConf = &MockVarConf{}

type MockVarConf struct {
	PlatName string
	Version  int32
	ConfData map[string]map[string]*conf.RawConf
}

func (v *MockVarConf) GetVersion() int32 {
	return v.Version
}
func (v *MockVarConf) GetConf(tp string, name string) (*conf.RawConf, error) {
	data, ok := v.ConfData[tp][name]
	if !ok {
		return nil, conf.ErrNoSetting
	}
	return data, nil
}
func (v *MockVarConf) GetConfVersion(tp string, name string) (int32, error) {
	data, ok := v.ConfData[tp][name]
	if !ok {
		return 0, conf.ErrNoSetting
	}
	return data.GetVersion(), nil
}
func (v *MockVarConf) GetObject(tp string, name string, res interface{}) (int32, error) {
	data, ok := v.ConfData[tp][name]
	if !ok {
		return 0, conf.ErrNoSetting
	}
	return data.GetVersion(), data.ToStruct(res)
}
func (v *MockVarConf) GetClone() conf.IVarConf {
	return conf.IVarConf(v)
}
func (v *MockVarConf) Has(tp string, name string) bool {
	_, ok := v.ConfData[tp][name]
	return ok
}
func (v *MockVarConf) Iter(f func(k string, conf *conf.RawConf) bool) {
	data := v.ConfData
	for k, v := range data {
		for sk, sv := range v {
			f(registry.Join(k, sk), sv)
		}
	}

}

//GetRLogPath 获取远程日志配置路径
func (c *MockVarConf) GetRLogPath() string {
	return c.GetVarPath("app", "rlog")
}

//GetVarPath 获取var配置路径
func (c *MockVarConf) GetVarPath(tp ...string) string {
	if len(tp) == 0 {
		return registry.Join(c.PlatName, "var")
	}
	l := make([]string, 0, len(tp)+2)
	l = append(l, c.PlatName)
	l = append(l, "var")
	l = append(l, tp...)
	return registry.Join(l...)
}
