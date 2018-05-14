package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/micro-plat/lib4go/types"
)

type IQueryRow interface {
	GetString(name string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string) float32
	GetFloat64(name string) float64
	Has(name string) bool
	GetMustString(name string) (string, error)
	GetMustInt(name string) (int, error)
	GetMustFloat32(name string) (float32, error)
	GetMustFloat64(name string) (float64, error)
	GetDatatime(name string, format ...string) (time.Time, error)
	ToStruct(o interface{}) error
}

type QueryRow map[string]interface{}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q QueryRow) GetString(name string) string {
	return fmt.Sprintf("%v", q[name])
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetInt(name string, def ...int) int {
	if value, err := strconv.Atoi(fmt.Sprintf("%v", q[name])); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetInt64 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetInt64(name string, def ...int64) int64 {
	if value, err := strconv.ParseInt(fmt.Sprintf("%v", q[name]), 10, 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetFloat32(name string, def ...float32) float32 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 32); err == nil {
		return float32(value)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetFloat64(name string, def ...float64) float64 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 64); err == nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetDatatime 获取时间字段
func (q QueryRow) GetDatatime(name string, format ...string) (time.Time, error) {
	t, b := q.GetMustString(name)
	if !b {
		return time.Now(), fmt.Errorf("%s列不存在", name)
	}
	f := "2006/01/02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	return time.ParseInLocation(f, t, time.Local)
}

//Has 检查对象中是否存在某个值
func (q QueryRow) Has(name string) bool {
	_, ok := q[name]
	return ok
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func (q QueryRow) GetMustString(name string) (string, bool) {
	if value, ok := q[name].(string); ok {
		return value, true
	}
	return "", false
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustInt(name string) (int, bool) {
	if value, err := strconv.Atoi(fmt.Sprintf("%v", q[name])); err == nil {
		return value, true
	}
	return 0, false
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustFloat32(name string) (float32, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 32); err == nil {
		return float32(value), true
	}
	return 0, false
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustFloat64(name string) (float64, bool) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 64); err == nil {
		return value, true
	}
	return 0, false
}

//ToStruct 将当前对象转换为指定的struct
func (q QueryRow) ToStruct(o interface{}) error {
	return types.Map2Struct(q, o)
}

//QueryRows 多行数据
type QueryRows []QueryRow

//ToStruct 将当前对象转换为指定的struct
func (q QueryRows) ToStruct(o interface{}) error {
	return types.Map2Struct(q, o)
}

//IsEmpty 当前数据集是否为空
func (q QueryRows) IsEmpty() bool {
	return q == nil || len(q) == 0
}

//Len 获取当前数据集的长度
func (q QueryRows) Len() int {
	return len(q)
}

//Get 获取指定索引的数据
func (q QueryRows) Get(i int) QueryRow {
	if q == nil || i >= len(q) || i < 0 {
		return QueryRow{}
	}
	return q[i]
}
