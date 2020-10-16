/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		address string
		opts    []api.Option
	}
	tests := []struct {
		name string
		args args
		want *api.Server
	}{
		{name: "默认初始化", args: args{}, want: &api.Server{}},
		{name: "设置address初始化", args: args{address: ":8080"}, want: &api.Server{Address: ":8080"}},
		{name: "设置option初始化", args: args{opts: []api.Option{api.WithTrace(), api.WithDNS("host1", "ip1"), api.WithHeaderReadTimeout(10), api.WithTimeout(11, 12)}},
			want: &api.Server{RTimeout: 11, WTimeout: 12, RHTimeout: 10, Host: "host1", Domain: "ip1", Trace: true}},
		{name: "设置disable初始化", args: args{opts: []api.Option{api.WithDisable()}}, want: &api.Server{Status: "stop"}},
		{name: "设置Enable初始化", args: args{opts: []api.Option{api.WithEnable()}}, want: &api.Server{Status: "start"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := api.New(tt.args.address, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIGetConf(t *testing.T) {
	conf := mocks.NewConfBy("hydar", "apiconftest") //构建对象
	wantC := api.New(":8081", api.WithHeaderReadTimeout(30))
	conf.API(":8081", api.WithHeaderReadTimeout(30))
	got, err := api.GetConf(conf.GetAPIConf().GetMainConf())
	if err != nil {
		t.Errorf("apiGetConf 获取配置对对象失败,err: %v", err)
	}
	assert.Equal(t, got, wantC, "检查对象是否满足预期")
}

func TestAPIGetConf1(t *testing.T) {
	conf := mocks.NewConfBy("hydar", "apiconftest") //构建对象
	conf.CRON()
	_, err := api.GetConf(conf.GetCronConf().GetMainConf())
	if err == nil {
		t.Errorf("apiGetConf 获取怕配置不能成功")
	}
}
