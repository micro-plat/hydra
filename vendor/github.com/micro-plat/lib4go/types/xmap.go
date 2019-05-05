package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

var _ IXMap = XMap{}

type IXMap interface {
	GetString(name string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	SetValue(name string, value interface{})
	Has(name string) bool
	GetMustString(name string) (string, bool)
	GetMustInt(name string) (int, bool)
	GetMustFloat32(name string) (float32, bool)
	GetMustFloat64(name string) (float64, bool)
	GetDatetime(name string, format ...string) (time.Time, error)
	ToStruct(o interface{}) error
}

type XMap map[string]interface{}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) Keys() []string {
	keys := make([]string, len(q))
	idx := 0
	for k := range q {
		keys[idx] = k
		idx++
	}
	return keys
}

func (q XMap) GetValue(name string) interface{} {
	return q[name]
}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) GetString(name string) string {
	parties := strings.Split(name, ":")
	if len(parties) == 1 {
		return GetString(q[name])
	}
	tmpv := q[parties[0]]
	for i, cnt := 1, len(parties); i < cnt; i++ {
		if v, ok := tmpv.(map[string]interface{}); ok {
			tmpv = v[parties[i]]
			continue
		}
		if v, ok := tmpv.(string); ok {
			tmp := map[string]interface{}{}
			json.Unmarshal([]byte(v), &tmp)
			tmpv = tmp[parties[i]]
			continue
		}
	}
	return GetString(tmpv)
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetInt(name string, def ...int) int {
	return GetInt(q[name], def...)
}

//GetInt64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetInt64(name string, def ...int64) int64 {
	return GetInt64(q[name], def...)
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetFloat32(name string, def ...float32) float32 {
	return GetFloat32(q[name], def...)
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetFloat64(name string, def ...float64) float64 {
	return GetFloat64(q[name], def...)
}

//GetBool 从对象中获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func (q XMap) GetBool(name string, def ...bool) bool {
	return GetBool(q[name], def...)
}

//GetDatetime 获取时间字段
func (q XMap) GetDatetime(name string, format ...string) (time.Time, error) {
	return GetDatetime(q[name], format...)
}

//SetValue 获取时间字段
func (q XMap) SetValue(name string, value interface{}) {
	q[name] = value
}

//Has 检查对象中是否存在某个值
func (q XMap) Has(name string) bool {
	_, ok := q[name]
	return ok
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) GetMustString(name string) (string, bool) {
	return MustString(q[name])
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetMustInt(name string) (int, bool) {
	return MustInt(q[name])
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetMustFloat32(name string) (float32, bool) {
	return MustFloat32(q[name])
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetMustFloat64(name string) (float64, bool) {
	return MustFloat64(q[name])
}

//ToStruct 将当前对象转换为指定的struct
func (q XMap) ToStruct(o interface{}) error {
	input := make(map[string]interface{})
	for k, v := range q {
		input[k] = fmt.Sprint(v)
	}
	return Map2Struct(&input, &o)
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m *XMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(*m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range *m {
		if v == nil || GetString(v) == "" {
			continue
		}
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: GetString(v)})
	}

	return e.EncodeToken(start.End())
}

func (m *XMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if m == nil {
		m = &XMap{}
	}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

//XMaps 多行数据
type XMaps []XMap

//ToStruct 将当前对象转换为指定的struct
func (q XMaps) ToStruct(o interface{}) error {
	return Map2Struct(q, o)
}

//IsEmpty 当前数据集是否为空
func (q XMaps) IsEmpty() bool {
	return q == nil || len(q) == 0
}

//Len 获取当前数据集的长度
func (q XMaps) Len() int {
	return len(q)
}

//Get 获取指定索引的数据
func (q XMaps) Get(i int) XMap {
	if q == nil || i >= len(q) || i < 0 {
		return XMap{}
	}
	return q[i]
}

//ParseBool 将字符串转换为bool值
func ParseBool(val interface{}) (value bool, err error) {
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
