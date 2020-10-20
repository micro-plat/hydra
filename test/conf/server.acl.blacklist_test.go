/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestBlackListNew(t *testing.T) {
	tests := []struct {
		name     string
		opts     []blacklist.Option
		want     *blacklist.BlackList
		denyIP   string
		wantDeny bool
	}{
		{name: "初始化空对象", opts: nil, want: &blacklist.BlackList{IPS: []string{}}, denyIP: "127.0.0.1", wantDeny: false},
		{name: "初始化单ip对象", opts: []blacklist.Option{blacklist.WithIP("19.168.0.101")}, want: &blacklist.BlackList{IPS: []string{"19.168.0.101"}}, denyIP: "19.168.0.101", wantDeny: true},
		{name: "初始化多对象Enable", opts: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("19.168.0.101")}, want: &blacklist.BlackList{Disable: false, IPS: []string{"19.168.0.101"}}, denyIP: "19.168.0.101", wantDeny: true},
		{name: "初始化多对象Disable", opts: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("19.168.0.101")}, want: &blacklist.BlackList{Disable: true, IPS: []string{"19.168.0.101"}}, denyIP: "19.168.0.101", wantDeny: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blacklist.New(tt.opts...)

			//比对初始化的disable是否相同
			assert.Equal(t, tt.want.Disable, got.Disable, tt.name+",disable")

			//比对初始化的iplist是否相同
			assert.Equal(t, tt.want.IPS, got.IPS, tt.name+",IPS")

			//测试私有匹配参数是否成功赋值
			denygot := got.IsDeny(tt.denyIP)
			assert.Equal(t, tt.wantDeny, denygot, tt.name+",deny")
		})
	}
}

//该方法依赖于conf-match函数的测试,match没有问题,该函数也就没有问题啦
func TestBlackList_IsDeny(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		ips  []string
		want bool
	}{
		{name: "空对象路径匹配", ip: "/t1/tt", ips: []string{}, want: false},
		{name: "全路径精确匹配", ip: "/t1/t2/tt", ips: []string{"/t1/t2/tt"}, want: true},
		{name: "单级模糊路径匹配", ip: "/t1/t2", ips: []string{"/t1/*"}, want: true},
		{name: "多级模糊路径匹配", ip: "/t1/t2/tt", ips: []string{"/t1/**"}, want: true},
		{name: "空对象ip匹配", ip: "192.168.5.101", ips: []string{}, want: false},
		{name: "全ip精确匹配", ip: "127.0.0.1", ips: []string{"127.0.0.1"}, want: true},
		{name: "单级模糊ip匹配", ip: "192.168.0.1", ips: []string{"192.168.0.*"}, want: true},
		{name: "多级模糊ip匹配", ip: "192.168.0.1", ips: []string{"192.168.**"}, want: true},
	}
	for _, tt := range tests {
		w := blacklist.New(blacklist.WithIP(tt.ips...))
		got := w.IsDeny(tt.ip)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestBlackListGetConf(t *testing.T) {
	tests := []struct {
		name string
		opts []blacklist.Option
		want *blacklist.BlackList
	}{
		{name: "节点不存在,获取默认对象", opts: []blacklist.Option{}, want: &blacklist.BlackList{Disable: true}},
		{name: "节点为空,获取默认对象", opts: []blacklist.Option{}, want: blacklist.New()},
		{name: "正常对象获取", opts: []blacklist.Option{blacklist.WithIP("192.168.0.*", "192.168.1.2")}, want: blacklist.New(blacklist.WithIP("192.168.0.*", "192.168.1.2"))},
	}

	//初始化服务conf配置对象
	conf := mocks.NewConf()
	confB := conf.API(":8081")
	for _, tt := range tests {
		if !strings.EqualFold(tt.name, "节点不存在,获取默认对象") {
			confB.BlackList(tt.opts...)
		}
		obj := blacklist.GetConf(conf.GetAPIConf().GetMainConf())
		assert.Equal(t, tt.want, obj, tt.name)
	}

	// json数据不合法,现在还不能测试   需要等待注册中心监听完善后测试
	// path := conf.GetAPIConf().GetMainConf().GetSubConfPath("acl", "black.list")
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		if !strings.Contains(e.(string), fmt.Sprintf("获取%s配置失败", path)) {
	// 			t.Error("json错误,返回了未知的错误信息")
	// 		}
	// 	}
	// }()
	// conf.Registry.Update(path, "错误的json字符串")
	// ch, _ := conf.Registry.WatchValue(path)
	// select {
	// case <-time.After(3 * time.Second):
	// 	return
	// case <-ch:
	// 	bobj = blacklist.GetConf(conf.GetAPIConf().GetMainConf())
	// 	t.Errorf("%v", bobj)
	// }
}
