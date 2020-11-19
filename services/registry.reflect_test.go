package services

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/test/assert"
)

// &Unit{Group: g, Path: path, Service: service, Actions: actions, Handle: h}
func Test_reflectHandle(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		h                 interface{}
		wantService       []string
		wantServicePath   []string
		wantServiceAction [][]string
		wantErr           string
	}{
		{name: "path为空", h: hander1{}, wantErr: "注册对象不能为空"},
		{name: "handler为空", path: "path", wantErr: "注册对象不能为空"},
		{name: "handler非函数接收类型", path: "path", h: "xxxx", wantErr: "只能接收引用类型或struct; 实际是 string"},
		{name: "handler为引用类型,但没有可用于注册的处理函数", path: "path", h: &testHandler1{}, wantErr: "path中，未找到可用于注册的处理函数"},
		{name: "handler为引用类型,无Handle函数", path: "path", h: &testHandler6{}, wantErr: "path中,未指定[/path/order]的Handle函数"},
		{name: "handler含有Suffix但签名不匹配", path: "path", h: &testHandlerSuffix{},
			wantErr: "函数【XxxHandle】是钩子类型（[Handling Handle Handled Fallback]）,但签名不是func(context.IContext) interface{}"},
		{name: "handler为错误的对象构建方法", path: "/path/*/request", h: func() int32 { return 0 }, wantErr: "输出参数第一个参数必须是结构体"},
		{name: "handler为func(context.IContext) interface{}", path: "path", h: func(context.IContext) interface{} { return nil },
			wantService:       []string{"path"},
			wantServicePath:   []string{"path"},
			wantServiceAction: [][]string{[]string{}}},
		{name: "handler为rpc协议", path: "path", h: "rpc://192.168.0.1:9091",
			wantService:       []string{"path"},
			wantServicePath:   []string{"path"},
			wantServiceAction: [][]string{[]string{}}},
		{name: "handler为引用类型,且有用于注册的处理函数", path: "/path", h: &testHandler{},
			wantService:       []string{"/path/$get", "/path/$post", "/path/order"},
			wantServicePath:   []string{"/path", "/path", "/path/order"},
			wantServiceAction: [][]string{[]string{"GET"}, []string{"POST"}, []string{"GET", "POST"}},
		},
		{name: "handler为引用类型,且需要对注册服务进行替换", path: "/path/*/request", h: &testHandler5{},
			wantService:       []string{"/path/order/request", "/path/post/request/$post"},
			wantServicePath:   []string{"/path/order/request", "/path/post/request"},
			wantServiceAction: [][]string{[]string{"GET", "POST"}, []string{"POST"}},
		},
		{name: "handler为引用类型,且需要对服务进行默认替换", path: "/path/*", h: &testHandler4{},
			wantService:       []string{"/path/handle"},
			wantServicePath:   []string{"/path/handle"},
			wantServiceAction: [][]string{[]string{}},
		},
		{name: "handler为正确的对象构建方法", path: "/path", h: newTestHandler,
			wantService:       []string{"/path/$get", "/path/$post", "/path/order"},
			wantServicePath:   []string{"/path", "/path", "/path/order"},
			wantServiceAction: [][]string{[]string{"GET"}, []string{"POST"}, []string{"GET", "POST"}},
		},
		{name: "handler为正确的对象", path: "/path", h: testHandler7{},
			wantService:       []string{"/path", "/path/$post", "/path/order"},
			wantServicePath:   []string{"/path", "/path", "/path/order"},
			wantServiceAction: [][]string{[]string{}, []string{"POST"}, []string{"GET", "POST"}},
		},
	}
	for _, tt := range tests {
		gotG, err := reflectHandle(tt.path, tt.h)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		if reflect.ValueOf(tt.h).MethodByName("Close").IsValid() {
			assert.Equal(t, true, gotG.Closing != nil, tt.name+"Close")
		}

		if reflect.ValueOf(tt.h).MethodByName("Handling").IsValid() {
			assert.Equal(t, true, gotG.Handling != nil, tt.name+"Handling")
		}

		if reflect.ValueOf(tt.h).MethodByName("Handled").IsValid() {
			assert.Equal(t, true, gotG.Handled != nil, tt.name+"Handled")
		}
		for k, v := range tt.wantService {
			u, ok := gotG.Services[v]
			assert.Equal(t, true, ok, tt.name+"Services")
			assert.Equal(t, v, u.Service, tt.name)
			assert.Equal(t, tt.wantServicePath[k], u.Path, tt.name)
			assert.Equal(t, tt.wantServiceAction[k], u.Actions, tt.name)
			//assert.Equal(t, tt.wantServiceHandler, u.Handle, tt.name) 地址无法比较
		}
	}
}

func Test_checkSuffix(t *testing.T) {
	type args struct {
		mName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "不包含", args: args{mName: "xxx"}, want: false},
		{name: "包含Handling", args: args{mName: "xxxHandling"}, want: true},
		{name: "包含Handle", args: args{mName: "xxxHandle"}, want: true},
		{name: "包含Handled", args: args{mName: "xxxHandled"}, want: true},
		{name: "包含Fallback", args: args{mName: "xxxFallback"}, want: true},
	}
	for _, tt := range tests {
		got := checkSuffix(tt.args.mName)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
