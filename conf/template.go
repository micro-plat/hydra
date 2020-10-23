package conf

import (
	"bytes"
	"text/template"
)

//TmpltTranslate 翻译模板
func TmpltTranslate(name string, tmplt string, funcs map[string]interface{}, input interface{}) (c string, err error) {
	tmpl, err := getTmplt(name, tmplt, funcs)
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
func getTmplt(name string, tmpl string, funcs map[string]interface{}) (*template.Template, error) {
	return template.New(name).Funcs(funcs).Parse(tmpl)
}
