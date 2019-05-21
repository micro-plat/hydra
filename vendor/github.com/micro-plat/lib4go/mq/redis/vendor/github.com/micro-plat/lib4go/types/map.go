package types

import (
	"encoding/json"
	"fmt"
)

func GetIMap(m map[string]string) map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range m {
		n[k] = v
	}
	return n
}
func GetMapValue(key string, m map[string]string, def ...string) string {
	if v, ok := m[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}
func ToStringMap(input map[string]interface{}) (map[string]string, error) {
	rmap := make(map[string]string)
	for k, v := range input {
		if s, ok := v.(string); ok {
			rmap[k] = s
		} else if s, ok := v.(interface{}); ok {
			buff, err := json.Marshal(s)
			if err != nil {
				return nil, err
			}
			rmap[k] = string(buff)
		} else {
			return nil, fmt.Errorf("不支持的数据类型:%+v", v)
		}
	}
	return rmap, nil

}
