package creator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/conf/server"
	varpub "github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/lib4go/types"
)

//getValue 当导入配置存在相同路径时，则将导入数据转换成结构体或直接返回
func getValue(path string, v interface{}, importConf types.XMap) (value interface{}, err error) {
	value = v

	importValue, ok := importConf[path]
	if !ok {
		return
	}
	//服务配置存在对应的导入配置项
	importMap, ok := importConf[path].(map[string]interface{})
	if !ok {
		value = importValue
		return
	}

	//导入配置序列化成对应的结构体
	temp := reflect.ValueOf(v)
	if temp.Kind() == reflect.Ptr {
		temp = temp.Elem()
	}
	if temp.Kind() != reflect.Struct {
		value = importConf[path]
		return
	}
	err = types.NewXMapByMap(importMap).ToStruct(&value)
	return
}

func checkImportConf(tps []string, platName string, systemName string, clusterName string, importConf types.XMap) (err error) {
	prefix := make([]string, 0, len(tps)+1)
	for _, tp := range tps {
		prefix = append(prefix, server.NewServerPub(platName, systemName, tp, clusterName).GetServerPath())
	}
	prefix = append(prefix, varpub.NewVarPub(platName).GetVarPath())
	for path := range importConf {
		allow := false
		for _, v := range prefix {
			if strings.HasPrefix(path, v) {
				allow = true
				break
			}
		}
		if !allow {
			return fmt.Errorf("导入配置的路径%s不正确", path)
		}
	}
	return nil
}
