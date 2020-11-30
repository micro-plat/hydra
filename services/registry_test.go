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
		{name: "1. 获取api对应的serverServices", tp: global.API, want: Def.servers[global.API]},
		{name: "2. 获取web对应的serverServices", tp: global.Web, want: Def.servers[global.Web]},
		{name: "3. 获取rpc对应的serverServices", tp: global.RPC, want: Def.servers[global.RPC]},
		{name: "4. 获取ws对应的serverServices", tp: global.WS, want: Def.servers[global.WS]},
		{name: "5. 获取cron对应的serverServices", tp: global.CRON, want: Def.servers[global.CRON]},
		{name: "6. 获取mqc对应的serverServices", tp: global.MQC, want: Def.servers[global.MQC]},
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
	f := func(g *Unit, ext ...interface{}) error {
		return nil
	}
	tests := []struct {
		name string
		tp   string
		f    []func(g *Unit, ext ...interface{}) error
	}{
		{name: "1.注册file服务类型", tp: "file", f: nil},
		{name: "2.注册socket服务类型", tp: "socket", f: []func(g *Unit, ext ...interface{}) error{f}},
	}
	s := Def
	for _, tt := range tests {
		s.RegisterServer(tt.tp, tt.f...)
	}

	//注册已经存在的服务类型
	assert.Panic(t, errors.New("服务api已存在，不能重复注册"), func() { s.RegisterServer(global.API) }, "注册已经存在的服务类型")
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
		{name: "1.指定api,添加启动预处理函数", tps: []string{global.API}, h: h},
		{name: "2.指定web,添加启动预处理函数", tps: []string{global.Web}, h: h},
		{name: "3.指定ws,添加启动预处理函数", tps: []string{global.WS}, h: h},
		{name: "4.指定cron,添加启动预处理函数", tps: []string{global.CRON}, h: h},
		{name: "5.指定mqc,添加启动预处理函数", tps: []string{global.MQC}, h: h},
		{name: "6.未指定服务类型,添加启动预处理函数", tps: []string{}, h: h},
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
		{name: "1.指定api,添加关闭服务", tps: []string{global.API}, h: h},
		{name: "2.指定web,添加关闭服务", tps: []string{global.Web}, h: h},
		{name: "3.指定ws,添加关闭服务", tps: []string{global.WS}, h: h},
		{name: "4.指定cron,添加关闭服务", tps: []string{global.CRON}, h: h},
		{name: "5.指定mqc,添加关闭服务", tps: []string{global.MQC}, h: h},
		{name: "6.未指定服务类型,添加关闭服务", tps: []string{}, h: h},
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
		{name: "1.指定api,添加空的业务预处理钩子", h: nil, tps: []string{global.API}},
		{name: "2.指定api,添加业务预处理钩子", tps: []string{global.API}, h: h},
		{name: "3.指定web,添加空的业务预处理钩子", h: nil, tps: []string{global.Web}},
		{name: "4.指定web,添加业务预处理钩子", tps: []string{global.Web}, h: h},
		{name: "5.指定ws,添加空的业务预处理钩子", h: nil, tps: []string{global.WS}},
		{name: "6.指定ws,添加业务预处理钩子", tps: []string{global.WS}, h: h},
		{name: "7.指定cron,添加空的业务预处理钩子", h: nil, tps: []string{global.CRON}},
		{name: "8.指定cron,添加业务预处理钩子", tps: []string{global.CRON}, h: h},
		{name: "9.指定mqc,添加空的业务预处理钩子", h: nil, tps: []string{global.MQC}},
		{name: "10.指定mqc,添加业务预处理钩子", tps: []string{global.MQC}, h: h},
		{name: "11.未指定服务类型,添加空的业务预处理钩子", h: nil, tps: []string{}},
		{name: "12.未指定服务类型,添加业务预处理钩子", tps: []string{}, h: h},
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
		{name: "1.指定api,添加空的业务后处理钩子", h: nil, tps: []string{global.API}},
		{name: "2.指定api,添加业务后处理钩子", tps: []string{global.API}, h: h},
		{name: "3.指定web,添加空的业务后处理钩子", h: nil, tps: []string{global.Web}},
		{name: "4.指定web,添加业务后处理钩子", tps: []string{global.Web}, h: h},
		{name: "5.指定ws,添加空的业务后处理钩子", h: nil, tps: []string{global.WS}},
		{name: "6.指定ws,添加业务后处理钩子", tps: []string{global.WS}, h: h},
		{name: "7.指定cron,添加空的业务后处理钩子", h: nil, tps: []string{global.CRON}},
		{name: "8.指定cron,添加业务后处理钩子", tps: []string{global.CRON}, h: h},
		{name: "9.指定mqc,添加空的业务后处理钩子", h: nil, tps: []string{global.MQC}},
		{name: "10.指定mqc,添加业务后处理钩子", tps: []string{global.MQC}, h: h},
		{name: "11.未指定服务类型,添加空的业务后处理钩子", h: nil, tps: []string{}},
		{name: "12.未指定服务类型,添加业务后处理钩子", tps: []string{}, h: h},
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
		{name: "1.注册api类型服务", tp: global.API, path: "/path", h: &testHandler{}},
		{name: "2.注册cron类型服务", tp: global.CRON, path: "/path1", h: &testHandler{}, ext: []interface{}{"taks1", "task2"}},
		{name: "3.注册web类型服务", tp: global.Web, path: "/path2", h: &testHandler{}},
		{name: "4.注册rpc类型服务", tp: global.RPC, path: "/path3", h: &testHandler{}},
		{name: "5.注册ws类型服务", tp: global.WS, path: "/path4", h: &testHandler{}},
		{name: "6.注册mqc类型服务", tp: global.MQC, path: "/path5", h: &testHandler{}, ext: []interface{}{"queue1", "queue2"}},
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
		assert.Equal(t, len(ext), len(MQC.queues.Queues), testName+" mqc")
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
		tp      string
		h       interface{}
		wantErr bool
		errStr  string
	}{
		{name: "1.1.api注册的Handler的Close()不存在", tp: global.API, path: "/api1", h: &hander1{}, wantErr: false},
		{name: "1.2.api注册的Handler的Close()未报错", tp: global.API, path: "/api2", h: &testHandler2{}, wantErr: false},
		{name: "1.3.api注册的Handler的Close()报错", tp: global.API, path: "/api4", h: &testHandler4{}, wantErr: true, errStr: "error"},
	}
	s := Def
	for _, tt := range tests {
		s.Custom(tt.tp, tt.path, tt.h)
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
		h    []interface{}
		ext  []router.Option
	}{
		{name: "1.api服务注册对象为结构体指针", path: "/api/reg/path1", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "2.api服务注册对象为结构体", path: "/api/reg/path2", h: []interface{}{testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "3.api服务注册对象为构建函数", path: "/api/reg/path3", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "4.api服务注册对象为函数", path: "/api/reg/path4", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "5.api服务同一个注册对象,注册两个不同地址", path: "/api/reg/path5", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "6.api服务注册对象为两个struct", path: "/api/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "7.api服务注册对象为一个struct,一个构建函数", path: "/api/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		for _, v := range tt.h {
			s.API(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.API, tt.name, tt.ext)
	}
}

func Test_regist_Web(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    []interface{}
		ext  []router.Option
	}{
		{name: "1.web服务注册对象为结构体指针", path: "/web/reg/path1", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "2.web服务注册对象为结构体", path: "/web/reg/path2", h: []interface{}{testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "3.web服务注册对象为构建函数", path: "/web/reg/path3", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "4.web服务注册对象为函数", path: "/web/reg/path4", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "5.web服务同一个注册对象,注册两个不同地址", path: "/web/reg/path5", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "6.web服务注册对象为两个struct", path: "/web/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "7.web服务注册对象为一个struct,一个构建函数", path: "/web/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		for _, v := range tt.h {
			s.Web(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.Web, tt.name, tt.ext)
	}
}

func Test_regist_RPC(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    []interface{}
		ext  []router.Option
	}{
		{name: "1.rpc服务注册对象为结构体指针", path: "/rpc/reg/path1", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "2.rpc服务注册对象为结构体", path: "/rpc/reg/path2", h: []interface{}{testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "3.rpc服务注册对象为构建函数", path: "/rpc/reg/path3", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "4.rpc服务注册对象为函数", path: "/rpc/reg/path4", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "5.rpc服务同一个注册对象,注册两个不同地址", path: "/rpc/reg/path5", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "6.rpc服务注册对象为两个struct", path: "/rpc/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "7.rpc服务注册对象为一个struct,一个构建函数", path: "/rpc/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		for _, v := range tt.h {
			s.RPC(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.RPC, tt.name, tt.ext)
	}
}

func Test_regist_WS(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    []interface{}
		ext  []router.Option
	}{
		{name: "1.ws服务注册对象为结构体指针", path: "/ws/reg/path1", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "2.ws服务注册对象为结构体", path: "/ws/reg/path2", h: []interface{}{testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "3.ws服务注册对象为构建函数", path: "/ws/reg/path3", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "4.ws服务注册对象为函数", path: "/ws/reg/path4", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "5.ws服务同一个注册对象,注册两个不同地址", path: "/ws/reg/path5", h: []interface{}{&testHandler{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "6.ws服务注册对象为两个struct", path: "/ws/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
		{name: "7.ws服务注册对象为一个struct,一个构建函数", path: "/ws/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []router.Option{router.WithPages("pages"), router.WithEncoding("utf-8")}},
	}
	s := Def
	for _, tt := range tests {
		for _, v := range tt.h {
			s.WS(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.WS, tt.name, tt.ext)
	}
}

func Test_regist_MQC(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    []interface{}
		ext  []string
	}{
		{name: "1.mqc服务注册对象为结构体指针", path: "/mqc/reg/path1", h: []interface{}{&testHandler{}}, ext: []string{"queue1", "queue2"}},
		{name: "2.mqc服务注册对象为结构体", path: "/mqc/reg/path2", h: []interface{}{testHandler8{}}, ext: []string{"queue1", "queue2"}},
		{name: "3.mqc服务注册对象为构建函数", path: "/mqc/reg/path3", h: []interface{}{&testHandler{}}, ext: []string{"queue1", "queue2"}},
		{name: "4.mqc服务注册对象为函数", path: "/mqc/reg/path4", h: []interface{}{&testHandler{}}, ext: []string{"queue1", "queue2"}},
		{name: "5.mqc服务同一个注册对象,注册两个不同地址", path: "/mqc/reg/path5", h: []interface{}{&testHandler{}}, ext: []string{"queue1", "queue2"}},
		{name: "6.mqc服务注册对象为两个struct", path: "/mqc/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []string{"queue1", "queue2"}},
		{name: "7.mqc服务注册对象为一个struct,一个构建函数", path: "/mqc/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []string{"queue1", "queue2"}},
	}
	s := Def
	global.MQConf.PlatNameAsPrefix(false)
	for _, tt := range tests {
		for _, v := range tt.h {
			s.MQC(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.MQC, tt.name, tt.ext[0], tt.ext[1])
	}
}

func Test_regist_CRON(t *testing.T) {
	tests := []struct {
		name string
		path string
		h    []interface{}
		ext  []string
	}{
		{name: "1.cron服务注册对象为结构体指针", path: "/cron/reg/path1", h: []interface{}{&testHandler{}}, ext: []string{"task1", "task2"}},
		{name: "2.cron服务注册对象为结构体", path: "/cron/reg/path2", h: []interface{}{testHandler8{}}, ext: []string{"task1", "task2"}},
		{name: "3.cron服务注册对象为构建函数", path: "/cron/reg/path3", h: []interface{}{&testHandler{}}, ext: []string{"task1", "task2"}},
		{name: "4.cron服务注册对象为函数", path: "/cron/reg/path4", h: []interface{}{&testHandler{}}, ext: []string{"task1", "task2"}},
		{name: "5.cron服务同一个注册对象,注册两个不同地址", path: "/cron/reg/path5", h: []interface{}{&testHandler{}}, ext: []string{"task1", "task2"}},
		{name: "6.cron服务注册对象为两个struct", path: "/cron/reg/path6", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []string{"task1", "task2"}},
		{name: "7.cron服务注册对象为一个struct,一个构建函数", path: "/cron/reg/path7", h: []interface{}{testHandler7{}, testHandler8{}}, ext: []string{"task1", "task2"}},
	}
	s := Def
	for _, tt := range tests {
		for _, v := range tt.h {
			s.CRON(tt.path, v, tt.ext...)
		}
		checkTestCustomeResult(t, s, global.CRON, tt.name, tt.ext[0], tt.ext[1])
	}
}
