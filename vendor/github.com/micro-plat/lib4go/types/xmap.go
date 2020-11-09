package types

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

var _ IXMap = XMap{}

//IXMap 扩展map
type IXMap interface {
	Keys() []string
	Get(name string) (interface{}, bool)
	GetValue(name string) interface{}
	GetString(name string, def ...string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	GetStrings(name string, def ...string) (r []string)
	GetArray(name string, def ...interface{}) (r []interface{})
	GetDatetime(name string, format ...string) (time.Time, error)
	SetValue(name string, value interface{})
	Has(name string) bool
	MustString(name string) (string, bool)
	MustInt(name string) (int, bool)
	MustFloat32(name string) (float32, bool)
	MustFloat64(name string) (float64, bool)

	Marshal() []byte
	GetJSON(name string) (r []byte, err error)
	IsXMap(name string) bool
	GetXMap(name string) (c XMap, err error)

	IsEmpty() bool
	Len() int
	ToStruct(o interface{}) error
	ToMap() map[string]interface{}
	ToSMap() map[string]string

	Cascade(m IXMap)
	Merge(m IXMap)
	MergeMap(anr map[string]interface{})
	MergeSMap(anr map[string]string)
}

//XMap map扩展对象
type XMap map[string]interface{}

//NewXMap 构建xmap对象
func NewXMap(len ...int) XMap {
	return make(map[string]interface{}, GetIntByIndex(len, 0, 1))
}

//NewXMapByMap 根据map[string]interface{}构建xmap
func NewXMapByMap(i map[string]interface{}) XMap {
	return i
}

//NewXMapBySMap  根据map[string]string构建xmap
func NewXMapBySMap(i map[string]string) XMap {
	return GetIMap(i)
}

//NewXMapByJSON 根据json创建XMap
func NewXMapByJSON(j string) (XMap, error) {
	var query XMap
	d := json.NewDecoder(bytes.NewBuffer([]byte(j)))
	d.UseNumber()
	err := d.Decode(&query)
	return query, err
}

//Merge 合并
func (q XMap) Merge(m IXMap) {
	keys := m.Keys()
	for _, key := range keys {
		q.SetValue(key, m.GetValue(key))
	}
}

//Cascade 对map进行级联累加，即将多级map转化为一级map,key使用"."进行边拉
func (q XMap) Cascade(m IXMap) {
	keys := m.Keys()
	for _, key := range keys {
		m := GetCascade(key, m.GetValue(key))
		q.Merge(XMap(m))
	}
}

//Keys 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) Keys() []string {
	keys := make([]string, len(q))
	idx := 0
	for k := range q {
		keys[idx] = k
		idx++
	}
	return keys
}

//IsEmpty 当前对象未包含任何数据
func (q XMap) IsEmpty() bool {
	return q == nil || len(q) == 0
}

//Len 获取当前对象包含的键值对个数
func (q XMap) Len() int {
	return len(q)
}

//Get 获取指定元素的值
func (q XMap) Get(name string) (interface{}, bool) {
	v, ok := q[name]
	return v, ok
}

//GetValue 获取指定参数的值
func (q XMap) GetValue(name string) interface{} {
	return q[name]
}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) GetString(name string, def ...string) string {
	parties := strings.Split(name, ":")
	if len(parties) == 1 {
		return GetString(q[name], def...)
	}
	tmpv := q[parties[0]]
	for i, cnt := 1, len(parties); i < cnt; i++ {
		if v, ok := tmpv.(map[string]interface{}); ok {
			tmpv = v[parties[i]]
			continue
		}
		if v, ok := tmpv.(XMap); ok {
			tmpv = v[parties[i]]
			continue
		}
		if v, ok := tmpv.(*XMap); ok {
			tmpv = v.GetValue(parties[i])
			continue
		}
		if v, ok := tmpv.(string); ok {
			tmp := map[string]interface{}{}
			json.Unmarshal([]byte(v), &tmp)
			tmpv = tmp[parties[i]]
			continue
		}
	}
	return GetString(tmpv, def...)
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

//GetStrings 获取字符串数组
func (q XMap) GetStrings(name string, def ...string) (r []string) {
	if v := q.GetString(name); v != "" {
		if r = strings.Split(v, ";"); len(r) > 0 {
			return r
		}
	}
	if len(def) > 0 {
		return def
	}
	return nil
}

//GetArray 获取数组对象
func (q XMap) GetArray(name string, def ...interface{}) (r []interface{}) {
	v, ok := q.Get(name)
	if !ok && len(def) > 0 {
		return def
	}
	if r, ok := v.([]interface{}); ok {
		return r
	}
	return nil
}

//Marshal 转换为json数据
func (q XMap) Marshal() []byte {
	r, _ := json.Marshal(q)
	return r
}

//GetJSON 获取JSON串
func (q XMap) GetJSON(name string) (r []byte, err error) {
	v, ok := q.Get(name)
	if !ok {
		return nil, fmt.Errorf("%s不存在或值为空", name)
	}

	buffer, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

//IsXMap 是否存在节点
func (q XMap) IsXMap(name string) bool {
	v, ok := q.Get(name)
	if !ok {
		return false
	}
	_, ok = v.(map[string]interface{})
	return ok
}

//GetXMap 指定节点名称获取JSONConf
func (q XMap) GetXMap(name string) (c XMap, err error) {
	v, ok := q.Get(name)
	if !ok {
		err = fmt.Errorf("%s不存在或值为空", name)
		return
	}
	if data, ok := v.(map[string]interface{}); ok {
		return data, nil
	}
	return nil, fmt.Errorf("%s不是有效的map", name)
}

//SetValue 获取时间字段
func (q XMap) SetValue(name string, value interface{}) {
	q[name] = value
}

//Has 检查对象中是否存在某个值
func (q XMap) Has(name string) bool {
	_, ok := q.Get(name)
	return ok
}

//MustString 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) MustString(name string) (string, bool) {
	return MustString(q[name])
}

//MustInt 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustInt(name string) (int, bool) {
	return MustInt(q[name])
}

//MustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustFloat32(name string) (float32, bool) {
	return MustFloat32(q[name])
}

//MustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustFloat64(name string) (float64, bool) {
	return MustFloat64(q[name])
}

//ToStruct 将当前对象转换为指定的struct
func (q XMap) ToStruct(out interface{}) error {
	buff, err := json.Marshal(q)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buff, &out); err != nil {
		return err
	}
	return nil
}

//ToMap 转换为map[string]interface{}
func (q XMap) ToMap() map[string]interface{} {
	return q
}

//ToSMap 转换为map[string]string
func (q XMap) ToSMap() map[string]string {
	v, _ := ToStringMap(q)
	return v
}

//MergeMap 将传入的xmap合并到当前xmap
func (q XMap) MergeMap(anr map[string]interface{}) {
	for k, v := range anr {
		q.SetValue(k, v)
	}
}

//MergeSMap 将传入的xmap合并到当前xmap
func (q XMap) MergeSMap(anr map[string]string) {
	for k, v := range anr {
		q.SetValue(k, v)
	}
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

//MarshalXML 转换为xml字符串
func (q XMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(q) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range q {
		if v == nil || GetString(v) == "" {
			continue
		}
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: GetString(v)})
	}

	return e.EncodeToken(start.End())
}

//UnmarshalXML xml转换为xmap
func (q XMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if q == nil {
		q = XMap{}
	}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(q)[e.XMLName.Local] = e.Value
	}
	return nil
}

//XMaps 多行数据
type XMaps []XMap

//NewXMaps 构建xmap对象
func NewXMaps(len ...int) XMaps {
	return make(XMaps, 0, GetIntByIndex(len, 0, 1))
}

//NewXMapsByJSON 根据json创建XMaps
func NewXMapsByJSON(j string) (XMaps, error) {
	var query XMaps
	d := json.NewDecoder(bytes.NewBuffer([]byte(j)))
	d.UseNumber()
	err := d.Decode(&query)
	return query, err
}

//Append 追加xmap
func (q *XMaps) Append(i ...XMap) XMaps {
	*q = append(*q, i...)
	return *q
}

//ToStructs 将当前对象转换为指定的struct
func (q XMaps) ToStructs(o interface{}) error {
	fval := reflect.ValueOf(o)
	if fval.Kind() == reflect.Interface || fval.Kind() == reflect.Ptr {
		fval = fval.Elem()
	} else {
		return fmt.Errorf("输入参数必须是指针:%v", fval.Kind())
	}
	// we only accept structs
	if fval.Kind() != reflect.Slice {
		return fmt.Errorf("传入参数错误，必须是切片类型:%v", fval.Kind())
	}
	val := reflect.Indirect(reflect.ValueOf(o))
	typ := val.Type()
	for _, r := range q {
		mVal := reflect.Indirect(reflect.New(typ.Elem().Elem())).Addr()
		if err := r.ToStruct(mVal.Interface()); err != nil {
			return err
		}
		val = reflect.Append(val, mVal)
	}
	deepCopy(o, val.Interface())
	return nil
}
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
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

//Copy 拷贝一个新的map,并追加新的键值对
func Copy(input map[string]interface{}, kv ...string) XMap {
	nmap := make(map[string]interface{}, len(input))
	for k, v := range input {
		nmap[k] = v
	}
	if len(kv) == 0 || len(kv)%2 != 0 {
		return nmap
	}
	for i := 0; i < len(kv)/2; i++ {
		nmap[kv[i]] = kv[i+1]
	}
	return nmap
}

//GetCascade 根据key将值转换为map[string]ineterface{}
func GetCascade(key string, value interface{}) map[string]interface{} {
	nmap := make(map[string]interface{})
	switch vlu := value.(type) {
	case string, []byte, int, int8, int32, int64, uint,
		uint8, uint32, uint64, float32, float64, time.Time,
		bool, complex64, complex128:
		nmap[key] = vlu
		return nmap
	case map[string]string:
		for k, v := range vlu {
			nmap[fmt.Sprintf("%s.%s", key, k)] = v
		}
		return nmap
	case map[string]interface{}:
		for k, v := range vlu {
			n := GetCascade(k, v)
			for a, b := range n {
				nmap[fmt.Sprintf("%s.%s", key, a)] = b
			}
		}
		return nmap
	default:

		m, err := IToMap(value)
		if err != nil {
			nmap[key] = value
			return nmap
		}
		return GetCascade(key, m)

	}
}

//IToMap struct类型转map[string]interface{}
func IToMap(o interface{}) (map[string]interface{}, error) {
	if o == nil {
		return nil, nil
	}
	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Map {
		switch v := o.(type) {
		case map[string]interface{}:
			return v, nil
		case map[string]string:
			return GetIMap(v), nil
		}
	}
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		buff, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		out := make(map[string]interface{})
		err = json.Unmarshal(buff, &out)
		return out, err
	case reflect.Slice:
		nmap := make(map[string]interface{})
		for i := 0; i < val.Len(); i++ {
			v := GetCascade(fmt.Sprint(i), val.Index(i).Interface())
			for a, b := range v {
				nmap[a] = b
			}
		}
		return nmap, nil
	default:
		return nil, fmt.Errorf("输入参数类型错误 accepts structs; got %s", val.Kind())
	}

}
