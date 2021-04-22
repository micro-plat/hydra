package creator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/conf/server"
	varpub "github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/lib4go/types"
)

func GetImportConfs(importPath string) (types.XMap, error) {
	if importPath == "" {
		return nil, nil
	}
	file, err := os.Open(importPath)
	if err != nil {
		return nil, fmt.Errorf("打开导入配置文件错误:%+v", err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取导入配置文件错误:%+v", err)
	}

	confs := make(types.XMap)
	if err := json.Unmarshal(content, &confs); err != nil {
		return nil, fmt.Errorf("导入配置格式转换错误:%+v", err)
	}
	return confs, nil
}

//getImportValue 当导入配置存在相同路径时，则将导入数据转换成结构体或直接返回
func getImportValue(path string, v interface{}, importConf types.XMap) (value interface{}, err error) {
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
