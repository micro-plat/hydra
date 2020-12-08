package creator

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/types"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *conf
	}{
		{name: "1. 初始化默认对象", want: &conf{
			data:         make(map[string]iCustomerBuilder),
			vars:         make(map[string]map[string]interface{}),
			routerLoader: services.GetRouter}},
	}
	for _, tt := range tests {
		got := New()
		assert.Equal(t, tt.want.data, got.data, tt.name+",data")
		assert.Equal(t, tt.want.vars, got.vars, tt.name+",vars")
		assert.Equal(t, reflect.TypeOf(tt.want.routerLoader), reflect.TypeOf(got.routerLoader), tt.name+",routerLoader")
	}
}

func TestNewByLoader(t *testing.T) {
	tests := []struct {
		name         string
		routerLoader func(string) *services.ORouter
		want         *conf
	}{
		{name: "1. 初始化Loader默认对象", routerLoader: func(string) *services.ORouter { return &services.ORouter{} }, want: &conf{
			data:         make(map[string]iCustomerBuilder),
			vars:         make(map[string]map[string]interface{}),
			routerLoader: func(string) *services.ORouter { return &services.ORouter{} }}},
	}
	for _, tt := range tests {
		got := NewByLoader(tt.routerLoader)
		assert.Equal(t, tt.want.data, got.data, tt.name+",data")
		assert.Equal(t, tt.want.vars, got.vars, tt.name+",vars")
		assert.Equal(t, reflect.TypeOf(tt.want.routerLoader), reflect.TypeOf(got.routerLoader), tt.name+",routerLoader")
	}
}

func Test_conf_Load(t *testing.T) {
	servers.Register(global.API, nil)
	servers.Register(global.WS, nil)
	servers.Register(global.Web, nil)
	servers.Register(global.RPC, nil)
	servers.Register(global.MQC, nil)
	servers.Register(global.CRON, nil)
	cuurConfDefault := New()
	cuurConfAPI := func() *conf { cuurConf := New(); cuurConf.API(":1122"); return cuurConf }()
	cuurConfWS := func() *conf { cuurConf := New(); cuurConf.WS(":1123"); return cuurConf }()
	cuurConfWEB := func() *conf { cuurConf := New(); cuurConf.Web(":1124"); return cuurConf }()
	cuurConfMQC := func() *conf { cuurConf := New(); cuurConf.MQC("redis://192.168.0.102"); return cuurConf }()
	cuurConfRPC := func() *conf { cuurConf := New(); cuurConf.RPC(":1125"); return cuurConf }()
	cuurConfCRON := func() *conf { cuurConf := New(); cuurConf.CRON(); return cuurConf }()
	cuurConfCustom := func() *conf { cuurConf := New(); cuurConf.Custom("test", "自定义配置"); return cuurConf }()
	cuurConfAll := func() *conf {
		cuurConf := New()
		cuurConf.Custom("test", "自定义配置")
		cuurConf.API(":1122")
		cuurConf.WS(":1123")
		cuurConf.MQC("redis://192.168.0.102")
		cuurConf.RPC(":1125")
		cuurConf.CRON()
		return cuurConf
	}()
	tests := []struct {
		name        string
		serverTypes []string
		fields      *conf
		want        map[string]iCustomerBuilder
		wantErr     bool
	}{
		{name: "1.1 没有设置api节点,加载默认节点", serverTypes: []string{global.API}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, api.DefaultAPIAddress, cuurConfDefault.routerLoader)}, wantErr: true},
		{name: "1.2 已经设置api节点,加载默认节点", serverTypes: []string{global.API}, fields: cuurConfAPI, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, ":1122", cuurConfAPI.routerLoader)}, wantErr: true},

		{name: "2.1 没有设置WS节点,加载默认节点", serverTypes: []string{global.WS}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, api.DefaultWSAddress, cuurConfDefault.routerLoader)}, wantErr: true},
		{name: "2.2 已经设置WS节点,加载默认节点", serverTypes: []string{global.WS}, fields: cuurConfWS, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, ":1123", cuurConfWS.routerLoader)}, wantErr: true},

		{name: "3.1 没有设置web节点,加载默认节点", serverTypes: []string{global.Web}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, api.DefaultWEBAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},
		{name: "3.2 已经设置web节点,加载默认节点", serverTypes: []string{global.Web}, fields: cuurConfWEB, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, ":1124", cuurConfWEB.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},

		{name: "4.1 没有设置RPC节点,加载默认节点", serverTypes: []string{global.RPC}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.RPC: newRPC(rpc.DefaultRPCAddress, cuurConfDefault.routerLoader)}, wantErr: true},
		{name: "4.2 已经设置RPC节点,加载默认节点", serverTypes: []string{global.RPC}, fields: cuurConfRPC, want: map[string]iCustomerBuilder{global.RPC: newRPC(":1125", cuurConfRPC.routerLoader)}, wantErr: true},

		{name: "5.1 没有设置MQC节点,加载默认节点", serverTypes: []string{global.MQC}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.101")}, wantErr: true},
		{name: "5.2 已经设置MQC节点,加载默认节点", serverTypes: []string{global.MQC}, fields: cuurConfMQC, want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.102")}, wantErr: true},

		{name: "6.1 没有设置CRON节点,加载默认节点", serverTypes: []string{global.CRON}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.CRON: newCron()}, wantErr: true},
		{name: "6.2 已经设置CRON节点,加载默认节点", serverTypes: []string{global.CRON}, fields: cuurConfCRON, want: map[string]iCustomerBuilder{global.CRON: newCron()}, wantErr: true},

		{name: "7.1 没有设置任何节点", serverTypes: []string{}, fields: cuurConfDefault, want: map[string]iCustomerBuilder{}, wantErr: true},
		{name: "7.2 已经设置自定义节点节点-test", serverTypes: []string{"test"}, fields: cuurConfCustom, want: map[string]iCustomerBuilder{"test": newCustomerBuilder("自定义配置")}, wantErr: true},
		{name: "7.3 同时加载所有服务节点", serverTypes: []string{global.API, global.WS, global.Web, global.RPC, global.CRON, "test"}, fields: cuurConfAll,
			want: map[string]iCustomerBuilder{"test": newCustomerBuilder("自定义配置"), global.API: newHTTP(global.API, ":1122", cuurConfAPI.routerLoader), global.WS: newHTTP(global.WS, ":1123", cuurConfWS.routerLoader),
				global.Web: newHTTP(global.Web, ":1124", cuurConfWEB.routerLoader).Static(static.WithArchive(global.AppName)), global.RPC: newRPC(":1125", cuurConfRPC.routerLoader),
				global.MQC: newMQC("redis://192.168.0.102"), global.CRON: newCron()}, wantErr: true},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = tt.serverTypes
		defer func(name string) {
			e := recover()
			if name == "5.1 没有设置MQC节点,加载默认节点" {
				assert.Equal(t, "未指定mqc服务器配置", types.GetString(e), name+",mqc-panic")
			} else {
				assert.Equal(t, true, e == nil, name+",panic1")
			}
		}(tt.name)
		err := tt.fields.Load()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		for k, _ := range tt.want {
			tt.want[k].Load()
			_, ok := tt.fields.data[k]
			assert.Equal(t, true, ok, tt.name+",ok")
			switch k {
			case global.API, global.WS, global.Web:
				assert.Equal(t, (tt.want[k].(*httpBuilder)).tp, (tt.fields.data[k].(*httpBuilder)).tp, tt.name+",http-tp")
				assert.Equal(t, (tt.want[k].(*httpBuilder)).CustomerBuilder, (tt.fields.data[k].(*httpBuilder)).CustomerBuilder, tt.name+",http-CustomerBuilder")
				assert.Equal(t, reflect.TypeOf((tt.want[k].(*httpBuilder)).fnGetRouter), reflect.TypeOf((tt.fields.data[k].(*httpBuilder)).fnGetRouter), tt.name+",http-fnGetRouter")
			case global.RPC:
				assert.Equal(t, (tt.want[k].(*rpcBuilder)).tp, (tt.fields.data[k].(*rpcBuilder)).tp, tt.name+",rpc-tp")
				assert.Equal(t, (tt.want[k].(*rpcBuilder)).CustomerBuilder, (tt.fields.data[k].(*rpcBuilder)).CustomerBuilder, tt.name+",rpc-CustomerBuilder")
				assert.Equal(t, reflect.TypeOf((tt.want[k].(*rpcBuilder)).fnGetRouter), reflect.TypeOf((tt.fields.data[k].(*rpcBuilder)).fnGetRouter), tt.name+",rpc-fnGetRouter")
			case global.MQC:
				assert.Equal(t, tt.want[k], tt.fields.data[k], tt.name+",mqc-CustomerBuilder")
			case global.CRON:
				assert.Equal(t, tt.want[k], tt.fields.data[k], tt.name+",CRON-CustomerBuilder")
			default:
				assert.Equal(t, tt.want[k], tt.fields.data[k], tt.name+",Custom-CustomerBuilder")
			}
			delete(tt.fields.data, k)
		}
		assert.Equal(t, 0, len(tt.fields.data), tt.name+",len")
		cuurConfDefault = New() //重置conf
	}
}

func Test_conf_API(t *testing.T) {
	cuurConfDefault := New()
	cuurConfAPI := func() *conf { cuurConf := New(); cuurConf.API(":1122"); return cuurConf }()
	tests := []struct {
		name    string
		address string
		fields  *conf
		opts    []api.Option
	}{
		{name: "1. 设置默认api配置对象", address: "", fields: cuurConfDefault, opts: []api.Option{}},
		{name: "2. 设置自定义api配置对象", address: ":9091", fields: cuurConfDefault, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
		{name: "3. 重复设置自定义api配置对象", address: ":9091", fields: cuurConfAPI, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newHTTP(global.API, tt.address, cuurConf.routerLoader, tt.opts...)
		obj := cuurConf.API(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf(want.fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_GetAPI(t *testing.T) {
	cuurConfDefault := New()
	cuurConfAPI := func() *conf { cuurConf := New(); cuurConf.API(":1122"); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取api配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, api.DefaultAPIAddress, cuurConfDefault.routerLoader)}},
		{name: "2. 已设置,获取api配置对象", fields: cuurConfAPI, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, ":1122", cuurConfAPI.routerLoader)}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetAPI()
		assert.Equal(t, (tt.want[global.API].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.API].(*httpBuilder)).CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.API].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_Web(t *testing.T) {
	cuurConfDefault := New()
	cuurConfWeb := func() *conf { cuurConf := New(); cuurConf.Web(":1122"); return cuurConf }()
	tests := []struct {
		name    string
		address string
		fields  *conf
		opts    []api.Option
	}{
		{name: "1. 设置默认web配置对象", address: "", fields: cuurConfDefault, opts: []api.Option{}},
		{name: "2. 设置自定义web配置对象", address: ":9092", fields: cuurConfDefault, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
		{name: "3. 重复设置自定义web配置对象", address: ":9092", fields: cuurConfWeb, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		want := newHTTP(global.Web, tt.address, tt.fields.routerLoader, tt.opts...)
		want.Static(static.WithArchive(global.AppName))
		obj := tt.fields.Web(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf(want.fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_GetWeb(t *testing.T) {
	cuurConfDefault := New()
	cuurConfWeb := func() *conf { cuurConf := New(); cuurConf.Web(":1122"); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取web配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, api.DefaultWEBAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}},
		{name: "2. 已设置,获取web配置对象", fields: cuurConfWeb, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, ":1122", cuurConfWeb.routerLoader).Static(static.WithArchive(global.AppName))}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetWeb()
		assert.Equal(t, (tt.want[global.Web].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.Web].(*httpBuilder)).CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.Web].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_WS(t *testing.T) {
	cuurConfDefault := New()
	cuurConfWs := func() *conf { cuurConf := New(); cuurConf.WS(":1122"); return cuurConf }()
	tests := []struct {
		name    string
		address string
		fields  *conf
		opts    []api.Option
	}{
		{name: "1. 设置默认ws配置对象", address: "", fields: cuurConfDefault, opts: []api.Option{}},
		{name: "2. 设置自定义ws配置对象", address: ":9092", fields: cuurConfDefault, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
		{name: "3. 重复设置自定义ws配置对象", address: ":9092", fields: cuurConfWs, opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		want := newHTTP(global.WS, tt.address, tt.fields.routerLoader, tt.opts...)
		obj := tt.fields.WS(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf(want.fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_GetWS(t *testing.T) {
	cuurConfDefault := New()
	cuurConfWs := func() *conf { cuurConf := New(); cuurConf.WS(":1122"); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取WS配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, api.DefaultWSAddress, cuurConfDefault.routerLoader)}},
		{name: "2. 已设置,获取ws配置对象", fields: cuurConfWs, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, ":1122", cuurConfWs.routerLoader)}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetWS()
		assert.Equal(t, (tt.want[global.WS].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.WS].(*httpBuilder)).CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.WS].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_RPC(t *testing.T) {
	cuurConfDefault := New()
	cuurConfRPC := func() *conf { cuurConf := New(); cuurConf.RPC(":1122"); return cuurConf }()
	tests := []struct {
		name    string
		address string
		fields  *conf
		opts    []rpc.Option
	}{
		{name: "1. 设置默认rpc配置对象", address: "", fields: cuurConfDefault, opts: []rpc.Option{}},
		{name: "2. 设置自定义rpc配置对象", address: ":9092", fields: cuurConfDefault, opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
		{name: "3. 重复设置自定义rpc配置对象", address: ":9092", fields: cuurConfRPC, opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
	}

	for _, tt := range tests {
		want := newRPC(tt.address, tt.fields.routerLoader, tt.opts...)
		obj := tt.fields.RPC(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf(want.fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_GetRPC(t *testing.T) {
	cuurConfDefault := New()
	cuurConfRPC := func() *conf { cuurConf := New(); cuurConf.RPC(":1122"); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取rpc配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.RPC: newRPC(rpc.DefaultRPCAddress, cuurConfDefault.routerLoader)}},
		{name: "2. 已设置,获取rpc配置对象", fields: cuurConfRPC, want: map[string]iCustomerBuilder{global.RPC: newRPC(":1122", cuurConfRPC.routerLoader)}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetRPC()
		assert.Equal(t, (tt.want[global.RPC].(*rpcBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.RPC].(*rpcBuilder)).CustomerBuilder, obj.CustomerBuilder, tt.name+",CustomerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.RPC].(*rpcBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_CRON(t *testing.T) {
	cuurConfDefault := New()
	cuurConfCRON := func() *conf { cuurConf := New(); cuurConf.CRON(cron.WithDisable()); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		opts   []cron.Option
	}{
		{name: "1. 设置默认cron配置对象", fields: cuurConfDefault, opts: []cron.Option{}},
		{name: "2. 设置自定义cron配置对象", fields: cuurConfDefault, opts: []cron.Option{cron.WithDisable(), cron.WithTrace()}},
		{name: "3. 重复设置cron配置对象", fields: cuurConfCRON, opts: []cron.Option{cron.WithDisable(), cron.WithTrace()}},
	}

	for _, tt := range tests {
		want := newCron(tt.opts...)
		obj := tt.fields.CRON(tt.opts...)
		assert.Equal(t, want, obj, tt.name)
	}
}

func Test_conf_GetCRON(t *testing.T) {
	cuurConfDefault := New()
	cuurConfCRON := func() *conf { cuurConf := New(); cuurConf.CRON(cron.WithDisable()); return cuurConf }()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取rpc配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.CRON: newCron()}},
		{name: "2. 已设置,获取rpc配置对象", fields: cuurConfCRON, want: map[string]iCustomerBuilder{global.CRON: newCron(cron.WithDisable())}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetCRON()
		assert.Equal(t, tt.want[global.CRON], obj, tt.name)
	}
}

func Test_conf_MQC(t *testing.T) {
	cuurConfDefault := New()
	cuurConfMQC := func() *conf {
		cuurConf := New()
		cuurConf.MQC("redis://192.168.0.102", mqc.WithDisable())
		return cuurConf
	}()
	tests := []struct {
		name    string
		address string
		fields  *conf
		opts    []mqc.Option
	}{
		{name: "1. 设置默认mqc配置对象", address: "redis://192.168.0.11", fields: cuurConfDefault, opts: []mqc.Option{}},
		{name: "2. 设置自定义mqc配置对象", address: "redis://192.168.0.12", fields: cuurConfDefault, opts: []mqc.Option{mqc.WithDisable(), mqc.WithTrace()}},
		{name: "3. 重复设置mqc配置对象", address: "redis://192.168.0.12", fields: cuurConfMQC, opts: []mqc.Option{mqc.WithDisable(), mqc.WithTrace()}},
	}

	for _, tt := range tests {
		want := newMQC(tt.address, tt.opts...)
		obj := tt.fields.MQC(tt.address, tt.opts...)
		assert.Equal(t, want, obj, tt.name)
	}
}

func Test_conf_GetMQC(t *testing.T) {
	cuurConfDefault := New()
	cuurConfMQC := func() *conf {
		cuurConf := New()
		cuurConf.MQC("redis://192.168.0.102", mqc.WithDisable())
		return cuurConf
	}()
	tests := []struct {
		name   string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 未设置,获取mqc配置对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.MQC: newCron()}},
		{name: "2. 已设置,获取mqc配置对象", fields: cuurConfMQC, want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.102", mqc.WithDisable())}},
	}

	for _, tt := range tests {
		defer func(name string) {
			e := recover()
			if name == "1. 未设置,获取mqc配置对象" {
				assert.Equal(t, "未指定mqc服务器配置", types.GetString(e), name+",mqc-panic")
			} else {
				assert.Equal(t, true, e == nil, name+",panic1")
			}
		}(tt.name)
		obj := tt.fields.GetMQC()
		assert.Equal(t, tt.want[global.CRON], obj, tt.name)
	}
}

func Test_conf_GetVar(t *testing.T) {
	tests := []struct {
		name     string
		tp       string
		tpname   string
		v        map[string]map[string]interface{}
		wantBool bool
		wantStr  interface{}
	}{
		{name: "1. 初始化数据为空-获取空key", tp: "", tpname: "", v: map[string]map[string]interface{}{}, wantBool: false, wantStr: nil},
		{name: "2. 初始化数据不未空-获取空key", tp: "", tpname: "", v: map[string]map[string]interface{}{"test": map[string]interface{}{"test1": "123"}}, wantBool: false, wantStr: nil},
		{name: "3. 初始化数据不未空-tp不存在", tp: "test1", tpname: "test1", v: map[string]map[string]interface{}{"test": map[string]interface{}{"test1": "123"}}, wantBool: false, wantStr: nil},
		{name: "4. 初始化数据不未空-tp存在-tpname不存在", tp: "test", tpname: "test", v: map[string]map[string]interface{}{"test": map[string]interface{}{"test1": "123"}}, wantBool: false, wantStr: nil},
		{name: "5. 初始化数据不未空-tp存在-tpname存在", tp: "test", tpname: "test1", v: map[string]map[string]interface{}{"test": map[string]interface{}{"test1": "123"}}, wantBool: true, wantStr: "123"},
	}
	for _, tt := range tests {
		conf := New()
		conf.vars = tt.v
		gotVal, gotOk := conf.GetVar(tt.tp, tt.tpname)
		assert.Equal(t, tt.wantStr, gotVal, tt.name)
		assert.Equal(t, tt.wantBool, gotOk, tt.name)
	}
}

func Test_conf_Custom(t *testing.T) {
	cuurConfC := func() *conf {
		cuurConf := New()
		cuurConf.Custom("redis", "测试")
		return cuurConf
	}()
	tests := []struct {
		name   string
		key    string
		args   []string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "1. 设置重复的自定义节点", key: "redis", args: []string{"测试"}, fields: cuurConfC, want: map[string]iCustomerBuilder{global.MQC: newCustomerBuilder()}},
		{name: "2. 设置自定义节点数据", key: "testconf", args: []string{"redis://192.168.0.102", "sssd"}, fields: cuurConfC, want: map[string]iCustomerBuilder{"testconf": newCustomerBuilder("redis://192.168.0.102", "sssd")}},
	}

	for _, tt := range tests {
		defer func(name, key string) {
			e := recover()
			if name == "1. 设置重复的自定义节点" {
				assert.Equal(t, "不能重复注册"+key, types.GetString(e), name+",Custome-panic")
			} else {
				assert.Equal(t, true, e == nil, name+",panic1")
			}
		}(tt.name, tt.key)
		obj := tt.fields.Custom(tt.key, tt.args)
		assert.Equal(t, tt.want[tt.key], obj, tt.name)
	}
}

//toml 文件格式化  暂时不用测试
func Test_conf_Encode(t *testing.T) {

	// //空对象序列化为toml格式
	// cuurConf := New()
	// cuurConf.GetAPI()

	// sss, err := cuurConf.Encode()
	// fmt.Println("sss:", sss)
	// fmt.Println("err:", err)
	// sss := make(map[string]CustomerBuilder)
	// sss["api"] = CustomerBuilder{}
	// sss["api"]["xxx"] = api.Server{Address: "ssssssss", Status: "start",
	// 	RTimeout:  10,
	// 	WTimeout:  10,
	// 	RHTimeout: 10,
	// 	Host:      "www.baidu.com",
	// 	Domain:    "sssss",
	// 	Trace:     true}

	// var buffer bytes.Buffer
	// encoder := toml.NewEncoder(&buffer)
	// err := encoder.Encode(sss)
	// // cuurConf.GetAPI()
	// // ss, err := cuurConf.Encode()
	// fmt.Println("sss:", buffer.String())
	// fmt.Println("err:", err)

	// f, err := os.OpenFile("./fs", os.O_RDWR|os.O_CREATE, 0755)
	// if err != nil {
	// 	fmt.Println("无法打开文件:%w", err)
	// 	return
	// }
	// encoder1 := toml.NewEncoder(f)
	// err = encoder1.Encode(sss)
	// if err != nil {
	// 	fmt.Println("11无法打开文件:%w", err)
	// 	return
	// }
	// if err := f.Close(); err != nil {
	// 	fmt.Println("22222:%w", err)
	// 	return
	// }
	// nodes := make(map[string]string)
	// vnodes := make(map[string]map[string]interface{})
	// if _, err := toml.DecodeFile("./fs", &vnodes); err != nil {
	// 	fmt.Println("11111:%w", err)
	// 	return
	// }
	// platName := "platName1"
	// systemName := "systemName1"
	// clusterName := "clusterName1"
	// for k, sub := range vnodes {
	// 	for name, value := range sub {
	// 		var path = registry.Join(platName, systemName, k, clusterName, "conf", name)
	// 		if name == "main" {
	// 			path = registry.Join(platName, systemName, k, clusterName, "conf")
	// 		}
	// 		buff, err := json.Marshal(&value)
	// 		if err != nil {
	// 			fmt.Println("2222222:%w", err)
	// 			return
	// 		}
	// 		nodes[path] = string(buff)
	// 	}
	// }

	// xxx, _ := json.Marshal(nodes)

	// fmt.Println("nodes:", string(xxx))
}

//toml 文件格式化  暂时不用测试
func Test_conf_Encode2File(t *testing.T) {

}

//toml 文件格式化  暂时不用测试
func Test_conf_Decode(t *testing.T) {

}
