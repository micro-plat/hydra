package pkgs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

var def = ""

//GetStringByHeader 设置头信息
func GetStringByHeader(content interface{}, hd ...string) string {
	header := types.NewXMap()
	header.Append(hd)

	out := types.NewXMap()
	out.SetValue("__data__", types.StringToBytes(GetString(content)))
	out.SetValue("__header__", header)
	return string(out.Marshal())
}

//GetString 将任意类型转换为字符串，map,struct等转换为json
func GetString(content interface{}) string {
	vtpKind := getTypeKind(content)
	if vtpKind == reflect.String {
		text := fmt.Sprint(content)
		switch {
		case json.Valid([]byte(text)) && (strings.HasPrefix(text, "{") ||
			strings.HasPrefix(text, "[")):
			return text
		default:
			panic("不支持非json字符串")
		}

	} else if vtpKind == reflect.Struct || vtpKind == reflect.Map {
		if buff, err := json.Marshal(content); err != nil {
			panic(err)
		} else {
			return string(buff)
		}
	}
	panic(fmt.Sprintf("不支持的数据类型:%s", vtpKind))
}
func getTypeKind(c interface{}) reflect.Kind {
	if c == nil {
		return reflect.String
	}
	value := reflect.ValueOf(c)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Kind()
}
