package tpl

import (
	"fmt"
	"regexp"
	"strings"
)

func isNil(input interface{}) bool {
	return input == nil || fmt.Sprintf("%v", input) == ""
}

//AnalyzeTPL 解析模板内容，并返回解析后的SQL语句，入输入参数
//@表达式，替换为参数化字符如: :1,:2,:3
//#表达式，替换为指定值，值为空时返回NULL
//~表达式，检查值，值为空时返加"",否则返回: , name=value
//&条件表达式，检查值，值为空时返加"",否则返回: and name=value
//|条件表达式，检查值，值为空时返回"", 否则返回: or name=value
func AnalyzeTPL(tpl string, input map[string]interface{}, prefix func() string, like func(string, func() string) string) (sql string, params []interface{}, names []string) {
	params = make([]interface{}, 0)
	names = make([]string, 0)
	defer func() {
		sql = replaceSpecialCharacter(sql)
	}()

	//@变量, 将数据放入params中
	word, _ := regexp.Compile(`[\\]?([@#&~!$?]|[|]{1,2})(\w+(\.\w+)?)`)
	sql = word.ReplaceAllStringFunc(tpl, func(s string) string {
		index := 1
		if strings.HasPrefix(s, "||") {
			index = 2
		}
		pre, key, name := s[:index], s[index:], s[index:]
		if strings.Index(key, ".") > 0 {
			name = strings.Split(key, ".")[1]
		}
		if strings.Contains(pre, "\\") {
			pre = "\\"
		}
		value := input[name]
		switch pre {
		case "@":
			if !isNil(value) {
				names = append(names, key)
				params = append(params, value)
			} else {
				names = append(names, key)
				params = append(params, "")
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
				return fmt.Sprintf("and %s like %s", key, like(key, prefix))
			}
			return ""
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
			return pre + key
		}
	})

	//@变量, 将数据放入params中
	word2, _ := regexp.Compile(`[\\][@#&~\|!\$\?><]`)
	sql = word2.ReplaceAllStringFunc(sql, func(s string) string {
		return s[1:]
	})
	return
}
