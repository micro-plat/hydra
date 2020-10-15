package mocks

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/test/assert"
)

func TestGetConf(t *testing.T) {
	conf := NewConf() //构建对象

	conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30))

	server := conf.GetAPIConf() //获取配置
	assert.Equal(t, server.GetMainConf().GetRootConf().GetString("address"), ":8081", "端口一致性检查")
	assert.Equal(t, server.GetMainConf().GetRootConf().GetInt("rTimeout"), 30, "超时时间检查")

}
