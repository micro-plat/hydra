package pkgs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/lib4go/envs"
	"github.com/micro-plat/lib4go/types"
)

var def = ""

//兼容hydra0-1版本的mqc服务
var queues01 = map[string]string{}

func init() {
	lst := strings.Split(envs.GetString("hydra01queues"), ",")
	for _, i := range lst {
		queues01[i] = ""
	}
}

//GetStringByHeader 设置头信息
func GetStringByHeader(name string, content interface{}, hd ...string) string {
	//兼容老版本
	if _, ok := queues01[name]; ok {
		return GetString(content)
	}

	header := make(map[string]string, 0)
	if len(hd)%2 == 0 {
		for i := 0; i < len(hd)/2; i++ {
			header[fmt.Sprint(hd[i])] = hd[i+1]
		}
	}

	//新版本hydra
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
