package creator

import (
	"fmt"
	"path/filepath"
)

var _ IPlatBinder = &PlatBinder{}

type IPlatBinder interface {
	SetVarConf(t string, s string, v string)
	Scan(platName string, nodeName string) error
	GetVarNames() []string
	NeedScanCount(nodeName string) int
	GetNodeConf(nodeName string) string
}

//PlatBinder 平台配置绑定
type PlatBinder struct {
	varConf           map[string]string   //var环境参数配置
	varParamsForInput map[string][]string //环境参数，用于用户输入
	//varConfParamsForTranslate map[string]map[string]string //环境参数，用于参数翻译
	params   map[string]string
	rvarConf map[string]string //翻译后的环境参数配置
}

//NewPlatBinder 平台绑定
func NewPlatBinder(params map[string]string) *PlatBinder {
	return &PlatBinder{
		varConf:           make(map[string]string),
		varParamsForInput: make(map[string][]string),
		//	varConfParamsForTranslate: make(map[string]map[string]string),
		params:   params,
		rvarConf: make(map[string]string),
	}
}

//SetVarConf 设置var配置内容
func (c *PlatBinder) SetVarConf(t string, s string, v string) {
	c.varConf[filepath.Join(t, s)] = v
	params := getParams(v)
	if len(params) > 0 {
		c.varParamsForInput[filepath.Join(t, s)] = params
	}
}
func (c *PlatBinder) GetVarNames() []string {
	v := make([]string, 0, len(c.varConf))
	for k := range c.varConf {
		v = append(v, k)
	}
	return v
}

//NeedScanCount 待输入个数
func (c *PlatBinder) NeedScanCount(nodeName string) int {
	count := 0
	for _, p := range c.varParamsForInput[nodeName] {
		if _, ok := c.params[p]; !ok {
			count++
		}
	}
	return count
}

//Scan 绑定参数
func (c *PlatBinder) Scan(platName string, nodeName string) error {
	//c.varConfParamsForTranslate[nodeName] = make(map[string]string)
	for _, p := range c.varParamsForInput[nodeName] {
		if _, ok := c.params[p]; ok {
			continue
		}

		fmt.Printf("请输入:%s中%s的值:", filepath.Join("/", platName, "var", nodeName), p)
		var value string
		fmt.Scan(&value)
		//	c.varConfParamsForTranslate[nodeName][p] = value
		c.params[p] = value
	}
	if v, ok := c.varConf[nodeName]; ok {
		c.rvarConf[nodeName] = translate(v, c.params)
	}
	return nil
}

//GetNodeConf 获取节点配置
func (c *PlatBinder) GetNodeConf(nodeName string) string {
	return c.rvarConf[nodeName]
}
