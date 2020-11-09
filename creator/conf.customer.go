package creator

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type ISUB interface {
	Sub(name string, s ...interface{}) ISUB
}

type iCustomerBuilder interface {
	Load()
	ISUB
	Map() map[string]interface{}
}

var _ iCustomerBuilder = CustomerBuilder{}

type CustomerBuilder map[string]interface{}

//newHTTP 构建http生成器
func newCustomerBuilder(s ...interface{}) CustomerBuilder {
	b := make(map[string]interface{})
	if len(s) == 0 {
		b[ServerMainNodeName] = make(map[string]interface{})
		return b
	}
	b[ServerMainNodeName] = s[0]
	return b
}

func (b CustomerBuilder) Sub(name string, s ...interface{}) ISUB {
	if len(s) == 0 {
		panic(fmt.Sprintf("配置：%s值不能为空", name))
	}
	tp := reflect.TypeOf(s[0])
	val := reflect.ValueOf(s[0])
	if tp.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch tp.Kind() {
	case reflect.String:
		b[name] = json.RawMessage([]byte(val.Interface().(string)))
	case reflect.Struct, reflect.Ptr, reflect.Map:
		b[name] = val.Interface()
	default:
		panic(fmt.Sprintf("配置：%s值类型不支持", name))
	}
	return b
}
func (b CustomerBuilder) Map() map[string]interface{} {
	return b
}
func (b CustomerBuilder) Load() {
}
