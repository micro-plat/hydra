package conf

import (
	"testing"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"

	"github.com/micro-plat/hydra/conf/server/api"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
)

func TestBlackListGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)).BlackList(blacklist.WithIP("192.168.0.*", "192.168.1.2"))
	server := conf.GetAPIConf() //获取配置
	bobj := blacklist.GetConf(server.GetMainConf())
	assert.Equal(t, bobj.Disable, false, "检查disblae值")
	assert.Equal(t, bobj.IPS, []string{"192.168.0.*", "192.168.1.2"}, "检查模板ip列表值")
	assert.Equal(t, bobj.IsDeny("192.168.1.2"), true, "检查IsDeny成功")
	assert.Equal(t, bobj.IsDeny("192.168.1.1"), false, "检查IsDeny失败")
	assert.Equal(t, bobj.IsDeny("192.168.0.1"), true, "检查IsDeny模糊匹配")
}
