package types

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

//Any2XML 将任意类型转换为xml
func Any2XML(v interface{}, header string, root ...string) (string, error) {
	body, err := any2XML(v, root...)
	if err != nil {
		return "", err
	}
	return header + body, nil

}

//any2XML 将任意类型转换为xml
func any2XML(v interface{}, root ...string) (string, error) {
	//1. 处理类型，及空值
	value := reflect.ValueOf(v)
	vt := reflect.TypeOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		vt = vt.Elem()
	}
	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	//2. 根据类型分别处理
	var builder = &strings.Builder{}
	switch value.Kind() {
	case reflect.Struct:
		//检查空值
		if !reflect.ValueOf(value).IsValid() || reflect.ValueOf(value).IsZero() {
			return "", nil
		}
		//转为xml
		str, err := struct2xml(value, "xml", root...)
		if err != nil || str == "" {
			return "", err
		}
		if len(root) > 0 {
			builder.Write(StringToBytes(fmt.Sprintf("<%s>", GetStringByIndex(root, 0))))
		}
		builder.Write(StringToBytes(str))
	case reflect.Map:
		//检查空值
		if value.IsNil() {
			return "", nil
		}
		////转为xml
		str, err := map2xml(value, root...)
		if err != nil || str == "" {
			return "", err
		}
		if len(root) > 0 {
			builder.Write(StringToBytes(fmt.Sprintf("<%s>", GetStringByIndex(root, 0))))
		}
		builder.Write(StringToBytes(str))
	case reflect.Array, reflect.Slice:
		//检查空值
		if value.Len() == 0 {
			return "", nil
		}
		//转为xml
		str, err := slice2xml(value, root...)
		if err != nil || str == "" {
			return "", err
		}

		//输出xml----不输出外边的标签，与标准库一致
		builder.Write(StringToBytes(str))
		return builder.String(), nil
	case reflect.String:
		//检查空值
		if value.Len() == 0 {
			return "", nil
		}
		//输出xml
		if len(root) > 0 {
			builder.Write(StringToBytes(fmt.Sprintf("<%s>", GetStringByIndex(root, 0))))
		}
		builder.Write(StringToBytes(fmt.Sprint(value.Interface())))
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return "", nil
	default:
		//输出xml
		if len(root) > 0 {
			builder.Write(StringToBytes(fmt.Sprintf("<%s>", GetStringByIndex(root, 0))))
		}
		builder.Write(StringToBytes(fmt.Sprint(value.Interface())))
	}
	//输出xml结束标签
	if len(root) > 0 {
		builder.Write(StringToBytes(fmt.Sprintf("</%s>", GetStringByIndex(root, 0))))
	}
	return builder.String(), nil
}

func map2xml(m reflect.Value, root ...string) (string, error) {
	var builder = &strings.Builder{}
	keys := m.MapKeys()
	for _, key := range keys {
		value := m.MapIndex(key)
		if !reflect.ValueOf(value).IsValid() || reflect.ValueOf(value).IsZero() {
			continue
		}
		str, err := any2XML(value.Interface(), fmt.Sprint(key.Interface()))
		if err != nil {
			return "", err
		}
		if str == "" {
			continue
		}
		builder.Write(StringToBytes(str))
	}
	return builder.String(), nil
}
func slice2xml(v reflect.Value, root ...string) (string, error) {
	var builder = &strings.Builder{}
	len := v.Len()
	for i := 0; i < len; i++ {
		if i == 0 {
			builder.Write(StringToBytes(fmt.Sprintf("<%s>", GetStringByIndex(root, 0, "xml"))))
		}
		value := v.Index(i)
		if !value.IsValid() || value.IsZero() {
			continue
		}
		str, err := any2XML(value.Interface(), "item")
		if err != nil {
			return "", err
		}
		if str == "" {
			continue
		}
		builder.Write(StringToBytes(str))
	}
	if len > 0 {
		builder.Write(StringToBytes(fmt.Sprintf("</%s>", GetStringByIndex(root, 0, "xml"))))
	}
	return builder.String(), nil

}
func struct2xml(value reflect.Value, tag string, root ...string) (string, error) {
	var builder = &strings.Builder{}
	var newValue interface{}
	vt := reflect.TypeOf(value.Interface())
	if value.Type().String() == "time.Time" {
		//获取字段的标签值
		newValue = (value.Interface().(time.Time)).String()
		str, err := any2XML(newValue)
		if err != nil {
			return "", err
		}

		builder.Write(StringToBytes(str))
		return builder.String(), nil
	}
	for i := 0; i < value.NumField(); i++ {
		tfield := vt.Field(i)
		vfield := value.Field(i)
		if tfield.PkgPath != "" && !tfield.Anonymous { // unexported
			continue
		}
		if !vfield.IsValid() || vfield.IsZero() {
			continue
		}
		//获取字段的标签值
		tagName, _ := head(tfield.Tag.Get(tag), ",")
		if tagName == "-" || tagName == "" {
			continue
		}

		str, err := any2XML(vfield.Interface(), tagName)
		if err != nil {
			return "", err
		}
		if str == "" {
			continue
		}
		builder.Write(StringToBytes(str))
	}

	return builder.String(), nil
}
