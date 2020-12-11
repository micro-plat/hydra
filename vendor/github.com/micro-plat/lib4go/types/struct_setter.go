package types

import (
	"fmt"
	"reflect"
)

//Map2Struct 将map转换为struct
func Map2Struct(v interface{}, input map[string]interface{}, tag string) (err error) {

	//检查输入对象的类型
	value := reflect.ValueOf(v)
	vt := reflect.TypeOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		vt = vt.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("输出对象必须为struct:%s", value.Kind().String())
	}

	//循环每个字段，对每个字段进行类型检查，并赋值
	for i := 0; i < value.NumField(); i++ {
		tfield := vt.Field(i)
		vfield := value.Field(i)
		if tfield.PkgPath != "" && !tfield.Anonymous { // unexported
			continue
		}

		//获取字段的标签值
		tagName, _ := head(tfield.Tag.Get(tag), ",")
		if tagName == "-" || tagName == "" {
			continue
		}

		//拿取实际值
		sourceValue, ok := input[tagName]
		if !ok {
			continue
		}

		//根据目标对象类型，处理数据源
		switch vfield.Kind() {
		case reflect.Array, reflect.Slice: //目标字段为数组，将数据源为非数据转化为数组
			rvalue := reflect.ValueOf(sourceValue)
			switch rvalue.Kind() {
			case reflect.Array, reflect.Slice: //值也为数组，构建新数组，直接并赋值
				array := make([]interface{}, 0, 1)
				for i := 0; i < rvalue.Len(); i++ {
					array = append(array, rvalue.Index(i).Interface())
				}
				if err = setSlice(array, value.Field(i), tfield); err != nil {
					return fmt.Errorf("向字段%s赋值失败%w，值是:%+v", tfield.Name, err, array)
				}
			case reflect.Chan, reflect.Func:
				return fmt.Errorf("无法将chan,func等放入数组字段：%s", tfield.Name)
			default:
				if err = setSlice([]interface{}{rvalue.Interface()}, value.Field(i), tfield); err != nil {
					return fmt.Errorf("向字段%s赋值失败%w，值是:%+v", tfield.Name, err, []interface{}{rvalue.Interface()})
				}
			}
		case reflect.Struct:
			v, ok := sourceValue.(map[string]interface{})
			if !ok {
				if err = setWithProperType(sourceValue, vfield, tfield); err != nil {
					return fmt.Errorf("向字段%s赋值失败%w，值是:%+v", tfield.Name, err, sourceValue)
				}
				continue
			}
			return Map2Struct(vfield.Addr().Interface(), v, tag)
		default: //目标字段为非数组，直接符值
			if err = setWithProperType(sourceValue, vfield, tfield); err != nil {
				return fmt.Errorf("向字段%s赋值失败%w，值是:%+v", tfield.Name, err, sourceValue)
			}
		}
	}
	return nil

}
