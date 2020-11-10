/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/registry/pub"
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
		{name: "默认初始化", opts: []rpc.Option{}, want: &rpc.Server{Status: "start"}},
		{name: "设置address初始化", address: ":8080", opts: []rpc.Option{}, want: &rpc.Server{Address: ":8080", Status: "start"}},
		{name: "设置option初始化", opts: []rpc.Option{rpc.WithTrace(), rpc.WithDNS("host1", "ip1")},
			want: &rpc.Server{Host: "host1", Domain: "ip1", Trace: true, Status: "start"}},
		{name: "设置disable初始化", opts: []rpc.Option{rpc.WithDisable()}, want: &rpc.Server{Status: "stop"}},
		{name: "设置Enable初始化", opts: []rpc.Option{rpc.WithEnable()}, want: &rpc.Server{Status: "start"}},
	}
	for _, tt := range tests {
		got := rpc.New(tt.address, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func xTestRPCGetConf(t *testing.T) {
	type test struct {
		name string
		opts []rpc.Option
		want *rpc.Server
	}

	conf := mocks.NewConf()
	pub.New(conf.GetRPCConf().GetServerConf())

	time.Sleep(time.Second)

	test1 := test{name: "节点不存在,获取默认对象", opts: []rpc.Option{}, want: &rpc.Server{Address: ":8090"}}
	obj, err := rpc.GetConf(conf.GetRPCConf().GetServerConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, test1.want, obj, test1.name)

	tests := []test{
		{name: "节点为空,获取默认对象", opts: []rpc.Option{}, want: rpc.New(":8090")},
		{name: "正常对象获取",
			opts: []rpc.Option{rpc.WithTrace()},
			want: rpc.New(":8090", rpc.WithTrace())},
	}
	for _, tt := range tests {
		conf.RPC(":8090", tt.opts...)
		obj, err := rpc.GetConf(conf.GetCronConf().GetServerConf())
		assert.Equal(t, nil, err, tt.name+",err")
		assert.Equal(t, tt.want, obj, tt.name)
	}
}
