/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
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
		got := rpc.New(tt.address, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRPCGetConf(t *testing.T) {
	type test struct {
		name string
		opts []rpc.Option
		want *rpc.Server
	}

	conf := mocks.NewConf()
	test1 := test{name: "节点不存在,获取默认对象", opts: []rpc.Option{}, want: &rpc.Server{Address: ":8090"}}
	obj, err := rpc.GetConf(conf.GetRPCConf().GetMainConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, test1.want, obj, test1.name)

	tests := []test{
		{name: "节点为空,获取默认对象", opts: []rpc.Option{}, want: rpc.New(":8090")},
		{name: "正常对象获取",
			opts: []rpc.Option{rpc.WithTrace(), rpc.WithHeaderReadTimeout(10)},
			want: rpc.New(":8090", rpc.WithTrace(), rpc.WithHeaderReadTimeout(10))},
	}
	for _, tt := range tests {
		conf.RPC(":8090", tt.opts...)
		obj, err := rpc.GetConf(conf.GetCronConf().GetMainConf())
		assert.Equal(t, nil, err, tt.name+",err")
		assert.Equal(t, tt.want, obj, tt.name)
	}

	//异常的json数据  需要完善注册中心后测试(借鉴blacklist的写法)
}
