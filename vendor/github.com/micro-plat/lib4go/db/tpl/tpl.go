package tpl

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

const (
	cOra    = "ora"
	cOracle = "oracle"
	cSqlite = "sqlite"
)

var (
	tpls      map[string]ITPLContext
	tplCaches cmap.ConcurrentMap
)

//ITPLContext 模板上下文
type ITPLContext interface {
	GetSQLContext(tpl string, input map[string]interface{}) (query string, args []interface{})
	GetSPContext(tpl string, input map[string]interface{}) (query string, args []interface{})
	Replace(sql string, args []interface{}) (r string)
}

func init() {
	tpls = make(map[string]ITPLContext)
	tplCaches = cmap.New(8)

	Register("oracle", ATTPLContext{name: "oracle", prefix: ":"})
	Register("ora", ATTPLContext{name: "ora", prefix: ":"})
	Register("mysql", MTPLContext{name: "mysql", prefix: "?"})
	Register("sqlite", MTPLContext{name: "sqlite", prefix: "?"})
	Register("postgres", ATTPLContext{name: "postgres", prefix: "$"})
}
func Register(name string, tpl ITPLContext) {
	if _, ok := tpls[name]; ok {
		panic("重复的注册:" + name)
	}
	tpls[name] = tpl
}

//GetDBContext 获取数据库上下文操作
func GetDBContext(name string) (ITPLContext, error) {
	if v, ok := tpls[strings.ToLower(name)]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("不支持的数据库类型:%s", name)
}
