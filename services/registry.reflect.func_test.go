package services

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func TestCreateObject(t *testing.T) {
	type v struct {
		Name string
	}
	obj := &v{Name: "abc"}
	cases := []struct {
		input interface{}
		value interface{}
		err   error
	}{
		{input: func() {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{input: func(int) {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{input: func(int, int) {}, err: fmt.Errorf("函数的输入参数或输出参数个数错误")},
		{input: func() int { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{input: func() map[string]string { return nil }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{input: func() string { return "" }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{input: func() float32 { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{input: func() int32 { return 0 }, err: fmt.Errorf("输出参数第一个参数必须是结构体")},
		{input: func() v { return *obj }, value: *obj},
		{input: func() *v { return obj }, value: obj},
		{input: func() (*v, int) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{input: func() (*v, map[string]string) { return obj, nil }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{input: func() (*v, string) { return obj, "" }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{input: func() (*v, float32) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{input: func() (*v, int32) { return obj, 0 }, err: fmt.Errorf("第二个输出参数必须为error类型")},
		{input: func() (*v, error) { return obj, fmt.Errorf("无法创建对象") }, err: fmt.Errorf("无法创建对象")},
		{input: func() (*v, error) { return obj, nil }, value: obj},
	}
	for _, k := range cases {
		v, err := createObject(k.input)
		assert.Equal(t, err, k.err)
		assert.Equal(t, v, k.value)
	}
}
