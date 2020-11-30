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
	"github.com/micro-plat/lib4go/types"
)

func TestMqcNew(t *testing.T) {

	tests := []struct {
		name string
		addr string
		opts []mqc.Option
		want *mqc.Server
	}{
		{name: "1. Conf-MqcNew-默认初始化", addr: "redis://11", opts: []mqc.Option{}, want: &mqc.Server{Addr: "redis://11", Status: "start"}},
		{name: "2. Conf-MqcNew-初始化MasterSlave对象", addr: "redis://11", opts: []mqc.Option{mqc.WithMasterSlave()}, want: &mqc.Server{Addr: "redis://11", Sharding: 1, Status: "start"}},
		{name: "3. Conf-MqcNew-初始化P2P对等模式对象", addr: "redis://11", opts: []mqc.Option{mqc.WithP2P()}, want: &mqc.Server{Addr: "redis://11", Sharding: 0, Status: "start"}},
		{name: "4. Conf-MqcNew-初始化分片模式对等模式对象", addr: "redis://11", opts: []mqc.Option{mqc.WithSharding(10)}, want: &mqc.Server{Addr: "redis://11", Sharding: 10, Status: "start"}},
		{name: "5. Conf-MqcNew-初始化剩余参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithTrace()}, want: &mqc.Server{Addr: "redis://11", Trace: true, Status: "start"}},
		{name: "6. Conf-MqcNew-初始化Disable参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithDisable()}, want: &mqc.Server{Addr: "redis://11", Status: "stop"}},
		{name: "7. Conf-MqcNew-初始化Enable参数对象", addr: "redis://11", opts: []mqc.Option{mqc.WithEnable()}, want: &mqc.Server{Addr: "redis://11", Status: "start"}},
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

	conf := mocks.NewConfBy("hydraconf_mqc_test2", "mqcmain")
	defer func() {
		e := recover()
		if e != nil {
			assert.Equal(t, "未指定mqc服务器配置", types.GetString(e), "节点不存在,获取默认对象")
		}
	}()
	test1 := test{name: "1. Conf-MQCGetConf-节点不存在,获取默认对象", opts: []mqc.Option{}, want: &mqc.Server{}}
	obj, err := mqc.GetConf(conf.GetMQCConf().GetServerConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, test1.want, obj, test1.name)
	tests := []test{
		{name: "2. Conf-MQCGetConf-正常对象获取",
			opts: []mqc.Option{mqc.WithTrace(), mqc.WithMasterSlave()},
			want: mqc.New("redis://192.196.0.1", mqc.WithTrace(), mqc.WithMasterSlave())},
	}
	for _, tt := range tests {
		conf.MQC("redis://192.196.0.1", tt.opts...)
		obj, err := mqc.GetConf(conf.GetMQCConf().GetServerConf())
		assert.Equal(t, nil, err, tt.name+",err")
		assert.Equal(t, tt.want, obj, tt.name)
	}
}
