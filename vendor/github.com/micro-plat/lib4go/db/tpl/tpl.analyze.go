package tpl

import (
	"fmt"
	"regexp"
	"strings"
)

func isNil(input interface{}) bool {
	return input == nil || fmt.Sprintf("%v", input) == ""
}

type tplCache struct {
	sql    string
	params []interface{}
	names  []string
}

//AnalyzeTPLFromCache 从缓存中获取已解析的SQL语句
func AnalyzeTPLFromCache(name string, tpl string, input map[string]interface{}, prefix func() string) (sql string, params []interface{}) {
	sql, params, _ = AnalyzeTPL(tpl, input, prefix)
	return
	/*key := fmt.Sprintf("%s_%s", name, tpl)
	b, cache, _ := tplCaches.SetIfAbsentCb(key, func(i ...interface{}) (interface{}, error) {
		sql, params, names := AnalyzeTPL(tpl, input, prefix)
		return &tplCache{sql: sql, params: params, names: names}, nil
	})
	value := cache.(*tplCache)
	if b {
		return value.sql, value.params
	}
	params = make([]interface{}, 0, len(value.names))
	for _, v := range value.names {
		va := input[v]
		if !isNil(va) {
			params = append(params, va)
		} else {
			params = append(params, nil)
		}
	}
	return value.sql, params*/
}

//AnalyzeTPL 解析模板内容，并返回解析后的SQL语句，入输入参数
//@表达式，替换为参数化字符如: :1,:2,:3
//#表达式，替换为指定值，值为空时返回NULL
//~表达式，检查值，值为空时返加"",否则返回: , name=value
//&条件表达式，检查值，值为空时返加"",否则返回: and name=value
//|条件表达式，检查值，值为空时返回"", 否则返回: or name=value
func AnalyzeTPL(tpl string, input map[string]interface{}, prefix func() string) (sql string, params []interface{}, names []string) {
	params = make([]interface{}, 0)
	names = make([]string, 0)
	defer func() {
		sql = strings.Replace(strings.Replace(strings.Replace(sql, "  ", " ", -1), "where and", "where", -1), "where or", "where", -1)
		sql = strings.Replace(strings.Replace(sql, "WHERE and", "WHERE", -1), "WHERE or", "WHERE", -1)
	}()
	word, _ := regexp.Compile(`[\\]?[@|#|&|~|\||!|\$|\?]\w?[\.]?\w+`)
	//@变量, 将数据放入params中
	sql = word.ReplaceAllStringFunc(tpl, func(s string) string {
		fullKey := s[1:]
		key := s[1:]
		name := s[1:]
		if strings.Index(fullKey, ".") > 0 {
			name = strings.Split(fullKey, ".")[1]
		}
		pre := s[:1]
		value := input[name]
		switch pre {
		case "@":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
			} else {
				names = append(names, key)
				params = append(params, nil)
			}
			return prefix()
		case "#":
			if !isNil(value) {
				return fmt.Sprintf("%v", value)
			}
			return "NULL"
		case "?":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
				return fmt.Sprintf("and %s like '%%'||%s||'%%'", key, prefix())
			}
			return ""
		// case ",":
		// 	if !isNil(value) {
		// 		names = append(names, key)
		// 		params = append(params, value)
		// 		return fmt.Sprintf("and %s > %s", key, prefix())
		// 	}
		// 	return ""
		// case ".":
		// 	if !isNil(value) {
		// 		names = append(names, key)
		// 		params = append(params, value)
		// 		return fmt.Sprintf("and %s < %s", key, prefix())
		// 	}
		// 	return ""
		case "$":
			if !isNil(value) {
				return fmt.Sprintf("%v", value)
			}
			return ""
		case "&":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
				return fmt.Sprintf("and %s=%s", key, prefix())
			}
			return ""
		case "|":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
				return fmt.Sprintf("or %s=%s", key, prefix())
			}
			return ""
		case "~":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
				return fmt.Sprintf(",%s=%s", key, prefix())
			}
			return ""
		default:
			return s
		}
	})

	word2, _ := regexp.Compile(`[\\][@|#|&|~|\||!|\$|\?|>|<]`)
	//@变量, 将数据放入params中
	sql = word2.ReplaceAllStringFunc(sql, func(s string) string {
		return s[1:]
	})
	return
}
