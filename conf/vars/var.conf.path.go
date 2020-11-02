package vars

import "github.com/micro-plat/hydra/registry"

type varPub struct {
	platName string
}

//NewVarPub 构建var路径管理器
func NewVarPub(platName string) *varPub {
	return &varPub{
		platName: platName,
	}
}

//GetVarPath 获取var配置路径
func (c *varPub) GetVarPath(tp ...string) string {
	if len(tp) == 0 {
		return registry.Join(c.platName, "var")
	}
	l := make([]string, 0, len(tp)+2)
	l = append(l, c.platName)
	l = append(l, "var")
	l = append(l, tp...)
	return registry.Join(l...)
}

//GetPlatName 获取平台名称
func (c *varPub) GetPlatName() string {
	return c.platName
}

//GetRLogPath 获取远程日志配置路径
func (c *varPub) GetRLogPath() string {
	return c.GetVarPath("app", "rlog")
}
