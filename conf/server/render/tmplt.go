package render

import (
	"bytes"
	"text/template"
)

//translate 翻译模板
func translate(path string, tmplt string, funcs map[string]interface{}, input interface{}) (c string, err error) {
	tmpl, err := getTmplt(path, tmplt, funcs)
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
func getTmplt(path string, tmpl string, funcs map[string]interface{}) (*template.Template, error) {
	return template.New(path).Funcs(funcs).Parse(tmpl)

}
