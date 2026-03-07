package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type IXMaps interface {
	ToStructs(o interface{}) error
	Append(i ...IXMap)
	ToAnyStructs(o interface{}) error
	Maps() []XMap
	IsEmpty() bool
	Len() int
	Get(i int) IXMap
}

var _ IXMaps = &XMaps{}

// XMaps 多行数据
type XMaps []XMap

// NewXMaps 构建xmap对象
func NewXMaps(len ...int) XMaps {
	v := make(XMaps, 0, GetIntByIndex(len, 0))
	return v
}

// NewXMapsByJSON 根据json创建XMaps
func NewXMapsByJSON(j string) (XMaps, error) {
	var query = make(XMaps, 0, 1)
	d := json.NewDecoder(bytes.NewBuffer(StringToBytes(j)))
	d.UseNumber()
	err := d.Decode(&query)
	return query, err
}

// Append 追加xmap
func (q *XMaps) Append(i ...IXMap) {
	for _, v := range i {
		*q = append(*q, v.ToMap())
	}
	return
}

// ToStructs 将当前对象转换为指定的struct
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
	DeepCopy(o, val.Interface())
	return nil
}

// ToStructsT 将XMaps转换为指定类型的切片
func ToStructsT[T any](q XMaps) ([]T, error) {
	if q.IsEmpty() {
		return nil, nil
	}

	result := make([]T, 0, len(q))
	for _, xmap := range q {
		var t T
		// 利用反射检查T是否为指针类型，并正确处理
		rt := reflect.TypeOf(t)
		var target interface{}
		if rt.Kind() == reflect.Ptr {
			// T是指针类型，创建指向新实例的指针
			elem := reflect.New(rt.Elem()).Interface()
			target = elem
		} else {
			// T是结构体类型，取地址
			target = &t
		}

		// 使用XMap的ToAnyStruct方法填充目标结构体
		if err := xmap.ToAnyStruct(target); err != nil {
			return nil, err
		}

		// 根据T的类型决定如何存储到结果切片
		if rt.Kind() == reflect.Ptr {
			result = append(result, target.(T))
		} else {
			result = append(result, reflect.ValueOf(target).Elem().Interface().(T))
		}
	}
	return result, nil
}

// ToAnyStructs 将当前对象转换为指定的struct
func (q XMaps) ToAnyStructs(o interface{}) error {
	fval := reflect.ValueOf(o)
	if fval.Kind() != reflect.Ptr {
		return fmt.Errorf("输入参数必须是指针: %v", fval.Kind())
	}

	sliceVal := fval.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("传入参数错误，必须是切片类型: %v", sliceVal.Kind())
	}

	elemType := sliceVal.Type().Elem() // 获取切片元素的类型
	for _, r := range q {
		var elemVal reflect.Value
		if elemType.Kind() == reflect.Ptr {
			// 如果元素是指针类型，创建一个新实例并取其地址
			newElem := reflect.New(elemType.Elem())
			elemVal = newElem
		} else {
			// 如果元素是结构体类型，创建一个新实例并取其地址（用于调用 ToAnyStruct）
			newElem := reflect.New(elemType)
			elemVal = newElem.Elem().Addr() // 取结构体的指针
		}

		if err := r.ToAnyStruct(elemVal.Interface()); err != nil {
			return err
		}

		// 将新元素追加到切片中
		sliceVal.Set(reflect.Append(sliceVal, elemVal))
	}
	return nil
}

// Maps map列表
func (q XMaps) Maps() []XMap {
	return q
}

// IsEmpty 当前数据集是否为空
func (q XMaps) IsEmpty() bool {
	return len(q) == 0
}

// Len 获取当前数据集的长度
func (q XMaps) Len() int {
	return len(q)
}

// Get 获取指定索引的数据
func (q XMaps) Get(i int) IXMap {
	if q == nil || i >= len(q) || i < 0 {
		return XMap{}
	}
	return q[i]
}
