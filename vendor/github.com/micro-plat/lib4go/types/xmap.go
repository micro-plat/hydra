package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/clbanning/mxj"
)

var _ IXMap = XMap{}

//IXMap 扩展map
type IXMap interface {
	//Keys map中的所有键名
	Keys() []string

	//Get 获取指定键对应的值，当存在时第二个值返回false
	Get(name string) (interface{}, bool)

	//GetValue 获取键对应的值，当值不存在时返回nil
	GetValue(name string) interface{}

	//Append 添加键值对，输入参数以:键，值，键，值...的顺序传入
	Append(kv ...interface{})

	//SetValue 设置值
	SetValue(name string, value interface{})

	//GetString 获取字符串
	GetString(name string, def ...string) string

	//GetInt 获取类型为int的值
	GetInt(name string, def ...int) int

	//GetInt32 获取类型为int32的值
	GetInt32(name string, def ...int32) int32

	//GetInt64 获取类型为int64的值
	GetInt64(name string, def ...int64) int64

	//GetFloat32 获取类型为float32的值
	GetFloat32(name string, def ...float32) float32

	//GetFloat64 获取类型为float64的值
	GetFloat64(name string, def ...float64) float64

	//GetDecimal 获取类型为Decimal的值
	GetDecimal(name string, def ...Decimal) Decimal

	//GetDatetime 获取日期类型的值
	GetDatetime(name string, format ...string) (time.Time, error)

	//GetStrings 获取值为[]string类型的值
	GetStrings(name string, def ...string) (r []string)

	//GetArray 获取值为数组类型的值
	GetArray(name string, def ...interface{}) (r []interface{})

	//GetBool 获取bool类型的值
	GetBool(name string, def ...bool) bool

	//Has 是否包含键
	Has(name string) bool

	//MustString 值是否是string并返回相关的值与判断结果
	MustString(name string) (string, bool)

	//MustInt 值是否是int并返回相关的值与判断结果
	MustInt(name string) (int, bool)

	//MustInt32 值是否是int32并返回相关的值与判断结果
	MustInt32(name string) (int32, bool)

	//MustInt64 值是否是int64并返回相关的值与判断结果
	MustInt64(name string) (int64, bool)

	//MustFloat32  值是否是int并返回相关的值与判断结果
	MustFloat32(name string) (float32, bool)

	//MustFloat64   值是否是int并返回相关的值与判断结果
	MustFloat64(name string) (float64, bool)

	//Delete 删除指定名称的值
	Delete(name string)

	//Marshal 将当前对转转换为json
	Marshal() []byte

	//GetJSON 将指定的键对应的值转换为json
	GetJSON(name string) (r []byte, err error)

	//IsXMap 指定的键是否是map[string]interface{}类型
	IsXMap(name string) bool

	//GetXMap 将指定的键的值转换为xmap
	GetXMap(name string) (c XMap)

	//IsEmpty 是否是空结构
	IsEmpty() bool

	//Len 获取元素个数
	Len() int

	//ToKV 转换为可通过URL传递的键值对
	ToKV(ecoding ...string) string

	//ToStruct 将当前map转换为结构值对象
	ToStruct(o interface{}) error

	//ToAnyStruct 转换为任意struct,struct中无须设置数据类型(性能较差)
	ToAnyStruct(out interface{}) error

	//ToMap 转换为map[string]interface{}
	ToMap() map[string]interface{}

	//ToSMap 转换为map[string]string
	ToSMap() map[string]string

	//Cascade 将多层map转换为单层map
	Cascade(m IXMap)

	//Translate 翻译带参数的变量支持格式有 @abc,{@abc}
	Translate(format string) string

	//Merge 合并多个xmap
	Merge(m IXMap)

	//Each 循环器，传入处理函数，内部循环每个数据并调用处理函数
	Each(fn func(string, interface{}))

	//Iterator 迭代处理器，传入处理函数，函数返回结果为新值，新建新的map并返回
	Iterator(fn func(string, interface{}) interface{}) XMap

	//Count 计数器，传入处理函数，函数返回值为true则为需要计数，最后返回符合条件的数量和
	Count(fn func(string, interface{}) bool) int

	//Filter 过滤器，传入过滤函数，函数返回值为true则为需要的参数装入map返回
	Filter(fn func(string, interface{}) bool) XMap

	//MergeMap 合并map[string]interface{}
	MergeMap(anr map[string]interface{})

	//MergeSMap 合并map[string]string
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
	n := make(map[string]interface{})
	for k, v := range i {
		n[k] = v
	}
	return n
}

//NewXMapByJSON 根据json创建XMap
func NewXMapByJSON(j string) (XMap, error) {
	var query XMap
	d := json.NewDecoder(bytes.NewBuffer(StringToBytes(j)))
	d.UseNumber()
	err := d.Decode(&query)
	return query, err
}

//NewXMapByXML 将xml转换为xmap
func NewXMapByXML(j string) (XMap, error) {
	data := make(map[string]interface{})
	mxj.PrependAttrWithHyphen(false) //修改成可以转换成多层map
	m, err := mxj.NewMapXml(StringToBytes(j))
	if err != nil {
		return nil, err
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("xml根节点错误:%s", j)
	}
	root := ""
	for k := range m {
		root = k
	}
	value := reflect.ValueOf(m[root])
	if value.Kind() != reflect.Map {
		data = m
		return data, nil
	}
	for _, key := range value.MapKeys() {
		data[GetString(key)] = value.MapIndex(key).Interface()
	}
	return data, nil
}

//Merge 合并
func (q XMap) Merge(m IXMap) {
	keys := m.Keys()
	for _, key := range keys {
		q.SetValue(key, m.GetValue(key))
	}
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

//Cascade 对map进行级联累加，即将多级map转化为一级map,key使用"."进行边拉
func (q XMap) Cascade(m IXMap) {
	keys := m.Keys()
	for _, key := range keys {
		m := GetCascade(key, m.GetValue(key))
		q.Merge(XMap(m))
	}
}

//Append 追加键值对
func (q XMap) Append(kv ...interface{}) {
	if len(kv) == 0 || len(kv)%2 != 0 {
		return
	}
	for i := 0; i < len(kv); i = i + 2 {
		q.SetValue(fmt.Sprint(kv[i]), kv[i+1])
	}
	return
}

//Keys 从对象中获取数据值，如果不是字符串则返回空
func (q XMap) Keys() []string {
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
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
	return GetString(q[name], def...)
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetInt(name string, def ...int) int {
	return GetInt(q[name], def...)
}

//GetInt32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) GetInt32(name string, def ...int32) int32 {
	return GetInt32(q[name], def...)
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

//GetDecimal 获取类型为Decimal的值
func (q XMap) GetDecimal(name string, def ...Decimal) Decimal {
	return GetDecimal(q[name], def...)
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
		if r = strings.Split(v, ","); len(r) > 0 {
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
	if !ok && len(def) > 0 || v == nil {
		return def
	}

	s := reflect.ValueOf(v)
	r = make([]interface{}, 0, s.Len())
	for i := 0; i < s.Len(); i++ {
		r = append(r, s.Index(i).Interface())
	}
	return r
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
func (q XMap) GetXMap(name string) (c XMap) {
	v, ok := q.Get(name)
	if !ok {
		return map[string]interface{}{}
	}
	switch value := v.(type) {
	case map[string]interface{}:
		return value
	case IXMap:
		return value.ToMap()
	case XMap:
		return value
	case map[string]string:
		return NewXMapBySMap(value)
	}
	return map[string]interface{}{}
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

//MustInt32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustInt32(name string) (int32, bool) {
	return MustInt32(q[name])
}

//MustInt64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustInt64(name string) (int64, bool) {
	return MustInt64(q[name])
}

//MustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustFloat32(name string) (float32, bool) {
	return MustFloat32(q[name])
}

//MustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q XMap) MustFloat64(name string) (float64, bool) {
	return MustFloat64(q[name])
}

//ToKV 转换为可通过URL传递的键值对
func (q XMap) ToKV(ecoding ...string) string {
	u := url.Values{}
	for k, v := range q {
		u.Set(k, fmt.Sprint(v))
	}
	return u.Encode()
}

//ToStruct 将当前对象转换为指定的struct
func (q XMap) ToStruct(out interface{}) error {
	buff, err := json.Marshal(q)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, &out)
	return err
}

//ToAnyStruct 转换为任意struct,struct中无须设置数据类型(性能较差)
func (q XMap) ToAnyStruct(out interface{}) error {
	return Map2Struct(out, q, "json")
}

//Each 循环器，传入处理函数，内部循环每个数据并调用处理函数
func (q XMap) Each(fn func(string, interface{})) {
	for k, v := range q {
		fn(k, v)
	}
}

//Delete 删除指定键名的值
func (q XMap) Delete(name string) {
	delete(q, name)
}

//Iterator 迭代处理器，传入处理函数，函数返回结果为新值，新建新的map并返回
func (q XMap) Iterator(fn func(string, interface{}) interface{}) XMap {
	n := NewXMap()
	for k, v := range q {
		nv := fn(k, v)
		n.SetValue(k, nv)
	}
	return n
}

//Count 计数器，传入处理函数，函数返回值为true则为需要计数，最后返回符合条件的数量和
func (q XMap) Count(fn func(string, interface{}) bool) int {
	var n = 0
	for k, v := range q {
		if fn(k, v) {
			n++
		}
	}
	return n
}

//Filter 过滤器，传入过滤函数，函数返回值为true则为需要的参数装入map返回
func (q XMap) Filter(fn func(string, interface{}) bool) XMap {
	n := NewXMap()
	for k, v := range q {
		if fn(k, v) {
			n.SetValue(k, v)
		}
	}
	return n
}

//ToMap 转换为map[string]interface{}
func (q XMap) ToMap() map[string]interface{} {
	return q
}

//ToSMap 转换为map[string]string
func (q XMap) ToSMap() map[string]string {
	rmap := make(map[string]string)
	for k, v := range q {
		if s, ok := v.(string); ok {
			rmap[k] = s
		} else if _, ok := v.(float64); ok {
			rmap[k] = q.GetString(k)
		} else if s, ok := v.(interface{}); ok {
			buff, err := json.Marshal(s)
			if err != nil {
				rmap[k] = fmt.Sprint(v)
				continue
			}
			rmap[k] = string(buff)
		} else {
			rmap[k] = fmt.Sprint(v)
		}
	}
	return rmap
}

//Translate 翻译带参数的变量支持格式有 @abc,{@abc},转义符@
func (q XMap) Translate(format string) string {
	word := regexp.MustCompile(`[\w^@]*(\{@\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*\})`)
	result := word.ReplaceAllStringFunc(format, func(s string) string {
		if strings.HasPrefix(s, "@{@") {
			return s[1:]
		}
		return q.GetString(s[2 : len(s)-1])
	})
	word = regexp.MustCompile(`[\w^#]*(\{#\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*\})`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		if strings.HasPrefix(s, "#{#") {
			return s[1:]
		}
		return url.QueryEscape(q.GetString(s[2 : len(s)-1]))
	})

	word = regexp.MustCompile(`[\w^@{}]*(@\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*)`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		if strings.HasPrefix(s, "@@") {
			return s[1:]
		}
		if strings.HasPrefix(s, "{@") {
			return s
		}
		return q.GetString(s[1:])
	})

	word = regexp.MustCompile(`[\w^#{}]*(#\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*)`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		if strings.HasPrefix(s, "##") {
			return s[1:]
		}
		if strings.HasPrefix(s, "{#") {
			return s
		}
		return url.QueryEscape(q.GetString(s[1:]))
	})

	return result
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
			return NewXMapBySMap(v), nil
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
