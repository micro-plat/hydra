package creator

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/types"
)

func Test_newRPC(t *testing.T) {
	type args struct {
		address string
		f       func(string) *services.ORouter
		opts    []rpc.Option
	}
	tests := []struct {
		name   string
		args   args
		repeat *args
		want   *rpcBuilder
	}{
		{name: "1. 初始化rpc空对想", args: args{address: ":1515", f: func(string) *services.ORouter { return nil }, opts: []rpc.Option{}},
			want: &rpcBuilder{httpBuilder: &httpBuilder{CustomerBuilder: map[string]interface{}{"main": rpc.New(":1515")}, fnGetRouter: func(string) *services.ORouter { return nil }}}},
		{name: "2. 初始化自定义rpc对象", args: args{address: ":1515", f: func(string) *services.ORouter { return nil }, opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
			want: &rpcBuilder{httpBuilder: &httpBuilder{CustomerBuilder: map[string]interface{}{"main": rpc.New(":1515", rpc.WithDisable(), rpc.WithTrace())},
				fnGetRouter: func(string) *services.ORouter { return nil }}}},
		{name: "3. 重复初始化自定义rpc对象", args: args{address: ":1515", f: func(string) *services.ORouter { return nil }, opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
			repeat: &args{address: ":1516", f: func(string) *services.ORouter { return nil }, opts: []rpc.Option{rpc.WithEnable(), rpc.WithDNS("0.0.0.0"), rpc.WithTrace()}},
			want: &rpcBuilder{httpBuilder: &httpBuilder{CustomerBuilder: map[string]interface{}{"main": rpc.New(":1516", rpc.WithEnable(), rpc.WithDNS("0.0.0.0"), rpc.WithTrace())},
				fnGetRouter: func(string) *services.ORouter { return nil }}}},
	}
	for _, tt := range tests {
		got := newRPC(tt.args.address, tt.args.f, tt.args.opts...)
		if tt.repeat != nil {
			got = newRPC(tt.repeat.address, tt.repeat.f, tt.repeat.opts...)
		}
		assert.Equal(t, tt.want.tp, got.tp, tt.name+",tp")
		assert.Equal(t, tt.want.fnGetRouter(""), got.fnGetRouter(""), tt.name+",fnGetRouter")
		assert.Equal(t, tt.want.CustomerBuilder, got.CustomerBuilder, tt.name+",CustomerBuilder")
	}
}

func Test_rpcBuilder_Load(t *testing.T) {
	tests := []struct {
		name   string
		fields *rpcBuilder
		want   CustomerBuilder
	}{
		{name: "1. 空路由,加载rpc路由配置", fields: &rpcBuilder{httpBuilder: &httpBuilder{fnGetRouter: func(string) *services.ORouter {
			return services.GetRouter(global.RPC)
		}, CustomerBuilder: make(map[string]interface{})}}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "2. 重复路由,加载rpc路由配置", fields: &rpcBuilder{httpBuilder: &httpBuilder{fnGetRouter: func(string) *services.ORouter {
			r := services.GetRouter(global.RPC)
			r.Add("path1", "service1", []string{"get"})
			r.Add("path1", "service1", []string{"get"})
			return r
		}, CustomerBuilder: make(map[string]interface{})}}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "3. 正常路由,加载rpc路由配置", fields: &rpcBuilder{httpBuilder: &httpBuilder{fnGetRouter: func(string) *services.ORouter {
			r := services.GetRouter(global.RPC)
			r.Add("path1", "service1", []string{"get"})
			return r
		}, CustomerBuilder: make(map[string]interface{})}}, want: CustomerBuilder{"router": router.NewRouters()}},
	}
	for _, tt := range tests {
		defer func() {
			if e := recover(); e != nil {
				assert.Equal(t, true, strings.Contains(types.GetString(e), "重复注册的服务"), tt.name)
			}
		}()
		tt.fields.fnGetRouter(global.RPC)
		tt.fields.Load()
		assert.Equal(t, tt.want, tt.fields.CustomerBuilder, tt.name)
	}
}
