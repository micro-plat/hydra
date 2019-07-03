package transform

import (
	"fmt"
	"net/url"
	"regexp"
	"sync"
)

type ITransformGetter interface {
	Set(string, string)
	Get(string) (string, error)
	Each(func(string, string))
}
type transformData map[string]string

func (t transformData) Get(key string) (string, error) {
	if v, ok := t[key]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	return "", fmt.Errorf("key(%s) not exist", key)
}
func (t transformData) Set(key string, value string) {
	t[key] = value
}
func (i transformData) Each(f func(string, string)) {
	for k, v := range i {
		f(k, v)
	}
}

//Transform 翻译组件
type Transform struct {
	Data  ITransformGetter
	mutex sync.Mutex
}

//New 创建翻译组件
func New() *Transform {
	var data transformData = make(map[string]string)
	return &Transform{Data: data}
}

//NewValues getter
func NewValues(t url.Values) *Transform {
	var data transformData = make(map[string]string)

	for k, v := range t {
		if len(v) > 1 {
			data[k] = fmt.Sprint(v)
		} else if len(v) > 0 {
			data[k] = fmt.Sprint(v[0])
		}
	}
	return &Transform{Data: data}
}
func NewGetter(t ITransformGetter) *Transform {
	return &Transform{Data: t}
}

//NewMaps 根据map创建组件
func NewMaps(d map[string]interface{}) *Transform {
	var data transformData = make(map[string]string)
	for k, v := range d {
		data[k] = fmt.Sprint(v)
	}
	return &Transform{Data: data}
}

//NewMap create by map
func NewMap(d map[string]string) *Transform {
	var data transformData = make(map[string]string)
	for k, v := range d {
		data[k] = fmt.Sprint(v)
	}
	return &Transform{Data: data}
}

//Append ITransformGetter 添加到当前对象
func (d *Transform) Append(t ITransformGetter) {
	t.Each(func(k, v string) {
		d.Set(k, v)
	})
}

//Set 设置变量的值
func (d *Transform) Set(k string, v string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Data.Set(k, v)
}

//Get 获取变量的值
func (d *Transform) Get(k string) (string, error) {
	return d.Data.Get(k)
}

//Translate 翻译带有@变量的字符串
func (d *Transform) Translate(format string) string {
	return d.TranslateAll(format, false)
}

//TranslateAll 翻译带有@变量的字符串
func (d *Transform) TranslateAll(format string, a bool) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	brackets, _ := regexp.Compile(`\{@\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		if v, err := d.Data.Get(s[2 : len(s)-1]); err == nil {
			return v
		}
		if a {
			return ""
		}
		return s
	})
	word, _ := regexp.Compile(`@\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		if v, err := d.Data.Get(s[1:]); err == nil {
			return v
		}
		if a {
			return ""
		}
		return s
	})
	return result
}

//Translate 翻译字符串kv为map[string]string，map[string]interface{}或以逗号分隔的健值对
func Translate(format string, kv ...interface{}) string {
	if len(kv) == 0 {
		panic(fmt.Sprintf("输入的kv必须为：map[string]string，map[string]interface{}，或健值对"))
	}
	trf := New()
	switch kv[0].(type) {
	case map[string]string:
		trf = NewMap(kv[0].(map[string]string))
	case map[string]interface{}:
		trf = NewMaps(kv[0].(map[string]interface{}))
	default:
		if len(kv)%2 != 0 {
			panic(fmt.Sprintf("输入的kv必须为2的倍数"))
		}
		for i := 0; i < len(kv)-1; i = i + 2 {
			trf.Set(fmt.Sprint(kv[i]), fmt.Sprint(kv[i+1]))
		}
	}
	return trf.Translate(format)
}
