package internal

import "github.com/d5/tengo/v2"

//IASANY  将返回结果为interface{}转为任意类型
func IASANY(fn func() interface{}) tengo.CallableFunc {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		s := fn()
		obj, err := tengo.FromInterface(s)
		if err != nil {
			return nil, err
		}

		return obj, nil
	}
}
