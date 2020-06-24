package server

import "github.com/micro-plat/hydra/registry"

type varPath struct {
	platName string
}

//NewVarPath 构建var路径管理器
func NewVarPath(platName string) *varPath {
	return &varPath{
		platName: platName,
	}
}

//GetVarPath 获取var配置路径
func (c *varPath) GetVarPath(tp ...string) string {
	if len(tp) == 0 {
		return registry.Join(c.platName, "var")
	}
	l := make([]string, 0, len(tp)+2)
	l = append(l, c.platName)
	l = append(l, "var")
	l = append(l, tp...)
	return registry.Join(l...)
}

//GetRLogPath 获取远程日志配置路径
func (c *varPath) GetRLogPath() string {
	return c.GetVarPath("app", "rlog")
}
