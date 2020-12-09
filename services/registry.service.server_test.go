package services

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/assert"
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
		{name: "1. extHandle为空", f: nil, g: &Unit{}},
		{name: "2. extHandle报错", f: func(u *Unit, ext ...interface{}) error { return fmt.Errorf("错误") }, g: &Unit{}, wantErr: true, errStr: "错误"},
	}
	for _, tt := range tests {
		err := newServerServices(tt.f).handleExt(tt.g, tt.ext...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
		}
	}
}

func Test_serverServices_Register_WithPanic(t *testing.T) {
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
		{name: "1.注册对象为空", pName: "", wantErr: true, errStr: "注册对象不能为空"},
		{name: "2.handleExt报错", pName: "name", h: &testHandler{}, f: func(g *Unit, ext ...interface{}) error { return fmt.Errorf("error") }, wantErr: true, errStr: "error"},
		{name: "3.AddClosingHanle报错", pName: "name", h: "123456", f: nil, wantErr: true, errStr: "只能接收引用类型或struct; 实际是 string"},
	}

	for _, tt := range tests {
		assert.Panics(t, func() {
			//errors.New(tt.errStr),
			s := newServerServices(tt.f)
			s.Register(tt.pName, tt.h, tt.ext...)
		}, tt.name)
	}
}

func Test_serverServices_Register(t *testing.T) {
	var API = NewORouter()
	var WS = NewORouter()
	var WEB = NewORouter()
	var RPC = NewORouter()
	var CRON = newCron()
	var MQC = newMQC()
	cronFunc := func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			CRON.Add(t.(string), g.Service)
		}
		return nil
	}
	global.MQConf.PlatNameAsPrefix(false)

	mqcFunc := func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			MQC.Add(t.(string), g.Service)
		}
		return nil
	}
	//注册正确
	tests := []struct {
		name        string
		pName       string
		h           interface{}
		ext         []interface{}
		f           func(g *Unit, ext ...interface{}) error
		wantService []string
	}{
		{name: "1.1.api注册对象为结构体指针", pName: "/api/path1", h: &testHandler{}, f: func(g *Unit, ext ...interface{}) error { return API.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/api/path1/$get", "/api/path1/$post", "/api/path1/order"}},
		{name: "1.2.api注册对象为结构体", pName: "/api/path2", h: testHandler7{}, f: func(g *Unit, ext ...interface{}) error { return API.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/api/path2", "/api/path2/$post", "/api/path2/order"}},
		{name: "1.3.api注册对象为func(context.IContext) interface{}", pName: "/api/path3", h: func(context.IContext) interface{} { return nil }, f: func(g *Unit, ext ...interface{}) error { return API.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/api/path3"}},
		{name: "1.4.api注册对象为对象构建方法", pName: "/api/path4", h: newTestHandler, f: func(g *Unit, ext ...interface{}) error { return API.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/api/path4/$get", "/api/path4/$post", "/api/path4/order"}},

		{name: "2.1.web注册对象为结构体指针", pName: "/web/path1", h: &testHandler{}, f: func(g *Unit, ext ...interface{}) error { return WEB.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/web/path1/$get", "/web/path1/$post", "/web/path1/order"}},
		{name: "2.2.web注册对象为结构体", pName: "/web/path2", h: testHandler7{}, f: func(g *Unit, ext ...interface{}) error { return WEB.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/web/path2", "/web/path2/$post", "/web/path2/order"}},
		{name: "2.3.web注册对象为func(context.IContext) interface{}", pName: "/web/path3", h: func(context.IContext) interface{} { return nil }, f: func(g *Unit, ext ...interface{}) error { return WEB.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/web/path3"}},
		{name: "2.4.web注册对象为对象构建方法", pName: "/web/path4", h: newTestHandler, f: func(g *Unit, ext ...interface{}) error { return WEB.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/web/path4/$get", "/web/path4/$post", "/web/path4/order"}},

		{name: "3.1.ws注册对象为结构体指针", pName: "/ws/path1", h: &testHandler{}, f: func(g *Unit, ext ...interface{}) error { return WS.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/ws/path1/$get", "/ws/path1/$post", "/ws/path1/order"}},
		{name: "3.2.ws注册对象为结构体", pName: "/ws/path2", h: testHandler7{}, f: func(g *Unit, ext ...interface{}) error { return WS.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/ws/path2", "/ws/path2/$post", "/ws/path2/order"}},
		{name: "3.3.ws注册对象为func(context.IContext) interface{}", pName: "/ws/path3", h: func(context.IContext) interface{} { return nil }, f: func(g *Unit, ext ...interface{}) error { return WS.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/ws/path3"}},
		{name: "3.4.ws注册对象为对象构建方法", pName: "/ws/path4", h: newTestHandler, f: func(g *Unit, ext ...interface{}) error { return WS.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/ws/path4/$get", "/ws/path4/$post", "/ws/path4/order"}},

		{name: "4.1.rpc注册对象为结构体指针", pName: "/rpc/path1", h: &testHandler{}, f: func(g *Unit, ext ...interface{}) error { return RPC.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/rpc/path1/$get", "/rpc/path1/$post", "/rpc/path1/order"}},
		{name: "4.2.rpc注册对象为结构体", pName: "/rpc/path2", h: testHandler7{}, f: func(g *Unit, ext ...interface{}) error { return RPC.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/rpc/path2", "/rpc/path2/$post", "/rpc/path2/order"}},
		{name: "4.3.rpc注册对象为func(context.IContext) interface{}", pName: "/rpc/path3", h: func(context.IContext) interface{} { return nil }, f: func(g *Unit, ext ...interface{}) error { return RPC.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/rpc/path3"}},
		{name: "4.4.rpc注册对象为对象构建方法", pName: "/rpc/path4", h: newTestHandler, f: func(g *Unit, ext ...interface{}) error { return RPC.Add(g.Path, g.Service, g.Actions, ext...) }, wantService: []string{"/rpc/path4/$get", "/rpc/path4/$post", "/rpc/path4/order"}},

		{name: "5.1.cron注册对象为结构体指针", pName: "/cron/path1", ext: []interface{}{"queue"}, h: &testHandler{}, f: cronFunc, wantService: []string{"/cron/path1/$get", "/cron/path1/$post", "/cron/path1/order"}},
		{name: "5.2.cron注册对象为结构体", pName: "/cron/path2", ext: []interface{}{"queue"}, h: testHandler7{}, f: cronFunc, wantService: []string{"/cron/path2", "/cron/path2/$post", "/cron/path2/order"}},
		{name: "5.3.cron注册对象为func(context.IContext) interface{}", pName: "/cron/path3", ext: []interface{}{"queue"}, h: func(context.IContext) interface{} { return nil }, f: cronFunc, wantService: []string{"/cron/path3"}},
		{name: "5.4.cron注册对象为对象构建方法", pName: "/cron/path4", ext: []interface{}{"queue"}, h: newTestHandler, f: cronFunc, wantService: []string{"/cron/path4/$get", "/cron/path4/$post", "/cron/path4/order"}},

		{name: "6.1.mqc注册对象为结构体指针", pName: "/mqc/path1", ext: []interface{}{"mqc1"}, h: &testHandler{}, f: mqcFunc, wantService: []string{"/mqc/path1/$get", "/mqc/path1/$post", "/mqc/path1/order"}},
		{name: "6.2.mqc注册对象为结构体", pName: "/mqc/path2", ext: []interface{}{"mqc2"}, h: testHandler7{}, f: mqcFunc, wantService: []string{"/mqc/path2", "/mqc/path2/$post", "/mqc/path2/order"}},
		{name: "6.3.mqc注册对象为func(context.IContext) interface{}", pName: "/mqc/path3", ext: []interface{}{"mqc3"}, h: func(context.IContext) interface{} { return nil }, f: mqcFunc, wantService: []string{"/mqc/path3"}},
		{name: "6.4.mqc注册对象为对象构建方法", pName: "/mqc/path4", ext: []interface{}{"mqc4"}, h: newTestHandler, f: mqcFunc, wantService: []string{"/mqc/path4/$get", "/mqc/path4/$post", "/mqc/path4/order"}},
	}

	for _, tt := range tests {

		s := newServerServices(tt.f)
		s.Register(tt.pName, tt.h, tt.ext...)
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
	r, _ := API.GetRouters()
	assert.Equal(t, 10, len(r.Routers), "api的services数量")
	r, _ = WEB.GetRouters()
	assert.Equal(t, 10, len(r.Routers), "web的services数量")
	r, _ = WS.GetRouters()
	assert.Equal(t, 10, len(r.Routers), "ws的services数量")
	r, _ = RPC.GetRouters()
	assert.Equal(t, 10, len(r.Routers), "rpc的services数量")
	assert.Equal(t, 10, len(CRON.dynamicTasks.Tasks), "cron数量")
	assert.Equal(t, 4, len(MQC.dynamicQueues.Queues), "MQC数量")
}
