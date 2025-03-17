package tpl

import (
	"fmt"
	"regexp"
	"strings"
)

func isNil(input interface{}) bool {
	return input == nil || fmt.Sprintf("%v", input) == ""
}

// AnalyzeTPL 解析模板内容，并返回解析后的SQL语句，入输入参数
// @表达式，替换为参数化字符如: :1,:2,:3
// #表达式，替换为指定值，值为空时返回NULL
// ~表达式，检查值，值为空时返加"",否则返回: , name=value
// &条件表达式，检查值，值为空时返加"",否则返回: and name=value
// |条件表达式，检查值，值为空时返回"", 否则返回: or name=value
func AnalyzeTPL(tpl string, input map[string]interface{}, prefix func() string, like func(string, func() string) string) (sql string, params []interface{}, names []string) {
	params = make([]interface{}, 0)
	names = make([]string, 0)
	defer func() {
		sql = replaceSpecialCharacter(sql)
	}()

	// 匹配模板中的变量表达式（包括转义字符）
	word, _ := regexp.Compile(`(\\?)([@#&~!$?]|\|{1,2})(\w+(\.\w+)?)`)
	sql = word.ReplaceAllStringFunc(tpl, func(s string) string {
		// 判断是否有转义字符
		groups := word.FindStringSubmatch(s)
		escapeChar := groups[1] // 转义字符（\）
		pre := groups[2]        // 前缀（@、#、&、~、| 等）
		key := groups[3]        // 变量名

		// 如果有转义字符，直接返回原始字符串（去掉转义字符）
		if escapeChar == "\\" {
			return pre + key
		}

		// 处理变量名（去掉可能的点号）
		name := key
		if strings.Contains(key, ".") {
			name = strings.Split(key, ".")[1]
		}

		// 获取输入值
		value := input[name]

		// 根据前缀处理不同的逻辑
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

	// 处理单独的转义字符（如 \@、\# 等）
	word2, _ := regexp.Compile(`\\([@#&~|!$?><_])`)
	sql = word2.ReplaceAllStringFunc(sql, func(s string) string {
		return s[1:] // 去掉转义字符
	})

	return
}
