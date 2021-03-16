package pkgs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/envs"
	"github.com/micro-plat/lib4go/types"
)

var def = ""

//Queues01 兼容hydra0-1版本的mqc服务
var queues01 = map[string]bool{}

func init() {
	lst := strings.Split(envs.GetString("hydra01queues"), ",")
	for _, i := range lst {
		queues01[i] = true
	}
	key := fmt.Sprintf("%s_%s", global.Def.PlatName, "hydra01queues")
	lst = strings.Split(envs.GetString(key), ",")
	for _, i := range lst {
		queues01[i] = true
	}

}

//GetStringByHeader 设置头信息
func GetStringByHeader(name string, content interface{}, hd ...string) string {
	//兼容老版本
	if IsOriginalQueue(name) {
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

//IsOriginalQueue 是否是老版本队列
func IsOriginalQueue(name string) bool {

	if _, ok := queues01[name]; ok {
		return true
	}
	name = global.MQConf.GetQueueName(name)
	if _, ok := queues01[name]; ok {
		return true
	}

	return false
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
