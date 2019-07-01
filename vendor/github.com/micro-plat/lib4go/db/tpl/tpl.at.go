package tpl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//ATTPLContext 参数化时使用@+参数名作为占位符的SQL数据库如:oracle,sql server
type ATTPLContext struct {
	name   string
	prefix string
}

func (o ATTPLContext) getSPName(query string) string {
	return fmt.Sprintf("begin %s;end;", strings.Trim(strings.Trim(query, ";"), ","))
}

//GetSQLContext 获取查询串
func (o ATTPLContext) GetSQLContext(tpl string, input map[string]interface{}) (sql string, args []interface{}) {
	index := 0
	f := func() string {
		index++
		return fmt.Sprint(o.prefix, index)
	}
	return AnalyzeTPLFromCache(o.name, tpl, input, f)
}

//GetSPContext 获取
func (o ATTPLContext) GetSPContext(tpl string, input map[string]interface{}) (sql string, args []interface{}) {
	q, args := o.GetSQLContext(tpl, input)
	sql = o.getSPName(q)
	return
}

//Replace 替换SQL中的占位符
func (o ATTPLContext) Replace(sql string, args []interface{}) (r string) {
	if sql == "" || args == nil {
		return sql
	}
	word, _ := regexp.Compile(fmt.Sprintf(`%s\d+([,|\) ;]|$)`, o.prefix))
	sql = word.ReplaceAllStringFunc(sql, func(s string) string {
		c := len(s)
		num := s[1 : c-1]
		// 处理匹配到结尾
		if num == "" {
			num = s[1:c]
			c++
		}
		k, err := strconv.Atoi(num)
		if err != nil || len(args) < k {
			return "NULL" + s[c-1:]
		}
		return fmt.Sprintf("'%v'%s", args[k-1], s[c-1:])
	})
	/*end*/
	return sql
}
