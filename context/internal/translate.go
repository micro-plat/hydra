package internal

import (
	"github.com/d5/tengo/v2"
	"github.com/micro-plat/lib4go/types"
)

//Translate 翻译带参数的变量支持格式有 @abc,{@abc}
func Translate(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	//获取第二个参数
	format, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	//获取第一个参数
	second := tengo.ToInterface(args[1])

	mapArgs, ok := second.(map[string]interface{})
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "map[string]interface{}",
			Found:    args[1].TypeName(),
		}
	}
	var xm types.XMap = mapArgs
	result := xm.Translate(format)

	return &tengo.String{Value: result}, nil
}
