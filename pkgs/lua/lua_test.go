package lua

import (
	"testing"
)

func TestMainFuncMode(t *testing.T) {
	vm, err := New(`function main(a,b)
		return get('abc')..a..b;
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
	if v != "hello" {
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
		if h != "hello" {
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
		if h != "hello" {
			b.Error("err:", h)
			return
		}
	}
}
func get(s string) string {
	return "hello"
}
