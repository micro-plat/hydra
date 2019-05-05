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
			rmap[k] = fmt.Sprint(v)
		}
	}
	return rmap, nil

}

//Struct2Map 将struct 转换成map[string]interface{}
func Struct2Map(i interface{}) (map[string]interface{}, error) {
	buff, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	if err := json.Unmarshal(buff, &out); err != nil {
		return nil, err
	}
	return out, nil
}

//Map2Struct 将map转换成struct
func Map2Struct(i interface{}, o interface{}) error {
	config := &DecoderConfig{
		WeaklyTypedInput: true,
		Result:           o,
		TagName:          "m2s",
	}
	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(i)
}
