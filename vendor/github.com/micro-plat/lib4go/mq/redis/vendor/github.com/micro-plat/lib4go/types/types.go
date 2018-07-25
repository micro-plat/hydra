package types

import (
	"fmt"
	"strconv"
	"strings"
)

//DecodeString 判断变量的值与指定相等时设置为另一个值，否则使用原值
func DecodeString(def interface{}, a interface{}, b interface{}, e ...interface{}) string {
	values := make([]interface{}, 0, len(e)+2)
	values = append(values, a)
	values = append(values, b)
	values = append(values, e...)

	for i := 0; i < len(values)-1; i = i + 2 {
		if def == values[i] {
			return fmt.Sprint(values[i+1])
		}
	}
	if len(values)%2 == 1 {
		return fmt.Sprint(values[len(values)-1])
	}
	if s, ok := def.(string); ok {
		return s
	}
	return ""
}

//DecodeInt 判断变量的值与指定相等时设置为另一个值，否则使用原值
func DecodeInt(def interface{}, a interface{}, b interface{}, e ...interface{}) int {
	values := make([]interface{}, 0, len(e)+2)
	values = append(values, a)
	values = append(values, b)
	values = append(values, e...)

	for i := 0; i < len(values)-1; i = i + 2 {
		if def == values[i] {
			v, err := Convert2Int(values[i+1])
			if err == nil {
				return v
			}
		}
	}
	if len(values)%2 == 1 {
		v, err := Convert2Int(values[len(values)-1])
		if err == nil {
			return v
		}
	}
	if r, ok := def.(int); ok {
		return r
	}
	return 0
}

//Convert2Int 转换为int类型
func Convert2Int(i interface{}) (int, error) {
	switch i.(type) {
	case int:
		return i.(int), nil
	case string:
		return strconv.Atoi(i.(string))
	default:
		return strconv.Atoi(fmt.Sprint(i))
	}
}

func ToInt(i interface{}, def ...int) int {
	v, err := Convert2Int(i)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return v
}

//IsEmpty 当前对像是否是字符串空
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	if t, ok := v.(string); ok && len(t) == 0 {
		return true
	}
	if t, ok := v.([]interface{}); ok && len(t) == 0 {
		return true
	}
	return false
}
func GetString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i.(type) {
	case []string:
		return strings.Join(i.([]string), ";")
	default:
		return fmt.Sprint(i)
	}
}
func IntContains(input []int, v int) bool {
	for _, i := range input {
		if i == v {
			return true
		}
	}
	return false
}
