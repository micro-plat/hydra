package types

import "fmt"

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
			v, b := MustInt(values[i+1])
			if b {
				return v
			}
		}
	}
	if len(values)%2 == 1 {
		v, b := MustInt(values[len(values)-1])
		if b {
			return v
		}
	}
	if r, ok := def.(int); ok {
		return r
	}
	return 0
}
