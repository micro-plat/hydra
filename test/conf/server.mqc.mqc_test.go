/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestMqcNew(t *testing.T) {

	tests := []struct {
		name string
		addr string
		opts []mqc.Option
		want *mqc.Server
	}{
		{name: "默认初始化", addr: "redis://11", opts: []mqc.Option{}, want: &mqc.Server{Addr: "redis://11"}},
		{name: "初始化MasterSlave对象", addr: "redis://11", opts: []mqc.Option{mqc.WithMasterSlave()}, want: &mqc.Server{Addr: "redis://11", Sharding: 1}},
		{name: "初始化P2P对等模式对象", addr: "redis://11", opts: []mqc.Option{mqc.WithP2P()}, want: &mqc.Server{Addr: "redis://11", Sharding: 0}},
		{name: "初始化分片模式对等模式对象", addr: "redis://11", opts: []mqc.Option{mqc.WithSharding(10)}, want: &mqc.Server{Addr: "redis://11", Sharding: 10}},
		{name: "初始化剩余参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithTrace(), mqc.WithTimeout(11)}, want: &mqc.Server{Addr: "redis://11", Trace: true, Timeout: 11}},
		{name: "初始化Disable参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithDisable()}, want: &mqc.Server{Addr: "redis://11", Status: "stop"}},
		{name: "初始化Enable参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithEnable()}, want: &mqc.Server{Addr: "redis://11", Status: "start"}},
	}
	for _, tt := range tests {
		got := mqc.New(tt.addr, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestMQCGetConf(t *testing.T) {
	type test struct {
		name string
		opts []mqc.Option
		want *mqc.Server
	}

	conf := mocks.NewConf()
	//mqc的节点不存在需要报 panic
	// test1 := test{name: "节点不存在,获取默认对象", opts: []mqc.Option{}, want: &mqc.Server{}}
	// obj := mqc.GetConf(conf.GetMQCConf().GetMainConf())
	// assert.Equal(t, test1.want, obj, test1.name)
	tests := []test{
		{name: "正常对象获取",
			opts: []mqc.Option{mqc.WithTrace(), mqc.WithMasterSlave()},
			want: mqc.New("redis://192.196.0.1", mqc.WithTrace(), mqc.WithMasterSlave())},
	}
	for _, tt := range tests {
		conf.MQC("redis://192.196.0.1", tt.opts...)
		obj := mqc.GetConf(conf.GetMQCConf().GetMainConf())
		assert.Equal(t, tt.want, obj, tt.name)
	}

	//异常的json数据  需要完善注册中心后测试(借鉴blacklist的写法)
}
