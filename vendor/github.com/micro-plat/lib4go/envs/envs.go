package envs

import (
	"os"
	"time"

	"github.com/micro-plat/lib4go/types"
)

//GetString 从环境变量中, 获取字符串
func GetString(name string, def ...string) string {
	value := os.Getenv(name)
	return types.GetString(value, def...)
}

//GetMax 从环境变量中, 获取指定参数的最大值
func GetMax(name string, o ...int) int {
	value := GetString(name)
	return types.GetMax(value, o...)
}

//GetMin 从环境变量中, 获取指定参数的最小值
func GetMin(name string, o ...int) int {
	value := GetString(name)
	return types.GetMin(value, o...)
}

//GetInt 从环境变量中, 获取int数据，不是有效的数字则返回默然值或0
func GetInt(name string, def ...int) int {
	value := GetString(name)
	return types.GetInt(value, def...)
}

//GetInt64 从环境变量中, 获取int64数据，不是有效的数字则返回默然值或0
func GetInt64(name string, def ...int64) int64 {
	value := GetString(name)
	return types.GetInt64(value, def...)
}

//GetFloat32 从环境变量中, 获取float32数据，不是有效的数字则返回默然值或0
func GetFloat32(name string, def ...float32) float32 {
	value := GetString(name)
	return types.GetFloat32(value, def...)
}

//GetFloat64 从环境变量中, 获取float64数据，不是有效的数字则返回默然值或0
func GetFloat64(name string, def ...float64) float64 {
	value := GetString(name)
	return types.GetFloat64(value, def...)
}

//GetBool 从环境变量中, 获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func GetBool(name string, def ...bool) bool {
	value := GetString(name)
	return types.GetBool(value, def...)
}

//GetDatatime 从环境变量中, 获取时间
func GetDatetime(name string, format ...string) (time.Time, error) {
	value := GetString(name)
	return types.GetDatetime(value, format...)
}

//MustInt 从环境变量中, 获取int，不是有效的数字则返回false
func MustInt(name string) (int, bool) {
	value := GetString(name)
	return types.MustInt(value)
}

//MustFloat32 从环境变量中, 获取float32，不是有效的数字则返回false
func MustFloat32(name string) (float32, bool) {
	value := GetString(name)
	return types.MustFloat32(value)
}

//MustFloat64 从环境变量中, 获取float64，不是有效的数字则返回false
func MustFloat64(name string) (float64, bool) {
	value := GetString(name)
	return types.MustFloat64(value)
}
