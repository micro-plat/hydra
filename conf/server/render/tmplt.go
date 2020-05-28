package render

import (
	"bytes"
	"text/template"

	"github.com/micro-plat/lib4go/security/md5"
)

//translate 翻译模板
func translate(s string, funcs map[string]interface{}, input interface{}) (c string, err error) {
	tmpl, err := getTmplt(s, funcs)
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
	key := md5.Encrypt(tmpl)
	return template.New(key).Funcs(funcs).Parse(tmpl)

}
