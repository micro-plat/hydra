package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_regist_get(t *testing.T) {
	tests := []struct {
		name string
		tp   string
		want *serverServices
	}{
		{name: "获取初始化的api对应的services", tp: global.API, want: Def.servers[global.API]},
		{name: "获取初始化的web对应的services", tp: global.Web, want: Def.servers[global.Web]},
		{name: "获取初始化的rpc对应的services", tp: global.RPC, want: Def.servers[global.RPC]},
		{name: "获取初始化的ws对应的services", tp: global.WS, want: Def.servers[global.WS]},
		{name: "获取初始化的cron对应的services", tp: global.CRON, want: Def.servers[global.CRON]},
		{name: "获取初始化的mqc对应的services", tp: global.MQC, want: Def.servers[global.MQC]},
	}
	s := Def
	for _, tt := range tests {
		got := s.get(tt.tp)
		assert.Equal(t, tt.want, got, tt.name)
	}

	//获取不支持的类型
	assert.Panic(t, "不支持的服务器类型:xxx", func() { s.get("xxx") }, "获取不支持的类型")
}

func Test_regist_RegisterServer(t *testing.T) {
	tests := []struct {
		name string
		tp   string
		f    []func(g *Unit, ext ...interface{}) error
	}{
		{name: "注册未初始化的file服务器", tp: "file", f: nil},
		{name: "注册未初始化的socket服务器", tp: "socket", f: []func(g *Unit, ext ...interface{}) error{
			func(g *Unit, ext ...interface{}) error {
				return nil
			}}},
	}
	s := Def
	for _, tt := range tests {
		s.RegisterServer(tt.tp, tt.f...)
	}

	//注册已经存在的服务器
	assert.Panic(t, errors.New("服务api已存在，不能重复注册"), func() { s.RegisterServer(global.API) }, "注册已经存在的服务器")
}

func Test_regist_OnStarting(t *testing.T) {
	h := func(app.IAPPConf) error {
		return fmt.Errorf("test")
	}
	tests := []struct {
		name string
		h    func(app.IAPPConf) error
		tps  []string
	}{
		{name: "tps不为空", h: h, tps: []string{global.API}},
		{name: "tps为空", h: h, tps: []string{}},
	}
	s := Def
	for _, tt := range tests {
		s.OnStarting(tt.h, tt.tps...)
		tps := tt.tps
		if len(tps) == 0 {
			tps = global.Def.ServerTypes
		}
		for _, v := range tps {
			gotErr := s.get(v).DoStarting(nil)
			assert.Equal(t, "test", gotErr.Error(), "启动添加的服务")
		}
	}

	//h为空 panic
	assert.Panic(t, errors.New("api OnServerStarting 启动服务不能为空"), func() { s.OnStarting(nil, []string{global.API}...) }, "handle为空")
}

func Test_regist_OnClosing(t *testing.T) {
	h := func(app.IAPPConf) error {
		return fmt.Errorf("test")
	}
	tests := []struct {
		name string
		h    func(app.IAPPConf) error
		tps  []string
	}{
		{name: "tps不为空", h: h, tps: []string{global.API}},
		{name: "tps为空", h: h, tps: []string{}},
	}
	s := Def
	for _, tt := range tests {
		s.OnClosing(tt.h, tt.tps...)
		tps := tt.tps
		if len(tps) == 0 {
			tps = global.Def.ServerTypes
		}
		for _, v := range tps {
			gotErr := s.get(v).DoClosing(nil)
			assert.Equal(t, "test", gotErr.Error(), "启动添加的关闭服务")
		}
	}

	//h为空 panic
	assert.Panic(t, errors.New("api OnServerClosing 关闭服务不能为空"), func() { s.OnClosing(nil, []string{global.API}...) }, "handle为空")
}

func Test_regist_OnHandleExecuting(t *testing.T) {
	h := func(c context.IContext) interface{} {
		return nil
	}

	tests := []struct {
		name string
		h    context.Handler
		tps  []string
	}{
		{name: "添加空接口", h: h, tps: []string{global.API}},
		{name: "添加api单个接口", tps: []string{global.API}, h: h},
		{name: "添加api单个接口,servertypes为空", tps: []string{}, h: h},
		{name: "添加cron单个接口", tps: []string{global.CRON}, h: h},
	}
	s := Def
	for _, tt := range tests {
		s.OnHandleExecuting(tt.h, tt.tps...)
		tps := tt.tps
		if len(tps) == 0 {
			tps = global.Def.ServerTypes
		}
		for _, v := range tps {
			gotHandle := s.GetHandleExecutings(v)
			assert.Equal(t, s.get(v).serverHook.handlings, gotHandle, tt.name)
		}
	}
}

func Test_regist_OnHandleExecuted(t *testing.T) {
	h := func(c context.IContext) interface{} {
		return nil
	}

	tests := []struct {
		name string
		h    context.Handler
		tps  []string
	}{
		{name: "添加空接口", h: h, tps: []string{global.API}},
		{name: "添加api单个接口", tps: []string{global.API}, h: h},
		{name: "添加api单个接口,servertypes为空", tps: []string{}, h: h},
		{name: "添加cron单个接口", tps: []string{global.CRON}, h: h},
	}
	s := Def
	for _, tt := range tests {
		s.OnHandleExecuted(tt.h, tt.tps...)
		tps := tt.tps
		if len(tps) == 0 {
			tps = global.Def.ServerTypes
		}
		for _, v := range tps {
			gotHandle := s.GetHandleExecuted(v)
			assert.Equal(t, s.get(v).serverHook.handleds, gotHandle, tt.name)
		}
	}
}

func Test_regist_Custome(t *testing.T) {
	tests := []struct {
		name string
		tp   string
		path string
		h    interface{}
		ext  []interface{}
	}{
		{name: "注册api类型", tp: global.API, path: "/path", h: &testHandler{}},
		{name: "注册cron类型", tp: global.CRON, path: "/path1", h: &testHandler{}, ext: []interface{}{"taks1", "task2"}},
		{name: "注册web类型", tp: global.Web, path: "/path2", h: &testHandler{}},
		{name: "注册rpc类型", tp: global.RPC, path: "/path3", h: &testHandler{}},
		{name: "注册ws类型", tp: global.WS, path: "/path4", h: &testHandler{}},
		{name: "注册mqc类型", tp: global.MQC, path: "/path5", h: &testHandler{}, ext: []interface{}{"queue1", "queue2"}},
	}
	s := Def
	global.MQConf.PlatNameAsPrefix(false)
	for _, tt := range tests {
		s.Custom(tt.tp, tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, tt.tp, tt.name, tt.ext...)
	}
}

func checkTestCustomeResult(t *testing.T, s *regist, tp, testName string, ext ...interface{}) {

	m := s.get(tp).metaServices

	//cron
	if tp == global.CRON {
		assert.Equal(t, len(ext)*len(m.services), len(CRON.tasks.Tasks), testName+" cron")
		return
	}

	//mqc
	if tp == global.MQC {
		assert.Equal(t, len(ext)*len(m.services), len(MQC.queues.Queues), testName+" mqc")
		return
	}

	//api,web,ws,rpc
	h := s.get(tp).handleHook
	o := GetRouter(tp)

	routers, err := o.GetRouters()
	assert.Equal(t, false, err != nil, testName+"getrouters")

	//检查routers

	for _, r := range routers.Routers {
		exist := false
		for _, v := range m.services {
			if r.Service == v {
				exist = true
			}
		}
		assert.Equal(t, true, exist, testName+r.Service+" not exists")
	}

	for _, v := range m.services {

		//检查fallback
		gotFallback, ok := s.GetFallback(tp, v)
		if ok {
			assert.Equal(t, fmt.Sprintf("%+v", m.fallbacks[v]), fmt.Sprintf("%+v", gotFallback), testName+"fallback")
		}

		//检查handler
		gotHandler, ok := s.GetHandler(tp, v)
		if ok {
			assert.Equal(t, fmt.Sprintf("%+v", m.handlers[v]), fmt.Sprintf("%+v", gotHandler), testName+"handler")
		}

		//检查handling
		if handling, ok := h.handlings[v]; ok {
			gotHandling := s.GetHandlings(tp, v)
			assert.Equal(t, handling, gotHandling, testName+"handling")
		}

		//检查handled
		if handled, ok := h.handleds[v]; ok {
			gotHandled := s.GetHandleds(tp, v)
			assert.Equal(t, handled, gotHandled, testName+"handled")
		}
	}
}

func Test_regist_Close(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		h       interface{}
		wantErr bool
		errStr  string
	}{
		{name: "注册的Handler的Close()不存在", path: "/api1", h: &hander1{}, wantErr: false},
		{name: "注册的Handler的Close()未报错", path: "/api2", h: &testHandler2{}, wantErr: false},
		{name: "注册的Handler的Close()报错", path: "/api4", h: &testHandler4{}, wantErr: true, errStr: "error"},
	}
	s := Def
	for _, tt := range tests {
		s.Custom(global.API, tt.path, tt.h)
		err := s.Close()
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
		}
	}
}

func Test_regist_API(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []router.Option
	}{
		{name: "注册api服务", path: "/path6", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "同一个handler,注册api服务", path: "/path7", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		s.API(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.API, tt.name, tt.ext)
	}
}

func Test_regist_Web(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []router.Option
	}{
		{name: "注册web服务", path: "/path8", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "同一个handler,注册web服务", path: "/path9", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		s.Web(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.Web, tt.name, tt.ext)
	}
}

func Test_regist_RPC(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []router.Option
	}{
		{name: "注册rpc服务", path: "/path10", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "同一个handler,注册rpc服务", path: "/path11", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		s.RPC(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.RPC, tt.name, tt.ext)
	}
}

func Test_regist_WS(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []router.Option
	}{
		{name: "注册ws服务", path: "/path12", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "同一个handler,注册ws服务", path: "/path13", h: &testHandler{}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		s.WS(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.WS, tt.name, tt.ext)
	}
}

func Test_regist_MQC(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []string
	}{
		{name: "注册mqc服务", path: "/path14", h: &testHandler{}, ext: []string{"queue1", "queue2"}},
		{name: "同一个handler,注册mqc服务", path: "/path15", h: &testHandler{}, ext: []string{"queue1", "queue2"}},
	}
	s := Def
	global.MQConf.PlatNameAsPrefix(false)
	for _, tt := range tests {
		s.MQC(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.MQC, tt.name, tt.ext[0], tt.ext[1])
	}
}

func Test_regist_CRON(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    interface{}
		ext  []string
	}{
		{name: "注册cron服务", path: "/path16", h: &testHandler{}, ext: []string{"task1", "task2"}},
		{name: "同一个handler,注册cron服务", path: "/path17", h: &testHandler{}, ext: []string{"task1", "task2"}},
	}
	s := Def
	for _, tt := range tests {
		s.CRON(tt.path, tt.h, tt.ext...)
		checkTestCustomeResult(t, s, global.CRON, tt.name, tt.ext[0], tt.ext[1])
	}
}
