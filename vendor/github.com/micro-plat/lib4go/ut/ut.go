package ut

import (
	"fmt"
	"reflect"
)

type test interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

func Expect(t test, a interface{}, b interface{}) {
	a1 := a
	b1 := b
	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice:
		a1 = fmt.Sprintf("%+v", a)
	}
	val = reflect.ValueOf(b)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice:
		b1 = fmt.Sprintf("%+v", b)
	}
	if a1 != b1 {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b1, reflect.TypeOf(b1), a1, reflect.TypeOf(a1))
	}
}
func ExpectSkip(t test, a interface{}, b interface{}) bool {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
		t.FailNow()
		return true
	}
	return false
}

func Refute(t test, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
func RefuteSkip(t test, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
		t.FailNow()
	}
}
