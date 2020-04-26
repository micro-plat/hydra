package tmpl

import (
	"bytes"
	"text/template"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/security/md5"
)

//GetStatus 获取翻译后的状态码
func (r *RspTmplItem) GetStatus(input interface{}) (c string, err error) {
	return r.translate(r.Status, input)
}

//GetContent 获取翻译后的返回内容
func (r *RspTmplItem) GetContent(input interface{}) (c string, err error) {
	return r.translate(r.Content, input)
}

//Translate 翻译模板
func (r *RspTmplItem) translate(s string, input interface{}) (c string, err error) {
	if s == "" {
		return "", nil
	}
	tmpl, err := getTemplate(s)
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

var tmpCache cmap.ConcurrentMap

func init() {
	tmpCache = cmap.New(2)
}

//getTemplate 获取模板信息
func getTemplate(ts string) (*template.Template, error) {
	key := md5.Encrypt(ts)
	_, tmp, err := tmpCache.SetIfAbsentCb(key, func(input ...interface{}) (c interface{}, err error) {
		t := input[0].(string)
		tmpl, err := template.New(key).Parse(t)
		return tmpl, err

	}, ts)
	if err != nil {
		return nil, err
	}
	return tmp.(*template.Template), nil
}
