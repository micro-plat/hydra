package creator

import (
	"github.com/micro-plat/hydra/registry"
)

// var _ IplatBinder = &platBinder{}

type iplatBinder interface {
	SetVarConf(t string, s string, v interface{})
	scan(platName string, nodeName string) error
	getVarNames() []string
	needScanCount(nodeName string) int
	getNodeConf(nodeName string) string
}

//platBinder 平台配置绑定
type platBinder struct {
	varConf           map[string]string   //var环境参数配置
	varParamsForInput map[string][]string //环境参数，用于用户输入
	inputs            map[string]*Input
	params            map[string]string
	rvarConf          map[string]string //翻译后的环境参数配置
}

//newPlatBinder 平台绑定
func newPlatBinder(params map[string]string, inputs map[string]*Input) *platBinder {
	return &platBinder{
		varConf:           make(map[string]string),
		varParamsForInput: make(map[string][]string),
		inputs:            inputs,
		params:            params,
		rvarConf:          make(map[string]string),
	}
}

//SetVarConf 设置var配置内容
func (c *platBinder) SetVarConf(t string, s string, input interface{}) {
	v, err := getConfig(input)
	if err != nil {
		panic(err)
	}
	c.varConf[registry.Join(t, s)] = v
	params := getParams(v)
	if len(params) > 0 {
		c.varParamsForInput[registry.Join(t, s)] = params
	}
}
func (c *platBinder) getVarNames() []string {
	v := make([]string, 0, len(c.varConf))
	for k := range c.varConf {
		v = append(v, k)
	}
	return v
}

//NeedScanCount 待输入个数
func (c *platBinder) needScanCount(nodeName string) int {
	count := 0
	for _, p := range c.varParamsForInput[nodeName] {
		if _, ok := c.params[p]; !ok {
			count++
		}
	}
	return count
}

//Scan 绑定参数
func (c *platBinder) scan(platName string, nodeName string) error {
	for _, p := range c.varParamsForInput[nodeName] {
		if _, ok := c.params[p]; ok {
			continue
		}

		nvalue, err := getInputValue(p, c.inputs, registry.Join("/", platName, "var", nodeName))
		if err != nil {
			return err
		}
		c.params[p] = nvalue

	}
	if v, ok := c.varConf[nodeName]; ok {
		c.rvarConf[nodeName] = Translate(v, c.params)
	}
	return nil
}

//GetNodeConf 获取节点配置
func (c *platBinder) getNodeConf(nodeName string) string {
	return c.rvarConf[nodeName]
}
