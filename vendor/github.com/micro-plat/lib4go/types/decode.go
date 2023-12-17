package types

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

//DecodeBool 判断变量的值与指定相等时设置为另一个值，否则使用原值
func DecodeBool(input interface{}, a interface{}, b interface{}, e ...interface{}) bool {
	values := make([]interface{}, 0, len(e)+2)
	values = append(values, a)
	values = append(values, b)
	values = append(values, e...)

	def, _ := ParseBool(input)
	for i := 0; i < len(values)-1; i = i + 2 {
		if def == values[i] {
			v, b := MustBool(values[i+1])
			if b {
				return v
			}
		}
	}
	if len(values)%2 == 1 {
		v, b := MustBool(values[len(values)-1])
		if b {
			return v
		}
	}
	return def
}

//DeepCopy 深拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
