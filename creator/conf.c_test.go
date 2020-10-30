package creator

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *conf
	}{
		{name: "初始化默认对象", want: &conf{
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
		{name: "初始化Loader默认对象", routerLoader: func(string) *services.ORouter {
			return &services.ORouter{}
		}, want: &conf{
			data: make(map[string]iCustomerBuilder),
			vars: make(map[string]map[string]interface{}),
			routerLoader: func(string) *services.ORouter {
				return &services.ORouter{}
			}}},
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
	tests := []struct {
		name        string
		serverTypes []string
		fields      *conf
		want        map[string]iCustomerBuilder
		wantErr     bool
	}{
		{name: "没有设置api节点,加载默认节点", serverTypes: []string{global.API}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, api.DefaultAPIAddress, cuurConfDefault.routerLoader)}, wantErr: true},
		{name: "没有设置WS节点,加载默认节点", serverTypes: []string{global.WS}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, api.DefaultWSAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},
		{name: "没有设置web节点,加载默认节点", serverTypes: []string{global.Web}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, api.DefaultWEBAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},
		{name: "没有设置RPC节点,加载默认节点", serverTypes: []string{global.RPC}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.RPC: newRPC(rpc.DefaultRPCAddress, cuurConfDefault.routerLoader)}, wantErr: true},
		{name: "没有设置MQC节点,加载默认节点", serverTypes: []string{global.MQC}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.101")}, wantErr: true},
		{name: "没有设置CRON节点,加载默认节点", serverTypes: []string{global.CRON}, fields: cuurConfDefault,
			want: map[string]iCustomerBuilder{global.CRON: newCron()}, wantErr: true},
		{name: "已经设置api节点,加载默认节点", serverTypes: []string{global.API}, fields: cuurConfAPI,
			want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, ":1122", cuurConfAPI.routerLoader)}, wantErr: true},
		{name: "已经设置WS节点,加载默认节点", serverTypes: []string{global.WS}, fields: cuurConfWS,
			want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, "1123", cuurConfWS.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},
		{name: "已经设置web节点,加载默认节点", serverTypes: []string{global.Web}, fields: cuurConfWEB,
			want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, ":1124", cuurConfWEB.routerLoader).Static(static.WithArchive(global.AppName))}, wantErr: true},
		{name: "已经设置RPC节点,加载默认节点", serverTypes: []string{global.RPC}, fields: cuurConfRPC,
			want: map[string]iCustomerBuilder{global.RPC: newRPC(":1125", cuurConfRPC.routerLoader)}, wantErr: true},
		{name: "已经设置MQC节点,加载默认节点", serverTypes: []string{global.MQC}, fields: cuurConfMQC,
			want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.102")}, wantErr: true},
		{name: "已经设置CRON节点,加载默认节点", serverTypes: []string{global.CRON}, fields: cuurConfCRON,
			want: map[string]iCustomerBuilder{global.CRON: newCron()}, wantErr: true},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = tt.serverTypes
		defer func(name string) {
			e := recover()
			if name == "没有设置MQC节点,加载默认节点" {
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
				assert.Equal(t, (tt.want[k].(*httpBuilder)).customerBuilder, (tt.fields.data[k].(*httpBuilder)).customerBuilder, tt.name+",http-customerBuilder")
				assert.Equal(t, reflect.TypeOf((tt.want[k].(*httpBuilder)).fnGetRouter), reflect.TypeOf((tt.fields.data[k].(*httpBuilder)).fnGetRouter), tt.name+",http-fnGetRouter")
			case global.RPC:
				assert.Equal(t, (tt.want[k].(*rpcBuilder)).tp, (tt.fields.data[k].(*rpcBuilder)).tp, tt.name+",rpc-tp")
				assert.Equal(t, (tt.want[k].(*rpcBuilder)).customerBuilder, (tt.fields.data[k].(*rpcBuilder)).customerBuilder, tt.name+",rpc-customerBuilder")
				assert.Equal(t, reflect.TypeOf((tt.want[k].(*rpcBuilder)).fnGetRouter), reflect.TypeOf((tt.fields.data[k].(*rpcBuilder)).fnGetRouter), tt.name+",rpc-fnGetRouter")
			case global.MQC:
				assert.Equal(t, tt.want[k], tt.fields.data[k], tt.name+",mqc-customerBuilder")
			case global.CRON:
				assert.Equal(t, tt.want[k], tt.fields.data[k], tt.name+",CRON-customerBuilder")
			}
			delete(tt.fields.data, k)
		}
		assert.Equal(t, 0, len(tt.fields.data), tt.name+",len")
		cuurConfDefault = New() //重置conf
	}
}

func Test_conf_API(t *testing.T) {

	tests := []struct {
		name    string
		address string
		opts    []api.Option
	}{
		{name: "设置默认对象", address: "", opts: []api.Option{}},
		{name: "设置自定义对象", address: ":9091", opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newHTTP(global.API, tt.address, cuurConf.routerLoader, tt.opts...)
		obj := cuurConf.API(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
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
		{name: "获取未设置api对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, api.DefaultAPIAddress, cuurConfDefault.routerLoader)}},
		{name: "获取已经设置对象", fields: cuurConfAPI, want: map[string]iCustomerBuilder{global.API: newHTTP(global.API, ":1122", cuurConfAPI.routerLoader)}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetAPI()
		assert.Equal(t, (tt.want[global.API].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.API].(*httpBuilder)).customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.API].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_Web(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []api.Option
	}{
		{name: "设置默认对象", address: "", opts: []api.Option{}},
		{name: "设置自定义对象", address: ":9092", opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newHTTP(global.Web, tt.address, cuurConf.routerLoader, tt.opts...)
		want.Static(static.WithArchive(global.AppName))
		obj := cuurConf.Web(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
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
		{name: "获取未设置web对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, api.DefaultWEBAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}},
		{name: "获取已经设置web对象", fields: cuurConfWeb, want: map[string]iCustomerBuilder{global.Web: newHTTP(global.Web, ":1122", cuurConfWeb.routerLoader).Static(static.WithArchive(global.AppName))}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetWeb()
		assert.Equal(t, (tt.want[global.Web].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.Web].(*httpBuilder)).customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.Web].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_WS(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []api.Option
	}{
		{name: "设置默认ws对象", address: "", opts: []api.Option{}},
		{name: "设置自定义ws对象", address: ":9092", opts: []api.Option{api.WithDisable(), api.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newHTTP(global.WS, tt.address, cuurConf.routerLoader, tt.opts...)
		want.Static(static.WithArchive(global.AppName))
		obj := cuurConf.WS(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
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
		{name: "获取未设置WS对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, api.DefaultWSAddress, cuurConfDefault.routerLoader).Static(static.WithArchive(global.AppName))}},
		{name: "获取已经设置ws对象", fields: cuurConfWs, want: map[string]iCustomerBuilder{global.WS: newHTTP(global.WS, ":1122", cuurConfWs.routerLoader).Static(static.WithArchive(global.AppName))}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetWS()
		assert.Equal(t, (tt.want[global.WS].(*httpBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.WS].(*httpBuilder)).customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.WS].(*httpBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_RPC(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []rpc.Option
	}{
		{name: "设置默认rpc对象", address: "", opts: []rpc.Option{}},
		{name: "设置自定义rpc对象", address: ":9092", opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newRPC(tt.address, cuurConf.routerLoader, tt.opts...)
		obj := cuurConf.RPC(tt.address, tt.opts...)
		assert.Equal(t, want.tp, obj.tp, tt.name+",tp")
		assert.Equal(t, want.customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
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
		{name: "获取未设置rpc对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.RPC: newRPC(rpc.DefaultRPCAddress, cuurConfDefault.routerLoader)}},
		{name: "获取已经设置rpc对象", fields: cuurConfRPC, want: map[string]iCustomerBuilder{global.RPC: newRPC(":1122", cuurConfRPC.routerLoader)}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetRPC()
		assert.Equal(t, (tt.want[global.RPC].(*rpcBuilder)).tp, obj.tp, tt.name+",tp")
		assert.Equal(t, (tt.want[global.RPC].(*rpcBuilder)).customerBuilder, obj.customerBuilder, tt.name+",customerBuilder")
		assert.Equal(t, reflect.TypeOf((tt.want[global.RPC].(*rpcBuilder)).fnGetRouter), reflect.TypeOf(obj.fnGetRouter), tt.name+",fnGetRouter")
	}
}

func Test_conf_CRON(t *testing.T) {
	tests := []struct {
		name string
		opts []cron.Option
	}{
		{name: "设置默认cron对象", opts: []cron.Option{}},
		{name: "设置自定义cron对象", opts: []cron.Option{cron.WithDisable(), cron.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newCron(tt.opts...)
		obj := cuurConf.CRON(tt.opts...)
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
		{name: "获取未设置rpc对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.CRON: newCron()}},
		{name: "获取已经设置rpc对象", fields: cuurConfCRON, want: map[string]iCustomerBuilder{global.CRON: newCron(cron.WithDisable())}},
	}

	for _, tt := range tests {
		obj := tt.fields.GetCRON()
		assert.Equal(t, tt.want[global.CRON], obj, tt.name)
	}
}

func Test_conf_MQC(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []mqc.Option
	}{
		{name: "设置默认mqc对象", address: "redis://192.168.0.11", opts: []mqc.Option{}},
		{name: "设置自定义mqc对象", address: "redis://192.168.0.12", opts: []mqc.Option{mqc.WithDisable(), mqc.WithTrace()}},
	}

	for _, tt := range tests {
		cuurConf := New()
		want := newMQC(tt.address, tt.opts...)
		obj := cuurConf.MQC(tt.address, tt.opts...)
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
		{name: "获取未设置mqc对象", fields: cuurConfDefault, want: map[string]iCustomerBuilder{global.MQC: newCron()}},
		{name: "获取已经设置mqc对象", fields: cuurConfMQC, want: map[string]iCustomerBuilder{global.MQC: newMQC("redis://192.168.0.102", mqc.WithDisable())}},
	}

	for _, tt := range tests {
		defer func(name string) {
			e := recover()
			if name == "获取未设置mqc对象" {
				assert.Equal(t, "未指定mqc服务器配置", types.GetString(e), name+",mqc-panic")
			} else {
				assert.Equal(t, true, e == nil, name+",panic1")
			}
		}(tt.name)
		obj := tt.fields.GetMQC()
		assert.Equal(t, tt.want[global.CRON], obj, tt.name)
	}
}

func Test_conf_Custome(t *testing.T) {
	cuurConfMQC := func() *conf {
		cuurConf := New()
		cuurConf.MQC("redis://192.168.0.102", mqc.WithDisable())
		return cuurConf
	}()
	tests := []struct {
		name   string
		key    string
		args   []string
		fields *conf
		want   map[string]iCustomerBuilder
	}{
		{name: "设置重复的自定义节点", key: global.MQC, args: []string{}, fields: cuurConfMQC, want: map[string]iCustomerBuilder{global.MQC: newCustomerBuilder()}},
		{name: "设置自定义节点数据", key: "testconf", args: []string{"redis://192.168.0.102", "sssd"}, fields: cuurConfMQC, want: map[string]iCustomerBuilder{"testconf": newCustomerBuilder("redis://192.168.0.102", "sssd")}},
	}

	for _, tt := range tests {
		defer func(name, key string) {
			e := recover()
			if name == "设置重复的自定义节点" {
				assert.Equal(t, "不能重复注册"+key, types.GetString(e), name+",Custome-panic")
			} else {
				assert.Equal(t, true, e == nil, name+",panic1")
			}
		}(tt.name, tt.key)
		obj := tt.fields.Custome(tt.key, tt.args)
		assert.Equal(t, tt.want[tt.key], obj, tt.name)
	}
}

func Test_conf_Encode(t *testing.T) {

	//空对象序列化为toml格式
	cuurConf := New()
	cuurConf.GetAPI()

	cuurConf.Encode()

	// sss := make(map[string]customerBuilder)
	// sss["api"] = customerBuilder{}
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

func Test_conf_Encode2File(t *testing.T) {

}

func Test_conf_Decode(t *testing.T) {

}
