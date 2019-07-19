package component

import (
	"reflect"

	"github.com/micro-plat/hydra/context"
)

func isCorrectType(h interface{}) bool {
	if isConstructor(h) {
		return true
	}
	return isHandler(h)
}

func isConstructor(h interface{}) bool {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	if fv.Kind() != reflect.Func || tp.NumIn() > 1 || tp.NumOut() > 2 || tp.NumOut() == 0 {
		return false
	}
	if tp.NumIn() == 1 && tp.In(0).Name() == "IContainer" {
		return true
	}
	if tp.NumIn() == 0 {
		return true
	}
	return false
}
func isHandler(h interface{}) bool {
	// fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	if tp.NumIn() != 1 || tp.NumOut() != 1 {
		return false
	}
	switch h.(type) {
	case ServiceFunc, Handler:
		return true
	default:
		_, ok := h.(func(*context.Context) interface{})
		return ok
	}
}
