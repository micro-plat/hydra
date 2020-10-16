/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestCronNew(t *testing.T) {
	tests := []struct {
		name string
		args []cron.Option
		want *cron.Server
	}{
		{name: "默认初始化", args: []cron.Option{}, want: &cron.Server{}},
		{name: "初始化MasterSlave对象", args: []cron.Option{cron.WithMasterSlave()}, want: &cron.Server{Sharding: 1}},
		{name: "初始化P2P对等模式对象", args: []cron.Option{cron.WithP2P()}, want: &cron.Server{Sharding: 0}},
		{name: "初始化分片模式对等模式对象", args: []cron.Option{cron.WithSharding(10)}, want: &cron.Server{Sharding: 10}},
		{name: "初始化剩余参数对象", args: []cron.Option{cron.WithTrace(), cron.WithTimeout(11)}, want: &cron.Server{Trace: true, Timeout: 11}},
		{name: "初始化Disable参数对象", args: []cron.Option{cron.WithDisable()}, want: &cron.Server{Status: "stop"}},
		{name: "初始化Enable参数对象", args: []cron.Option{cron.WithEnable()}, want: &cron.Server{Status: "start"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cron.New(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	wantC := cron.New(cron.WithTrace())
	conf.CRON(cron.WithTrace())
	got, err := cron.GetConf(conf.GetCronConf().GetMainConf())
	if err != nil {
		t.Errorf("cronGetConf 获取配置对对象失败,err: %v", err)
	}
	assert.Equal(t, got, wantC, "检查对象是否满足预期")
}

func TestCronGetConf1(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	conf.API(":8000")
	_, err := cron.GetConf(conf.GetAPIConf().GetMainConf())
	if err == nil {
		t.Errorf("cronGetConf 获取怕配置不能成功")
	}
}
