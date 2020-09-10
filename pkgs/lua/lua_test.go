package lua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLib(t *testing.T) {
	vm, err := New(`
	require "math"
	return math.abs(-300)`)

	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.Call()
	if err != nil {
		t.Error(err)
	}
	if v != "300" {
		t.Error("结果不一致:", v)
	}
}

func TestModule(t *testing.T) {
	vm, err := New(`r=require("request");
	return r.get("abc")`, WithModule("request", map[string]lua.LGFunction{
		"get": func(ls *lua.LState) int {
			s := lua.LVAsString(ls.Get(1))
			ls.Push(lua.LString(s))
			return 1
		},
	}))

	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.Call()
	if err != nil {
		t.Error(err)
	}
	if v != "abc" {
		t.Error("结果不一致:", v)
	}
}
func TestModules(t *testing.T) {
	vm, err := New(`r=require("request");
	return r.get("abc")`, WithModules(map[string]map[string]lua.LGFunction{
		"request": map[string]lua.LGFunction{
			"get": func(ls *lua.LState) int {
				s := lua.LVAsString(ls.Get(1))
				ls.Push(lua.LString(s))
				return 1
			},
		},
	}))

	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.Call()
	if err != nil {
		t.Error(err)
	}
	if v != "abc" {
		t.Error("结果不一致:", v)
	}
}
func TestUserType(t *testing.T) {
	vm, err := New(`local r=req();
	return r:Get("abc")`, WithType("req", hello{}))

	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.Call()
	if err != nil {
		t.Error(err)
	}
	if v != "abc" {
		t.Error("结果不一致:", v)
	}
}

func TestMainFuncMode(t *testing.T) {
	vm, err := New(`function main(a,b)
		return get('hello')..a..b;
	end`, With("get", get), WithMainFuncMode())

	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.CallByMethod("main", " colin ", 10)
	if err != nil {
		t.Error(err)
	}
	if v[0] != "hello colin 10" {
		t.Error("结果不一致:", v)
	}
}

func TestCodeBlockMode(t *testing.T) {
	vm, err := New(`return get('abc')`, With("get", get))
	if err != nil {
		t.Error(err)
		return
	}
	v, err := vm.Call()
	if err != nil {
		t.Error(err)
	}
	if v != "abc" {
		t.Error("结果不一致:", v)
	}
}
func BenchmarkMainFuncMode(b *testing.B) {
	vm, err := New(`function main()
		return get('abc');
	end`, With("get", get), WithMainFuncMode())

	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		h, err := vm.Call()
		if err != nil {
			b.Error(err)
			return
		}
		if h != "abc" {
			b.Error("err:", h)
			return
		}
	}
}
func BenchmarkCodeBlockMode(b *testing.B) {
	vm, err := New(`return get('abc')`, With("get", get))
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		h, err := vm.Call()
		if err != nil {
			b.Error(err)
			return
		}
		if h != "abc" {
			b.Error("err:", h)
			return
		}
	}
}
func get(s string) string {
	return s
}

type hello struct {
}

func (h *hello) Get(s string) string {
	return s
}
