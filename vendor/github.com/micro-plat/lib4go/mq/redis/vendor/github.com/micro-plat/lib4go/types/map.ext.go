package types

import (
	"fmt"
	"time"
)

type ExtendMap map[string]interface{}

//NewExtendMap 构建包含扩展函数的 map[string]interface{}
func NewExtendMap(mp ...map[string]interface{}) ExtendMap {
	if len(mp) == 0 {
		return make(map[string]interface{})
	}
	if len(mp) == 1 {
		return mp[0]
	}
	nmap := make(map[string]interface{})
	for _, m := range mp {
		for k, v := range m {
			nmap[k] = v
		}
	}
	return nmap
}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q ExtendMap) GetString(name string) string {
	if value, ok := q[name].(string); ok {
		return value
	}
	return ""
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetInt(name string) int {
	if value, ok := q[name].(int); ok {
		return value
	}
	return 0
}

//GetTime 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetTime(name string) time.Time {
	if value, ok := q[name].(time.Time); ok {
		return value
	}
	return time.Time{}
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetFloat32(name string) float32 {
	if value, ok := q[name].(float32); ok {
		return value
	}
	return 0
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetFloat64(name string) float64 {
	if value, ok := q[name].(float64); ok {
		return value
	}
	return 0
}

//Has 检查对象中是否存在某个值
func (q ExtendMap) Has(name string) bool {
	_, ok := q[name]
	return ok
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func (q ExtendMap) GetMustString(name string) (string, error) {
	if value, ok := q[name].(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("不存在列:%s", name)
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetMustInt(name string) (int, error) {
	if value, ok := q[name].(int); ok {
		return value, nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}

//GetMustTime 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetMustTime(name string) (time.Time, error) {
	if value, ok := q[name].(time.Time); ok {
		return value, nil
	}
	return time.Time{}, fmt.Errorf("不存在列:%s", name)
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetMustFloat32(name string) (float32, error) {
	if value, ok := q[name].(float32); ok {
		return value, nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q ExtendMap) GetMustFloat64(name string) (float64, error) {
	if value, ok := q[name].(float64); ok {
		return value, nil
	}
	return 0, fmt.Errorf("不存在列:%s", name)
}
