package services

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
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
		{name: "1.1 注册路由为空", h: hander1{}, wantErr: "注册对象不能为空"},
		{name: "1.2 注册对象为空", path: "path", h: nil, wantErr: "注册对象不能为空"},
		{name: "1.3 注册对象为string", path: "path", h: "xxxx", wantErr: "只能接收引用类型或struct; 实际是 string"},
		{name: "1.4 注册对象为int", path: "path", h: 0, wantErr: "只能接收引用类型或struct; 实际是 int"},
		{name: "1.5 注册对象为float", path: "path", h: 0.1, wantErr: "只能接收引用类型或struct; 实际是 float64"},
		{name: "1.6 注册对象为map", path: "path", h: map[string]string{}, wantErr: "只能接收引用类型或struct; 实际是 map"},
		{name: "1.7 注册对象为结构体指针,但没有可用于注册的处理函数", path: "path", h: &testHandler1{}, wantErr: "path中，未找到可用于注册的处理函数"},
		{name: "1.8 注册对象为结构体指针,无Handle函数", path: "path", h: &testHandler6{}, wantErr: "path中,未指定[/path/order]的Handle函数"},
		{name: "1.9 注册对象为结构体指针,函数签名不正确", path: "path", h: &testHandlerSuffix{}, wantErr: "函数【XxxHandle】是钩子类型（[Handling Handle Handled Fallback]）,但签名不是func(context.IContext) interface{}或者func(context.IContext)"},
		{name: "1.10 注册对象为构建函数,函数签名不正确", path: "/path/*/request", h: func() int32 { return 0 }, wantErr: "输出参数第一个参数必须是结构体"},

		{name: "2.1 注册对象为func(context.IContext) interface{}", path: "path", h: func(context.IContext) interface{} { return nil }, wantService: []string{"path"}, wantServicePath: []string{"path"}, wantServiceAction: [][]string{[]string{}}},
		{name: "2.2 注册对象为rpc协议", path: "path", h: "rpc://192.168.0.1:9091", wantService: []string{"rpc://192.168.0.1:9091"}, wantServicePath: []string{"path"}, wantServiceAction: [][]string{[]string{}}},
		{name: "2.3 注册对象为结构体指针,且有用于注册的处理函数", path: "/path", h: &testHandler{}, wantService: []string{"/path/$get", "/path/$post", "/path/order"}, wantServicePath: []string{"/path", "/path", "/path/order"}, wantServiceAction: [][]string{[]string{"GET"}, []string{"POST"}, []string{"GET", "POST"}}},
		{name: "2.4 注册对象为结构体指针,且需要对注册服务进行替换", path: "/path/*/request", h: &testHandler5{}, wantService: []string{"/path/order/request", "/path/post/request/$post"}, wantServicePath: []string{"/path/order/request", "/path/post/request"}, wantServiceAction: [][]string{[]string{"GET", "POST"}, []string{"POST"}}},
		{name: "2.5 注册对象为结构体指针,且需要对服务进行默认替换", path: "/path/*", h: &testHandler4{}, wantService: []string{"/path/handle"}, wantServicePath: []string{"/path/handle"}, wantServiceAction: [][]string{[]string{}}},
		{name: "2.6 注册对象为正确的对象构建方法", path: "/path", h: newTestHandler, wantService: []string{"/path/$get", "/path/$post", "/path/order"}, wantServicePath: []string{"/path", "/path", "/path/order"}, wantServiceAction: [][]string{[]string{"GET"}, []string{"POST"}, []string{"GET", "POST"}}},
		{name: "2.7 注册对象为正确的对象", path: "/path", h: testHandler7{}, wantService: []string{"/path", "/path/$post", "/path/order"}, wantServicePath: []string{"/path", "/path", "/path/order"}, wantServiceAction: [][]string{[]string{}, []string{"POST"}, []string{"GET", "POST"}}},
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
		{name: "1.1 不包含后缀", args: args{mName: "xxx"}, want: false},
		{name: "1.2 包含Handling后缀", args: args{mName: "xxxHandling"}, want: true},
		{name: "1.3 包含Handle后缀", args: args{mName: "xxxHandle"}, want: true},
		{name: "1.4 包含Handled后缀", args: args{mName: "xxxHandled"}, want: true},
		{name: "1.5 包含Fallback后缀", args: args{mName: "xxxFallback"}, want: true},
	}
	for _, tt := range tests {
		got := checkSuffix(tt.args.mName)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
