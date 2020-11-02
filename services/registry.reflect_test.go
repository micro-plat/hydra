package services

import (
	"fmt"
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
		wantErr           bool
	}{
		{name: "path为空", h: hander1{}, wantErr: true},
		{name: "handler为空", path: "path", wantErr: true},
		{name: "handler非函数接收类型", path: "path", h: "xxxx", wantErr: true},
		{name: "handler为引用类型,但没有可用于注册的处理函数", path: "path", h: &testHandler1{}, wantErr: true},
		{name: "handler为func(context.IContext) interface{}", path: "path", h: func(context.IContext) interface{} { return nil },
			wantService:       []string{"path"},
			wantServicePath:   []string{"path"},
			wantServiceAction: [][]string{[]string{}}},
		{name: "handler为rpc协议", path: "path", h: "rpc://192.168.5.115:9091",
			wantService:       []string{"path"},
			wantServicePath:   []string{"path"},
			wantServiceAction: [][]string{[]string{}}},
		{name: "handler为引用类型,且有用于注册的处理函数", path: "path", h: &testHandler{},
			wantService:       []string{"/path/$get", "/path/$post", "/path/order"},
			wantServicePath:   []string{"path", "path", "/path/order"},
			wantServiceAction: [][]string{[]string{"GET"}, []string{"POST"}, []string{"GET", "POST"}},
		},
		{name: "handler为引用类型,且需要对注册服务进行替换", path: "/path/*/request", h: &testHandler5{},
			wantService:       []string{"/path/order/request", "/path/*/request/$post"},
			wantServicePath:   []string{"/path/order/request", "/path/*/request"},
			wantServiceAction: [][]string{[]string{"GET", "POST"}, []string{"POST"}},
		},
	}
	for _, tt := range tests {
		gotG, err := reflectHandle(tt.path, tt.h)
		fmt.Println(gotG)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		if reflect.ValueOf(tt.h).MethodByName("Close").IsValid() {
			assert.Equal(t, true, gotG.Closing != nil, tt.name)
		}

		if reflect.ValueOf(tt.h).MethodByName("Handling").IsValid() {
			assert.Equal(t, true, gotG.Handling != nil, tt.name+"2")
		}

		if reflect.ValueOf(tt.h).MethodByName("Handled").IsValid() {
			assert.Equal(t, true, gotG.Handled != nil, tt.name+"3")
		}
		for k, v := range tt.wantService {
			u, ok := gotG.Services[v]
			fmt.Println(tt.name, ok)
			assert.Equal(t, true, ok, tt.name)
			assert.Equal(t, v, u.Service, tt.name)
			assert.Equal(t, tt.wantServicePath[k], u.Path, tt.name)
			assert.Equal(t, tt.wantServiceAction[k], u.Actions, tt.name)
			//assert.Equal(t, tt.wantServiceHandler, u.Handle, tt.name) 地址无法比较
		}
	}
}
