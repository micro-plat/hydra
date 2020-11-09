/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestAPINew(t *testing.T) {
	tests := []struct {
		name    string
		address string
		opts    []api.Option
		want    *api.Server
	}{
		{name: "默认初始化", want: &api.Server{Status: "start"}},
		{name: "设置address初始化", address: ":8080", want: &api.Server{Address: ":8080", Status: "start"}},
		{name: "设置option初始化", opts: []api.Option{api.WithTrace(), api.WithDNS("ip1"), api.WithHeaderReadTimeout(10), api.WithTimeout(11, 12)},
			want: &api.Server{RTimeout: 11, WTimeout: 12, RHTimeout: 10, Domain: "ip1", Trace: true, Status: "start"}},
		{name: "设置disable初始化", opts: []api.Option{api.WithDisable()}, want: &api.Server{Status: "stop"}},
		{name: "设置Enable初始化", opts: []api.Option{api.WithEnable()}, want: &api.Server{Status: "start"}},
	}
	for _, tt := range tests {
		got := api.New(tt.address, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestAPIGetConf(t *testing.T) {
	type test struct {
		name string
		opts []api.Option
		want *api.Server
	}

	conf := mocks.NewConf()
	//getconf的address默认值是8080
	test1 := test{name: "节点不存在,获取默认对象", opts: []api.Option{}, want: &api.Server{Address: ":8080", Status: "start"}}
	obj, err := api.GetConf(conf.GetAPIConf().GetServerConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, test1.want, obj, test1.name)

	tests := []test{
		{name: "节点为空,获取默认对象", opts: []api.Option{}, want: api.New(":8080")},
		{name: "正常对象获取",
			opts: []api.Option{api.WithTrace(), api.WithDNS("ip1"), api.WithHeaderReadTimeout(10), api.WithTimeout(11, 12)},
			want: api.New(":8080", api.WithTrace(), api.WithDNS("ip1"), api.WithHeaderReadTimeout(10), api.WithTimeout(11, 12))},
	}
	for _, tt := range tests {
		conf.API(":8080", tt.opts...)
		obj, err := api.GetConf(conf.GetAPIConf().GetServerConf())
		assert.Equal(t, nil, err, tt.name+",err")
		assert.Equal(t, tt.want, obj, tt.name)
	}
}
