package types

import (
	"fmt"
	"strconv"
	"time"
)

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func GetInt(name interface{}, def ...int) int {
	if value, err := strconv.Atoi(fmt.Sprintf("%v", name)); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetInt64 从对象中获取数据值，如果不是字符串则返回0
func GetInt64(name interface{}, def ...int64) int64 {
	if value, err := strconv.ParseInt(fmt.Sprintf("%v", name), 10, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func GetFloat32(name interface{}, def ...float32) float32 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", name), 32); err == nil {
		return float32(value)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func GetFloat64(name interface{}, def ...float64) float64 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", name), 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetBool 从对象中获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func GetBool(name interface{}, def ...bool) bool {
	if value, err := parseBool(name); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

//GetDatatime 获取时间字段
func GetDatatime(name interface{}, format ...string) (time.Time, error) {
	t, b := GetMustString(name)
	if !b {
		return time.Now(), fmt.Errorf("%s列不存在", name)
	}
	f := "2006/01/02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	return time.ParseInLocation(f, t, time.Local)
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func GetMustString(name interface{}) (string, bool) {
	if value, ok := name.(string); ok {
		return value, true
	}
	return "", false
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func GetMustInt(name interface{}) (int, bool) {
	if value, err := strconv.Atoi(fmt.Sprintf("%v", name)); err == nil {
		return value, true
	}
	return 0, false
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func GetMustFloat32(name interface{}) (float32, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", name), 32); err == nil {
		return float32(value), true
	}
	return 0, false
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func GetMustFloat64(name interface{}) (float64, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", name), 64); err == nil {
		return value, true
	}
	return 0, false
}

//ParseBool 将字符串转换为bool值
func parseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%s", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1 {
				return true, nil
			} else if v == 0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}
