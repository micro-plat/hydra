package gray

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

//Check 检查当前是否需要转到上游服务器处理
func (g *Gray) Check(funcs map[string]interface{}, i interface{}) (bool, error) {
	r, err := translate(g.Filter, funcs, i)
	if err != nil {
		return true, fmt.Errorf("%s 过滤器转换出错 %w", g.Filter, err)
	}
	return strings.EqualFold(r, "true"), nil
}

//translate 翻译模板
func translate(tmplt string, funcs map[string]interface{}, input interface{}) (c string, err error) {
	tmpl, err := getTmplt(tmplt, funcs)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, input)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

//getTmplt 获取模板信息
func getTmplt(tmpl string, funcs map[string]interface{}) (*template.Template, error) {
	return template.New(tmpl).Funcs(funcs).Parse(tmpl)

}
