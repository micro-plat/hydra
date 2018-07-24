package creator

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/micro-plat/hydra/component"
)

var _ IMainBinder = &MainBinder{}

type IMainBinder interface {
	SetMainConf(s string)
	SetSubConf(n string, s string)
	GetSubConfNames() []string
	Scan(mainConf string, nodeName string) error
	NeedScanCount(nodeName string) int
	GetNodeConf(nodeName string) string
	Installer(func(c component.IContainer) error)
}

//MainBinder 主配置绑定
type MainBinder struct {
	mainConf           string              //系统主配置
	subConf            map[string]string   //子系统配置
	mainParamsForInput []string            //主配置参数，用于用户输入
	subParamsForInput  map[string][]string //子系统参数,用于用户输入
	//	mainConfParamsForTranslate map[string]string            //主配置参数，用于参数翻译
	//	subConfParamsForTranslate  map[string]map[string]string //子系统参数,用于参数翻译
	params     map[string]string
	rmainConf  string            //翻译后的主配置
	rsubConf   map[string]string //翻译后的子系统配置
	installers []func(c component.IContainer) error
}

//NewMainBinder 构建主配置绑定
func NewMainBinder(params map[string]string) *MainBinder {
	return &MainBinder{
		subConf:            make(map[string]string),
		mainParamsForInput: make([]string, 0, 2),
		subParamsForInput:  make(map[string][]string),
		params:             params,
		//mainConfParamsForTranslate: make(map[string]string),
		//	subConfParamsForTranslate:  make(map[string]map[string]string),
		rsubConf:   make(map[string]string),
		installers: make([]func(c component.IContainer) error, 0, 2),
	}
}
func (c *MainBinder) GetInstallers() []func(c component.IContainer) error {
	return c.installers
}
func (c *MainBinder) Installer(f func(c component.IContainer) error) {
	c.installers = append(c.installers, f)
}

//SetMainConf 设置主配置内容
func (c *MainBinder) SetMainConf(s string) {
	c.mainConf = s
	c.mainParamsForInput = getParams(s)
}

//SetSubConf 设置子配置内容
func (c *MainBinder) SetSubConf(n string, s string) {
	c.subConf[n] = s
	params := getParams(s)
	if len(params) > 0 {
		c.subParamsForInput[n] = params
	}
}

//GetSubConfNames 获取子系统名称
func (c *MainBinder) GetSubConfNames() []string {
	v := make([]string, 0, len(c.subConf))
	for k := range c.subConf {
		v = append(v, k)
	}
	return v
}

//NeedScanCount 待输入个数
func (c *MainBinder) NeedScanCount(nodeName string) int {
	count := 0
	if nodeName == "" {
		for _, p := range c.mainParamsForInput {
			if _, ok := c.params[p]; !ok {
				count++
			}
		}
	}
	for _, p := range c.subParamsForInput[nodeName] {
		if _, ok := c.params[p]; !ok {
			count++
		}
	}
	return count
}

//Scan 绑定参数
func (c *MainBinder) Scan(mainConf string, nodeName string) error {
	if nodeName == "" {
		for _, p := range c.mainParamsForInput {
			if _, ok := c.params[p]; ok {
				continue
			}
			fmt.Printf("请输入:%s中%s的值:", mainConf, p)
			var value string
			fmt.Scan(&value)
			//	c.mainConfParamsForTranslate[p] = value
			c.params[p] = value
		}
		c.rmainConf = translate(c.mainConf, c.params)
	} else {
		//c.subConfParamsForTranslate[nodeName] = make(map[string]string)
		for _, p := range c.subParamsForInput[nodeName] {
			if _, ok := c.params[p]; ok {
				continue
			}
			fmt.Printf("请输入:%s中%s的值:", filepath.Join(mainConf, nodeName), p)
			var value string
			fmt.Scan(&value)
			//c.subConfParamsForTranslate[nodeName][p] = value
			c.params[p] = value
		}
		if v, ok := c.subConf[nodeName]; ok {
			c.rsubConf[nodeName] = translate(v, c.params)
		}
	}

	return nil
}

//GetNodeConf 获取节点配置
func (c *MainBinder) GetNodeConf(nodeName string) string {
	if nodeName == "" {
		return c.rmainConf
	}
	if v, ok := c.rsubConf[nodeName]; ok {
		return v
	}
	return ""
}

//getParams 翻译带有@变量的字符串
func getParams(format string) []string {
	brackets, _ := regexp.Compile(`\{#\w+\}`)
	p1 := brackets.FindAllString(format, -1)
	brackets, _ = regexp.Compile(`#\w+`)
	p2 := brackets.FindAllString(format, -1)
	r := make([]string, 0, len(p1)+len(p2))
	for _, v := range p1 {
		r = append(r, v[2:len(v)-1])
	}
	for _, v := range p2 {
		r = append(r, v[1:])
	}
	return r
}

//translate 翻译带有@变量的字符串
func translate(format string, data map[string]string) string {
	brackets, _ := regexp.Compile(`\{#\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		key := s[2 : len(s)-1]
		if v, ok := data[key]; ok {
			return v
		}
		return s
	})
	word, _ := regexp.Compile(`#\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		key := s[1:]
		if v, ok := data[key]; ok {
			return v
		}
		return s
	})
	return result
}
