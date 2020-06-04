package conf

import (
	"fmt"
	"reflect"
	"time"

	"github.com/micro-plat/lib4go/types"
)

var _ IMeta = Meta{}

//IMeta 无数据接口
type IMeta interface {
	Keys() []string
	Get(name string) (interface{}, bool)
	GetString(name string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	Set(name string, value interface{})
	Has(name string) bool
	GetMustString(name string) (string, bool)
	GetMustInt(name string) (int, bool)
	GetMustFloat32(name string) (float32, bool)
	GetMustFloat64(name string) (float64, bool)
	GetDatetime(name string, format ...string) (time.Time, error)
	IsEmpty() bool
	Len() int
	ToStruct(o interface{}) error
	ToMap() map[string]interface{}
	MergeMap(anr map[string]interface{})
	MergeSMap(anr map[string]string)
}

//Meta 元数据
type Meta map[string]interface{}

//NewMeta 构建元数据
func NewMeta() Meta {
	return make(map[string]interface{})
}

//Keys 从对象中获取数据值，如果不是字符串则返回空
func (q Meta) Keys() []string {
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	return keys
}

//IsEmpty 当前对象未包含任何数据
func (q Meta) IsEmpty() bool {
	return q == nil || len(q) == 0
}

//Len 获取当前对象包含的键值对个数
func (q Meta) Len() int {
	return len(q)
}

//Get 获取指定元素的值
func (q Meta) Get(name string) (interface{}, bool) {
	v, ok := q[name]
	return v, ok
}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (q Meta) GetString(name string) string {
	return types.GetString(q[name])
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetInt(name string, def ...int) int {
	return types.GetInt(q[name], def...)
}

//GetInt64 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetInt64(name string, def ...int64) int64 {
	return types.GetInt64(q[name], def...)
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetFloat32(name string, def ...float32) float32 {
	return types.GetFloat32(q[name], def...)
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetFloat64(name string, def ...float64) float64 {
	return types.GetFloat64(q[name], def...)
}

//GetBool 从对象中获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func (q Meta) GetBool(name string, def ...bool) bool {
	return types.GetBool(q[name], def...)
}

//GetDatetime 获取时间字段
func (q Meta) GetDatetime(name string, format ...string) (time.Time, error) {
	return types.GetDatetime(q[name], format...)
}

//Set 获取时间字段
func (q Meta) Set(name string, value interface{}) {
	q[name] = value
}

//Has 检查对象中是否存在某个值
func (q Meta) Has(name string) bool {
	_, ok := q[name]
	return ok
}

//GetMustString 从对象中获取数据值，如果不是字符串则返回空
func (q Meta) GetMustString(name string) (string, bool) {
	return types.MustString(q[name])
}

//GetMustInt 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetMustInt(name string) (int, bool) {
	return types.MustInt(q[name])
}

//GetMustFloat32 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetMustFloat32(name string) (float32, bool) {
	return types.MustFloat32(q[name])
}

//GetMustFloat64 从对象中获取数据值，如果不是字符串则返回0
func (q Meta) GetMustFloat64(name string) (float64, bool) {
	return types.MustFloat64(q[name])
}

//ToStruct 将当前对象转换为指定的struct
func (q Meta) ToStruct(o interface{}) error {
	fval := reflect.ValueOf(o)
	if fval.Kind() != reflect.Ptr {
		return fmt.Errorf("输入参数必须是指针:%v", fval.Kind())
	}
	return types.Map2Struct(q, o)
}

//ToMap 转换为map[string]interface{}
func (q Meta) ToMap() map[string]interface{} {
	return q
}

//ToSMap 转换为map[string]string
func (q Meta) ToSMap() map[string]string {
	v, _ := types.ToStringMap(q)
	return v
}

//MergeMap 将传入的Meta合并到当前Meta
func (q Meta) MergeMap(anr map[string]interface{}) {
	for k, v := range anr {
		q.Set(k, v)
	}
}

//MergeSMap 将传入的Meta合并到当前Meta
func (q Meta) MergeSMap(anr map[string]string) {
	for k, v := range anr {
		q.Set(k, v)
	}
}
