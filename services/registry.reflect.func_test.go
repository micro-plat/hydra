package services

import (
	"fmt"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestCreateObject(t *testing.T) {
	type v struct {
		Name string
	}
	obj := &v{Name: "abc"}
	cases := []struct {
		name  string
		input interface{}
		value interface{}
		err   error
	}{
		{name: "1.输出参数个数为0,输入参数个数为0", input: func() {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{name: "2.输出参数个数为0,输入参数个数为1", input: func(int) {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{name: "3.输出参数个数为0,输入参数个数为2", input: func(int, int) {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{name: "4.输入参数个数为0,第一个输出参数为int", input: func() int { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{name: "5.输入参数个数为0,第一个输出参数为MAP", input: func() map[string]string { return nil }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{name: "6.输入参数个数为0,第一个输出参数为string", input: func() string { return "" }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{name: "7.输入参数个数为0,第一个输出参数为float", input: func() float32 { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{name: "8.输入参数个数为0,第一个输出参数为int32", input: func() int32 { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{name: "9.输入参数个数为0,第一个输出参数为struct", input: func() v { return *obj }, value: *obj},
		{name: "10.输入参数个数为0,第一个输出参数为struct指针", input: func() *v { return obj }, value: obj},
		{name: "11.输入参数个数为0,第二个输出参数为int", input: func() (*v, int) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{name: "12.输入参数个数为0,第二个输出参数为map", input: func() (*v, map[string]string) { return obj, nil }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{name: "13.输入参数个数为0,第二个输出参数为string", input: func() (*v, string) { return obj, "" }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{name: "14.输入参数个数为0,第二个输出参数为float32", input: func() (*v, float32) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{name: "15.输入参数个数为0,第二个输出参数为int32", input: func() (*v, int32) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{name: "16.输入参数个数为0,第二个输出参数为非空的error", input: func() (*v, error) { return obj, fmt.Errorf("无法创建对象") }, err: fmt.Errorf("无法创建对象")},
		{name: "17.输入参数个数为0,第二个输出参数为空的error", input: func() (*v, error) { return obj, nil }, value: obj},
	}
	for _, k := range cases {
		v, err := createObject(k.input)
		assert.Equal(t, err, k.err, k.name)
		assert.Equal(t, v, k.value, k.name)
	}
}
