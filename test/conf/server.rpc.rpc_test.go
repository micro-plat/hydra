/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestRPCNew(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []rpc.Option
		want    *rpc.Server
	}{
		{name: "默认初始化", opts: []rpc.Option{}, want: &rpc.Server{}},
		{name: "设置address初始化", address: ":8080", opts: []rpc.Option{}, want: &rpc.Server{Address: ":8080"}},
		{name: "设置option初始化", opts: []rpc.Option{rpc.WithTrace(), rpc.WithDNS("host1", "ip1"), rpc.WithHeaderReadTimeout(10), rpc.WithTimeout(11, 12)},
			want: &rpc.Server{RTimeout: 11, WTimeout: 12, RHTimeout: 10, Host: "host1", Domain: "ip1", Trace: true}},
		{name: "设置disable初始化", opts: []rpc.Option{rpc.WithDisable()}, want: &rpc.Server{Status: "stop"}},
		{name: "设置Enable初始化", opts: []rpc.Option{rpc.WithEnable()}, want: &rpc.Server{Status: "start"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rpc.New(tt.address, tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRpcGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	wantC := rpc.New(":8088", rpc.WithTrace())
	conf.RPC(":8088", rpc.WithTrace())
	got, err := rpc.GetConf(conf.GetRPCConf().GetMainConf())
	if err != nil {
		t.Errorf("rpcGetConf 获取配置对对象失败,err: %v", err)
	}
	assert.Equal(t, got, wantC, "检查对象是否满足预期")
}

func TestRpcGetConf1(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	conf.RPC(":8088")
	_, err := rpc.GetConf(conf.GetRPCConf().GetMainConf())
	if err == nil {
		t.Errorf("rpcGetConf 获取怕配置不能成功")
	}
}
