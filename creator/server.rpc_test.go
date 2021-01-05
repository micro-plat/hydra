package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/assert"
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
		{name: "1. 初始化rpc空对想", args: args{address: ":1515", opts: []rpc.Option{}},
			want: &rpcBuilder{httpBuilder: &httpBuilder{BaseBuilder: map[string]interface{}{"main": rpc.New(":1515")}}}},
		{name: "2. 初始化自定义rpc对象", args: args{address: ":1515", opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
			want: &rpcBuilder{httpBuilder: &httpBuilder{BaseBuilder: map[string]interface{}{"main": rpc.New(":1515", rpc.WithDisable(), rpc.WithTrace())}}}},
		{name: "3. 重复初始化自定义rpc对象", args: args{address: ":1515", opts: []rpc.Option{rpc.WithDisable(), rpc.WithTrace()}},
			repeat: &args{address: ":1516", opts: []rpc.Option{rpc.WithEnable(), rpc.WithDNS("0.0.0.0"), rpc.WithTrace()}},
			want:   &rpcBuilder{httpBuilder: &httpBuilder{BaseBuilder: map[string]interface{}{"main": rpc.New(":1516", rpc.WithEnable(), rpc.WithDNS("0.0.0.0"), rpc.WithTrace())}}}},
	}
	for _, tt := range tests {
		got := newRPC(tt.args.address, tt.args.opts...)
		if tt.repeat != nil {
			got = newRPC(tt.repeat.address, tt.repeat.opts...)
		}
		assert.Equal(t, tt.want.tp, got.tp, tt.name+",tp")
		assert.Equal(t, tt.want.BaseBuilder, got.BaseBuilder, tt.name+",BaseBuilder")
	}
}
