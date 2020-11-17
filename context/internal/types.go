package internal

import (
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/micro-plat/lib4go/types"
)

func GetStringByIndex(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	//获取第一个参数
	inputArray := make([]string, 0, 2)
	switch arg0 := args[0].(type) {
	case *tengo.Array:
		for idx, a := range arg0.Value {
			as, ok := tengo.ToString(a)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("first[%d]", idx),
					Expected: "string(compatible)",
					Found:    a.TypeName(),
				}
			}
			inputArray = append(inputArray, as)
		}
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}

	//获取第二个参数
	index, ok := tengo.ToInt(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int",
			Found:    args[0].TypeName(),
		}
	}
	value := types.GetStringByIndex(inputArray, index)
	return &tengo.String{Value: value}, nil
}

func GetIntByIndex(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	//获取第一个参数
	inputArray := make([]int, 0, 2)
	switch arg0 := args[0].(type) {
	case *tengo.Array:
		for idx, a := range arg0.Value {
			as, ok := tengo.ToInt(a)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("first[%d]", idx),
					Expected: "int(compatible)",
					Found:    a.TypeName(),
				}
			}
			inputArray = append(inputArray, as)
		}
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}

	//获取第二个参数
	index, ok := tengo.ToInt(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int",
			Found:    args[0].TypeName(),
		}
	}
	value := types.GetIntByIndex(inputArray, index)
	return &tengo.Int{Value: int64(value)}, nil
}

func GetFloatByIndex(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	//获取第一个参数
	inputArray := make([]float64, 0, 2)
	switch arg0 := args[0].(type) {
	case *tengo.Array:
		for idx, a := range arg0.Value {
			as, ok := tengo.ToFloat64(a)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("first[%d]", idx),
					Expected: "float64(compatible)",
					Found:    a.TypeName(),
				}
			}
			inputArray = append(inputArray, as)
		}
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}

	//获取第二个参数
	index, ok := tengo.ToInt(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int",
			Found:    args[0].TypeName(),
		}
	}
	value := types.GetFloat64ByIndex(inputArray, index)
	return &tengo.Float{Value: value}, nil
}

//exclude  从数组中排除指定的值
func Exclude(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	//获取第一个参数
	inputArray := make([]string, 0, 2)
	switch arg0 := args[0].(type) {
	case *tengo.Array:
		for idx, a := range arg0.Value {
			as, ok := tengo.ToString(a)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("first[%d]", idx),
					Expected: "string(compatible)",
					Found:    a.TypeName(),
				}
			}
			inputArray = append(inputArray, as)
		}
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}

	//获取第二个参数
	second, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	//排除指定的值
	array := &tengo.Array{}
	for _, v := range inputArray {
		if len(v) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		if v != second {
			array.Value = append(array.Value, &tengo.String{Value: v})
		}
	}

	return array, nil
}
