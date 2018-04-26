package context

import (
	"errors"
	"fmt"
	"strings"
)

type extParams struct {
	ext map[string]interface{}
}

func (w *extParams) Get(name string) (interface{}, bool) {
	v, ok := w.ext[name]
	return v, ok
}
func (w *extParams) GetMethod() string {
	if m, ok := w.ext["__method_"].(string); ok {
		return m
	}
	return ""
}
func (w *extParams) GetHeader() map[string]string {
	if h, ok := w.ext["__header_"].(map[string]string); ok {
		return h
	}
	header := make(map[string]string)
	if h, ok := w.ext["__header_"].(map[string][]string); ok {
		for k, v := range h {
			header[k] = strings.Join(v, ",")
		}
	}
	return header
}

func (w *extParams) GetBindingFunc() func(i interface{}) error {
	v, ok := w.ext["__binding_"].(func(i interface{}) error)
	if ok {
		return v
	}
	panic(errors.New("不存在__binding_函数"))
}
func (w *extParams) GetBindWithFunc() func(i interface{}, c string) error {
	v, ok := w.ext["__binding_with_"].(func(i interface{}, c string) error)
	if ok {
		return v
	}
	panic(errors.New("不存在__binding_with_函数"))
}

//GetSharding 获取任务分片信息(分片索引[从1开始]，分片总数)
func (w *extParams) GetSharding() (int, int) {
	v, ok := w.ext["__get_sharding_index_"]
	if !ok {
		return 0, 0
	}
	if f, ok := v.(func() (int, int)); ok {
		return f()
	}
	return 0, 0
}

func (w *extParams) GetBodyMap(encoding ...string) map[string]string {
	if fun, ok := w.ext["__get_request_values_"].(func() map[string]string); ok {
		return fun()
	}
	return nil

}
func (w *extParams) GetBody(encoding ...string) (string, error) {
	e := "utf-8"
	if len(encoding) > 0 {
		e = encoding[0]
	}
	if fun, ok := w.ext["__func_body_get_"].(func(ch string) (string, error)); ok {
		return fun(e)
	}
	return "", fmt.Errorf("无法根据%s格式转换数据", e)
}

//GetJWTBody 获取jwt存储数据
func (w *extParams) GetJWTBody() interface{} {
	return w.ext["__jwt_"]
}

//GetUUID
func (w *extParams) GetUUID() string {
	return w.ext["__hydra_sid_"].(string)
}
