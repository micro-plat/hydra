package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//GetString 获取字符串
func GetString(v interface{}, def ...string) string {
	if v != nil {
		if r := fmt.Sprintf("%v", v); r != "" {
			return r
		}
	}
	if len(def) > 0 {
		return def[0]
	}

	return ""
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
	if strings.Contains(strings.ToUpper(value), "E+") {
		var n float64
		_, err := fmt.Sscanf(value, "%e", &n)
		if err == nil {
			return int(n)
		}
		if len(def) > 0 {
			return def[0]
		}
	}
	if value, err := strconv.Atoi(value); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetInt64 获取int64数据，不是有效的数字则返回默然值或0
func GetInt64(v interface{}, def ...int64) int64 {
	value := fmt.Sprintf("%v", v)
	if strings.Contains(strings.ToUpper(value), "E+") {
		var n float64
		_, err := fmt.Sscanf(value, "%e", &n)
		if err == nil {
			return int64(n)
		}
		if len(def) > 0 {
			return def[0]
		}
	}
	if value, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat32 获取float32数据，不是有效的数字则返回默然值或0
func GetFloat32(v interface{}, def ...float32) float32 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 32); err == nil {
		return float32(value)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat64 获取float64数据，不是有效的数字则返回默然值或0
func GetFloat64(v interface{}, def ...float64) float64 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetBool 获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func GetBool(v interface{}, def ...bool) bool {
	if value, err := parseBool(v); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

//GetDatatime 获取时间
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
	if value, err := strconv.Atoi(fmt.Sprintf("%v", v)); err == nil {
		return value, true
	}
	return 0, false
}

//MustFloat32 获取float32，不是有效的数字则返回false
func MustFloat32(v interface{}) (float32, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 32); err == nil {
		return float32(value), true
	}
	return 0, false
}

//MustFloat64 获取float64，不是有效的数字则返回false
func MustFloat64(v interface{}) (float64, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); err == nil {
		return value, true
	}
	return 0, false
}

//IsEmpty 值是否为空
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

//IntContains int数组中是否包含指定值
func IntContains(input []int, v int) bool {
	for _, i := range input {
		if i == v {
			return true
		}
	}
	return false
}
