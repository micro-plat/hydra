package assert

// Testing helpers for doozer.

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func assert(t testing.TB, result bool, f func(), cd int) {
	if !result {
		_, file, line, _ := runtime.Caller(cd + 1)
		t.Errorf("%s:%d", file, line)
		f()
		t.FailNow()
	}
}

func equal(t testing.TB, exp, got interface{}, cd int, args ...interface{}) {
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

func tt(t testing.TB, result bool, cd int, args ...interface{}) {
	fn := func() {
		t.Errorf("!  Failure")
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, result, fn, cd+1)
}

func T(t testing.TB, result bool, args ...interface{}) {
	tt(t, result, 1, args...)
}

func Tf(t testing.TB, result bool, format string, args ...interface{}) {
	tt(t, result, 1, fmt.Sprintf(format, args...))
}

func IsNil(t testing.TB, ext bool, got interface{}, args ...interface{}) {
	fn := func() {
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, reflect.DeepEqual(nil, got) == ext, fn, 1)
}

func IsNilf(t testing.TB, ext bool, got interface{}, format string, args ...interface{}) {
	fn := func() {
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	assert(t, (got == nil) == ext, fn, 1)
}

func Equal(t testing.TB, exp, got interface{}, args ...interface{}) {
	equal(t, exp, got, 1, args...)
}

func Equalf(t testing.TB, exp, got interface{}, format string, args ...interface{}) {
	equal(t, exp, got, 1, fmt.Sprintf(format, args...))
}

func NotEqual(t testing.TB, exp, got interface{}, args ...interface{}) {
	fn := func() {
		t.Errorf("!  Unexpected: <%#v>", exp)
		if len(args) > 0 {
			t.Error("!", " -", fmt.Sprint(args...))
		}
	}
	result := !reflect.DeepEqual(exp, got)
	assert(t, result, fn, 1)
}
func NotEqualf(t testing.TB, exp, got interface{}, format string, args ...interface{}) {
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
func Panic(t testing.TB, expect interface{}, fn func(), args ...interface{}) {
	defer func() {
		equal(t, expect, recover(), 3, args...)
	}()
	fn()
}

func PanicError(t testing.TB, expect interface{}, fn func(), args ...interface{}) {
	defer func() {
		equal(t, expect, fmt.Sprintf("%s", recover()), 3, args...)
	}()
	fn()
}

var (
	wd, _ = os.Getwd()
)

func Expect(t testing.TB, a interface{}, b interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	fn = strings.Replace(fn, wd+"/", "", -1)

	if !reflect.DeepEqual(a, b) {
		t.Errorf("(%s:%d) Expected %v (type %v) - Got %v (type %v)", fn, line, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Refute(t testing.TB, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
