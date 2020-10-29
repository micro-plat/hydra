package assert

// Testing helpers for doozer.

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/kr/pretty"
)

func assert(t *testing.T, result bool, f func(), cd int) {
	if !result {
		_, file, line, _ := runtime.Caller(cd + 1)
		t.Errorf("%s:%d", file, line)
		f()
		t.FailNow()
	}
}

func equal(t *testing.T, exp, got interface{}, cd int, args ...interface{}) {
	fn := func() {
		for _, desc := range pretty.Diff(exp, got) {
			t.Error("!", desc)
		}
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	result := reflect.DeepEqual(exp, got)
	assert(t, result, fn, cd+1)
}

func tt(t *testing.T, result bool, cd int, args ...interface{}) {
	fn := func() {
		t.Errorf("!  Failure")
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, result, fn, cd+1)
}

func T(t *testing.T, result bool, args ...interface{}) {
	tt(t, result, 1, args...)
}

func Tf(t *testing.T, result bool, format string, args ...interface{}) {
	tt(t, result, 1, fmt.Sprintf(format, args...))
}

func IsNil(t *testing.T, ext bool, got interface{}, args ...interface{}) {
	fn := func() {
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, reflect.DeepEqual(nil, got) == ext, fn, 1)
}

func IsNilf(t *testing.T, ext bool, got interface{}, format string, args ...interface{}) {
	fn := func() {
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, (got == nil) == ext, fn, 1)
}

func Equal(t *testing.T, exp, got interface{}, args ...interface{}) {
	equal(t, exp, got, 1, args...)
}

func Equalf(t *testing.T, exp, got interface{}, format string, args ...interface{}) {
	equal(t, exp, got, 1, fmt.Sprintf(format, args...))
}

func NotEqual(t *testing.T, exp, got interface{}, args ...interface{}) {
	fn := func() {
		t.Errorf("!  Unexpected: <%#v>", exp)
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	result := !reflect.DeepEqual(exp, got)
	assert(t, result, fn, 1)
}
func NotEqualf(t *testing.T, exp, got interface{}, format string, args ...interface{}) {
	fn := func() {
		t.Errorf("!  Unexpected: <%#v>", exp)
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprintf(format, args...))
		}
	}
	result := !reflect.DeepEqual(exp, got)
	assert(t, result, fn, 1)
}

//Panic Panic
func Panic(t *testing.T, expect interface{}, fn func(), args ...interface{}) {
	defer func() {
		equal(t, expect, recover(), 3, args...)
	}()
	fn()
}
