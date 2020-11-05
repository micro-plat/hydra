package services

import (
	"fmt"
	"reflect"
)

func createObject(f interface{}) (interface{}, error) {

	//验证函数签名
	tp := reflect.TypeOf(f)
	tv := reflect.ValueOf(f)
	if tp.NumIn() != 0 || tp.NumOut() == 0 || tp.NumOut() > 2 {
		return nil, fmt.Errorf("函数的输入参数或输出参数个数错误")
	}

	//检查第一个输出参数
	outKind := tp.Out(0).Kind()
	if outKind == reflect.Ptr {
		outKind = tp.Out(0).Elem().Kind()
	}
	if outKind != reflect.Struct {
		return nil, fmt.Errorf("输出参数第一个参数必须是结构体")
	}

	//检查第二个输出参数
	if tp.NumOut() == 2 {
		if tp.Out(1).Name() != "error" {
			return nil, fmt.Errorf("第二个输出参数必须为error类型")
		}
	}

	//调用函数
	rvalues := tv.Call(nil)

	//处理返回值
	if len(rvalues) > 1 {
		if err := rvalues[1].Interface(); err != nil {
			return nil, err.(error)
		}
	}
	return rvalues[0].Interface(), nil

}
