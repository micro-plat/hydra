package types

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

//GetString 获取字符串
func GetString(v interface{}, def ...string) string {
	if v != nil {
		switch v.(type) {
		case float32:
			d := decimal.NewFromFloat32(v.(float32))
			return d.String()
		case float64:
			d := decimal.NewFromFloat(v.(float64))
			return d.String()
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return GetStringByIndex(def, 0)
}

//GetMax 获取指定参数的最大值
func GetMax(v interface{}, o ...int) int {
	r := GetInt(v)
	if len(o) > 0 && o[0] > r {
		return o[0]
	}
	return r
}

//GetMin 获取指定参数的最小值
func GetMin(v interface{}, o ...int) int {
	r := GetInt(v)
	if len(o) > 0 && o[0] < r {
		return o[0]
	}
	return r
}

//GetInt 获取int数据，不是有效的数字则返回默然值或0
func GetInt(v interface{}, def ...int) int {
	value := fmt.Sprintf("%v", v)
	d, err := decimal.NewFromString(value)
	if err != nil {
		return GetIntByIndex(def, 0)
	}

	//如果分母!=1  说明是小数
	if d.Rat().Denom().Int64() != 1 {
		return GetIntByIndex(def, 0)
	}

	res := d.BigInt()
	if res.IsInt64() {
		return int(res.Int64())
	}
	return GetIntByIndex(def, 0)
}

//GetInt32 获取int32数据，不是有效的数字则返回默然值或0
func GetInt32(v interface{}, def ...int32) int32 {
	value := fmt.Sprintf("%v", v)
	d, err := decimal.NewFromString(value)
	if err != nil {
		return GetInt32ByIndex(def, 0)
	}

	//如果分母!=1  说明是小数
	if d.Rat().Denom().Int64() != 1 {
		return GetInt32ByIndex(def, 0)
	}

	res := d.BigInt()
	if res.IsInt64() {
		if res.Int64() > math.MaxInt32 || res.Int64() < math.MinInt32 {
			return GetInt32ByIndex(def, 0)
		}
		return int32(res.Int64())
	}
	return GetInt32ByIndex(def, 0)
}

//GetInt64 获取int64数据，不是有效的数字则返回默然值或0
func GetInt64(v interface{}, def ...int64) int64 {
	value := fmt.Sprintf("%v", v)
	d, err := decimal.NewFromString(value)
	if err != nil {
		return GetInt64ByIndex(def, 0)
	}

	//如果分母!=1  说明是小数
	if d.Rat().Denom().Int64() != 1 {
		return GetInt64ByIndex(def, 0)
	}

	res := d.BigInt()
	if res.IsInt64() {
		return res.Int64()
	}
	return GetInt64ByIndex(def, 0)
}

//GetFloat32 获取float32数据，不是有效的数字则返回默然值或0
func GetFloat32(v interface{}, def ...float32) float32 {
	value := fmt.Sprintf("%v", v)
	d, err := decimal.NewFromString(value)
	if err != nil {
		return GetFloat32ByIndex(def, 0)
	}
	nv, _ := d.BigFloat().Float32()
	if float64(nv) == math.Inf(-1) || float64(nv) == math.Inf(1) {
		return GetFloat32ByIndex(def, 0)
	}
	return nv
}

//GetFloat64 获取float64数据，不是有效的数字则返回默然值或0
func GetFloat64(v interface{}, def ...float64) float64 {
	value := fmt.Sprintf("%v", v)
	d, err := decimal.NewFromString(value)
	if err != nil {
		return GetFloat64ByIndex(def, 0)
	}
	nv, _ := d.BigFloat().Float64()
	if float64(nv) == math.Inf(-1) || float64(nv) == math.Inf(1) {
		return GetFloat64ByIndex(def, 0)
	}
	return nv
}

//GetDecimal 获取float64数据，不是有效的数字则返回默然值或0
func GetDecimal(v interface{}, def ...Decimal) Decimal {
	if value, err := decimal.NewFromString(fmt.Sprintf("%v", v)); err == nil {
		return value
	}
	return GetDecimalByIndex(def, 0)
}

//GetBool 获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func GetBool(v interface{}, def ...bool) bool {
	if value, err := ParseBool(v); err == nil {
		return value
	}
	return GetBoolByIndex(def, 0)
}

//GetDatetime 获取时间
func GetDatetime(v interface{}, format ...string) (time.Time, error) {
	t, b := MustString(v)
	if !b {
		return time.Now(), errors.New("值不能为空")
	}
	f := "2006/01/02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	return time.ParseInLocation(f, t, time.Local)
}

//MustString 获取字符串，不是字符串格式则返回false
func MustString(v interface{}) (string, bool) {
	if value, ok := v.(string); ok {
		return value, true
	}
	return "", false
}

//MustInt 获取int，不是有效的数字则返回false
func MustInt(v interface{}) (int, bool) {
	if value, ok := v.(int); ok {
		return value, true
	}
	return 0, false
}

//MustInt32 获取int32，不是有效的数字则返回false
func MustInt32(v interface{}) (int32, bool) {
	if value, ok := v.(int32); ok {
		return value, true
	}
	return 0, false
}

//MustInt64 获取int64，不是有效的数字则返回false
func MustInt64(v interface{}) (int64, bool) {
	if value, ok := v.(int64); ok {
		return value, true
	}
	return 0, false
}

//MustFloat32 获取float32，不是有效的数字则返回false
func MustFloat32(v interface{}) (float32, bool) {
	if value, ok := v.(float32); ok {
		vn, _ := decimal.NewFromFloat32(value).BigFloat().Float32()
		return vn, true
	}
	return 0, false
}

//MustFloat64 获取float64，不是有效的数字则返回false
func MustFloat64(v interface{}) (float64, bool) {
	if value, ok := v.(float64); ok {
		vn, _ := decimal.NewFromFloat(value).BigFloat().Float64()
		return vn, true
	}
	return 0, false
}

//IsEmpty 值是否为空
func IsEmpty(vs ...interface{}) bool {
	for _, v := range vs {
		if v == nil {
			return true
		}
		tp := reflect.TypeOf(v).Kind()
		value := reflect.ValueOf(v)
		if tp == reflect.Ptr {
			value = value.Elem()
		}
		switch tp {
		case reflect.Chan, reflect.Map, reflect.Slice:
			if value.Len() == 0 {
				return true
			}
		default:
			if value.IsZero() {
				return true
			}
		}
	}
	return false
}

//IntContains int数组中是否包含指定值
func IntContains(input []int, v int) bool {
	for _, i := range input {
		if i == v {
			return true
		}
	}
	return false
}

//StringContains string数组中是否包含指定值
func StringContains(input []string, v string) bool {
	for _, i := range input {
		if i == v {
			return true
		}
	}
	return false
}

//GetStringByIndex 获取数组中的指定元素
func GetStringByIndex(v []string, index int, def ...string) string {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

//GetIntByIndex 获取数组中的指定元素
func GetIntByIndex(v []int, index int, def ...int) int {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetBoolByIndex 获取数组中的指定元素
func GetBoolByIndex(v []bool, index int, def ...bool) bool {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

//GetInt32ByIndex 获取数组中的指定元素
func GetInt32ByIndex(v []int32, index int, def ...int32) int32 {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetInt64ByIndex 获取数组中的指定元素
func GetInt64ByIndex(v []int64, index int, def ...int64) int64 {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat32ByIndex 获取数组中的指定元素
func GetFloat32ByIndex(v []float32, index int, def ...float32) float32 {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat64ByIndex 获取数组中的指定元素
func GetFloat64ByIndex(v []float64, index int, def ...float64) float64 {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetDecimalByIndex 获取数组中的指定元素
func GetDecimalByIndex(v []Decimal, index int, def ...Decimal) Decimal {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return decimal.Zero
}

//ParseBool 将字符串转换为bool值
func ParseBool(val interface{}) (value bool, err error) {
	if val == nil {
		return false, fmt.Errorf("parsing <nil>: invalid syntax")
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch strings.ToUpper(v) {
		case "1", "T", "TRUE", "YES", "Y", "ON":
			return true, nil
		case "0", "F", "FALSE", "NO", "N", "OFF":
			return false, nil
		}
	case int, int8, int16, int32, int64, float32, float64:
		if v == 0 {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("parsing %q: invalid syntax", val)
}

//Translate 翻译带参数的变量支持格式有 @abc,{@abc}
func Translate(format string, kv ...interface{}) string {
	trf := NewXMap()
	trf.Append(kv...)
	return trf.Translate(format)
}
