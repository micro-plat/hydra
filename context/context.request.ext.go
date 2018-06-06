package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

type extParams struct {
	ext map[string]interface{}
	ctx *Context
}

// func (w *extParams) Get(name string) (interface{}, bool) {
// 	v, ok := w.ext[name]
// 	return v, ok
// }
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

func (w *extParams) GetBodyMap(encoding ...string) map[string]interface{} {

	if fun, ok := w.ext["__get_request_values_"].(func() map[string]interface{}); ok {
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

//GetJWTBody 获取jwt存储数据
func (w *extParams) GetJWT(out interface{}) error {
	jwt := w.ext["__jwt_"]
	if jwt == nil || reflect.ValueOf(jwt).IsNil() {
		return fmt.Errorf("未找到jwt,用户未登录")
	}
	switch v := jwt.(type) {
	case string:
		return json.Unmarshal([]byte(v), &out)
	case map[string]interface{}:
		return types.Map2Struct(v, out)
	default:
		buff, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(buff, out)
	}
}

//GetUUID
func (w *extParams) GetUUID() string {
	return w.ext["__hydra_sid_"].(string)
}

//GetJWTConfig 获取jwt配置信息
func (w *extParams) GetJWTConfig() (*conf.Auth, error) {
	var auths conf.Authes
	var jwt *conf.Auth
	if _, err := w.ctx.GetContainer().GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("jwt配置有误:%v", err)
		return nil, err
	}
	jwt, enable := auths["jwt"]
	if !enable {
		return nil, fmt.Errorf("jwt:%v", conf.ErrNoSetting)
	}
	return jwt, nil
}
