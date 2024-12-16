package global

import (
	"errors"
	"sync"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

var lock sync.Mutex

func TestOnReady_Normal(t *testing.T) {
	type args struct {
		fs []interface{}
	}
	tests := []struct {
		name    string
		isReady bool
		funcs   []func() error
		args    args
		wantLen int
	}{
		{name: "1. 参数类型为func()-isReady=false,noerror", isReady: false, funcs: nil, args: args{fs: []interface{}{func() {}}}, wantLen: 1},
		{name: "2. 参数类型为func()-isReady=true,noerror", isReady: true, funcs: nil, args: args{fs: []interface{}{func() {}}}, wantLen: 0},
		{name: "3. 参数类型为func()error-isReady=false,noerror", isReady: false, funcs: nil, args: args{fs: []interface{}{func() error { return nil }}}, wantLen: 1},
		{name: "4. 参数类型为func()error-isReady=true,noerror", isReady: true, funcs: nil, args: args{fs: []interface{}{func() error { return nil }}}, wantLen: 0},
		{name: "5. 参数类型为func()error-isReady=false,error", isReady: false, funcs: nil, args: args{fs: []interface{}{func() error { return errors.New("error") }}}, wantLen: 1},
	}
	for _, tt := range tests {
		lock.Lock()
		isReady = tt.isReady
		funcs = tt.funcs
		OnReady(tt.args.fs...)
		assert.Equal(t, tt.wantLen, len(funcs), tt.name)
		lock.Unlock()
	}
}

func TestOnReady_Panic(t *testing.T) {
	type args struct {
		fs []interface{}
	}
	tests := []struct {
		name      string
		isReady   bool
		funcs     []func() error
		args      args
		wantPanic interface{}
	}{
		{name: "1. 参数类型为func()error-isReady=true,error", isReady: true, funcs: nil, args: args{fs: []interface{}{func() error { return errors.New("error") }}}, wantPanic: errors.New("error")},
		{name: "2. 参数类型为不匹配", isReady: true, funcs: nil, args: args{fs: []interface{}{func(x int) error { return nil }}}, wantPanic: "函数签名格式不正确，支持的格式有func(){} 或 func()error{}"},
	}
	for _, tt := range tests {
		assert.Panics(t, func() {
			isReady = tt.isReady
			funcs = tt.funcs
			OnReady(tt.args.fs...)
		}, tt.name)
	}
}
