package db

import (
	"fmt"
	"strconv"
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
func (q QueryRow) GetFloat32(name string) float32 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 32); err == nil {
		return float32(value)
	}
	return 0
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetFloat64(name string) float64 {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 64); err == nil {
		return value
	}
	return 0
}

//Has 检查对象中是否存在某个值
func (q QueryRow) Has(name string) bool {
	_, ok := q[name]
	return ok
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func (q QueryRow) GetMustString(name string) (string, error) {
	if value, ok := q[name].(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("不存在列:%s", name)
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustInt(name string) (int, error) {
	if value, err := strconv.Atoi(fmt.Sprintf("%v", q[name])); err == nil {
		return value, nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustFloat32(name string) (float32, error) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 32); err == nil {
		return float32(value), nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q QueryRow) GetMustFloat64(name string) (float64, error) {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%v", q[name]), 64); err == nil {
		return value, nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}
