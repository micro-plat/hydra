package creator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	vc "github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/conf/server/metric"
)

type ISUB interface {
	Sub(name string, s ...interface{}) ISUB
}

type BaseBuilder map[string]interface{}

//Sub 子配置
func (b BaseBuilder) Sub(name string, s ...interface{}) ISUB {
	if len(s) == 0 {
		panic(fmt.Sprintf("配置：%s值不能为空", name))
	}
	tp := reflect.TypeOf(s[0])
	val := reflect.ValueOf(s[0])
	switch tp.Kind() {
	case reflect.String:
		str := val.String()
		if strings.HasPrefix(str, vc.ByInstall) {
			b[name] = str
			return b
		}
		b[name] = json.RawMessage([]byte(str))
	case reflect.Struct, reflect.Ptr, reflect.Map:
		b[name] = val.Interface()
	default:
		panic(fmt.Sprintf("配置：%s值类型不支持", name))
	}
	return b
}

//Metric 监控配置
func (b BaseBuilder) Metric(host string, db string, cron string, opts ...metric.Option) BaseBuilder {
	b[metric.TypeNodeName] = metric.New(host, db, cron, opts...)
	return b
}

//APM 构建APM配置
func (b BaseBuilder) APM(address string) BaseBuilder {
	b[apm.TypeNodeName] = apm.New(address)
	return b
}

//Map 将监控配置返回为map
func (b BaseBuilder) Map() map[string]interface{} {
	return b
}

//Load 加载配置内容
func (b BaseBuilder) Load() {
}
