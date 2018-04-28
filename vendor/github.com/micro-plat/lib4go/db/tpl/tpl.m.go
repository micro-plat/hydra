package tpl

import (
	"fmt"
	"regexp"
	"strings"
)

//MTPLContext  SQLite模板
type MTPLContext struct {
	name string
}

//GetSQLContext 获取查询串
func (o MTPLContext) GetSQLContext(tpl string, input map[string]interface{}) (query string, args []interface{}) {
	f := func() string {
		return "?"
	}
	return AnalyzeTPLFromCache(o.name, tpl, input, f)
}

//GetSPContext 获取存储过程
func (o MTPLContext) GetSPContext(tpl string, input map[string]interface{}) (query string, args []interface{}) {
	return o.GetSQLContext(tpl, input)
}

//Replace 替换SQL中的占位符
func (o MTPLContext) Replace(sql string, args []interface{}) (r string) {
	if strings.EqualFold(sql, "") || args == nil {
		return sql
	}
	word, _ := regexp.Compile(`\?([,|\ ;)]|$)`)
	index := -1
	sql = word.ReplaceAllStringFunc(sql, func(s string) string {
		index++
		if index >= len(args) {
			return "NULL" + s[1:]
		}
		return fmt.Sprintf("'%v'%s", args[index], s[1:])
	})
	return sql
}
func init() {
	// Register("sqlite", MTPLContext{name: "sqlite"})
}
