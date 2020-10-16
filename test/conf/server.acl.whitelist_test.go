/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestWhiteListGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	confObj := conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30))
	confObj.WhiteList(
		whitelist.WithIPList(
			whitelist.NewIPList("/t1/t2", []string{"192.168.0.*", "192.168.1.2"}...),
			whitelist.NewIPList("/t1/t2/*", []string{"192.168.5.*", "192.168.5.2"}...),
		))
	server := conf.GetAPIConf() //获取配置
	bobj := whitelist.GetConf(server.GetMainConf())
	assert.Equal(t, bobj.Disable, false, "检查disblae值")
	assert.Equal(t, bobj.IPS[0].Requests, []string{"/t1/t2"}, "检查模板路径列表值")
	assert.Equal(t, bobj.IPS[0].IPS, []string{"192.168.0.*", "192.168.1.2"}, "检查模板ip列表值")
	assert.Equal(t, bobj.IPS[1].Requests, []string{"/t1/t2/*"}, "检查模板路径列表值1")
	assert.Equal(t, bobj.IPS[1].IPS, []string{"192.168.5.2", "192.168.5.*"}, "检查模板ip列表值1")
	// assert.Equal(t, bobj.IsAllow("path", "192.168.1.2"), true, "匹配测试")  待匹配方案确定后在测试
}
