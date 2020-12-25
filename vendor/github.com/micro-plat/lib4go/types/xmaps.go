package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

//XMaps 多行数据
type XMaps []XMap

//NewXMaps 构建xmap对象
func NewXMaps(len ...int) XMaps {
	return make(XMaps, 0, GetIntByIndex(len, 0, 1))
}

//NewXMapsByJSON 根据json创建XMaps
func NewXMapsByJSON(j string) (XMaps, error) {
	var query XMaps
	d := json.NewDecoder(bytes.NewBuffer(StringToBytes(j)))
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
	DeepCopy(o, val.Interface())
	return nil
}

//ToAnyStructs 将当前对象转换为指定的struct
func (q XMaps) ToAnyStructs(o interface{}) error {
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
		if err := r.ToAnyStruct(mVal.Interface()); err != nil {
			return err
		}
		val = reflect.Append(val, mVal)
	}
	DeepCopy(o, val.Interface())
	return nil
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
