/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"reflect"
	"strings"
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
		t.Run(tt.name, func(t *testing.T) {
			if got := mqc.New(tt.addr, tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMqcGetConf(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("MqcGetConf 获取失败,err:%v", err)
		}
	}()
	conf := mocks.NewConf() //构建对象
	wantC := mqc.New("redis://192.196.0.1", mqc.WithTrace())
	conf.MQC("redis://192.196.0.1", mqc.WithTrace())
	got := mqc.GetConf(conf.GetMQCConf().GetMainConf())
	assert.Equal(t, got, wantC, "检查对象是否满足预期")
}

func TestMqcGetConf1(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			err1 := err.(error)
			if !strings.Contains(err1.Error(), "配置有误") {
				t.Errorf("MqcGetConf 获取怕配置不能成功")
			} else {
				t.Errorf("MqcGetConf 配置有误,err:%v", err)
			}
		}
	}()
	conf := mocks.NewConf() //构建对象
	conf.API(":8000")
	mqc.GetConf(conf.GetAPIConf().GetMainConf())

}
