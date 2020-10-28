package services

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func Test_serverServices_handleExt(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		f       func(u *Unit, ext ...interface{}) error
		g       *Unit
		ext     []interface{}
		wantErr bool
		errStr  string
	}{
		{name: "extHandle为空", f: nil, g: &Unit{}},
		{name: "extHandle不为空", f: func(u *Unit, ext ...interface{}) error {
			return fmt.Errorf("错误")
		}, g: &Unit{}, wantErr: true, errStr: "错误"},
	}
	for _, tt := range tests {
		err := newServerServices(tt.f).handleExt(tt.g, tt.ext...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
		}
	}
}

func Test_serverServices_Register(t *testing.T) {
	tests := []struct {
		name        string
		pName       string
		h           interface{}
		ext         []interface{}
		f           func(g *Unit, ext ...interface{}) error
		wantErr     bool
		errStr      string
		wantService []string
	}{
		// Register会panic报错 测试需要单条测试
		// {name: "注册对象为空", pName: "", wantErr: true, errStr: "注册对象不能为空"},
		// {name: "handleExt报错", pName: "name", h: &testHandler{},
		// 	f: func(g *Unit, ext ...interface{}) error { return fmt.Errorf("error") }, wantErr: true, errStr: "error"},
		// {name: "AddClosingHanle报错", pName: "name", h: "123456", f: nil, wantErr: true, errStr: "只能接收引用类型; 实际是 string"},
		{name: "注册正确", pName: "name", h: &testHandler{}, f: nil, wantService: []string{"/name/$get", "/name/$post", "/name/order"}},
	}

	for _, tt := range tests {
		s := newServerServices(tt.f)
		defer func() {
			e := recover()
			assert.Equal(t, tt.wantErr, e != nil, tt.name)
			if tt.wantErr {
				assert.Equal(t, tt.errStr, fmt.Sprintf("%s", e), tt.name)
			}
		}()
		s.Register(tt.pName, tt.h, tt.ext...)
		if !tt.wantErr {
			g, _ := reflectHandle(tt.pName, tt.h)
			for _, v := range tt.wantService {
				//检验handling
				assert.Equal(t, len(g.Services[v].GetHandlings()), len(s.handleHook.GetHandlings(v)), tt.name)

				//检验Handle
				_, ok := s.metaServices.GetHandlers(v)
				assert.Equal(t, true, ok, tt.name)
				//	assert.Equal(t, g.Services[v].Handle, handler, tt.name)地址无法比较

				//检验Handled
				assert.Equal(t, len(g.Services[v].GetHandleds()), len(s.handleHook.GetHandleds(v)), tt.name)

				//检验Fallback
				if g.Services[v].Fallback != nil {
					_, ok = s.metaServices.GetFallback(v)
					assert.Equal(t, true, ok, tt.name)
				}
			}
		}
	}
}
