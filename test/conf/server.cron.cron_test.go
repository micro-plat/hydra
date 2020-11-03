/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
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
		{name: "默认初始化", args: []cron.Option{}, want: &cron.Server{Status: "start"}},
		{name: "初始化MasterSlave对象", args: []cron.Option{cron.WithMasterSlave()}, want: &cron.Server{Sharding: 1, Status: "start"}},
		{name: "初始化P2P对等模式对象", args: []cron.Option{cron.WithP2P()}, want: &cron.Server{Sharding: 0, Status: "start"}},
		{name: "初始化分片模式对等模式对象", args: []cron.Option{cron.WithSharding(10)}, want: &cron.Server{Sharding: 10, Status: "start"}},
		{name: "初始化剩余参数对象", args: []cron.Option{cron.WithTrace(), cron.WithTimeout(11)}, want: &cron.Server{Trace: true, Timeout: 11, Status: "start"}},
		{name: "初始化Disable参数对象", args: []cron.Option{cron.WithDisable()}, want: &cron.Server{Status: "stop"}},
		{name: "初始化Enable参数对象", args: []cron.Option{cron.WithEnable()}, want: &cron.Server{Status: "start"}},
	}
	for _, tt := range tests {
		got := cron.New(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCronGetConf(t *testing.T) {
	type test struct {
		name string
		opts []cron.Option
		want *cron.Server
	}

	conf := mocks.NewConf()
	test1 := test{name: "节点不存在,获取默认对象", opts: []cron.Option{}, want: &cron.Server{}}
	obj, err := cron.GetConf(conf.GetCronConf().GetServerConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, test1.want, obj, test1.name)

	tests := []test{
		{name: "节点为空,获取默认对象", opts: []cron.Option{}, want: cron.New()},
		{name: "正常对象获取",
			opts: []cron.Option{cron.WithTrace(), cron.WithMasterSlave()},
			want: cron.New(cron.WithTrace(), cron.WithMasterSlave())},
	}
	for _, tt := range tests {
		conf.CRON(tt.opts...)
		obj, err := cron.GetConf(conf.GetCronConf().GetServerConf())
		assert.Equal(t, nil, err, tt.name+",err")
		assert.Equal(t, tt.want, obj, tt.name)
	}

	//异常的json数据  需要完善注册中心后测试(借鉴blacklist的写法)
}
