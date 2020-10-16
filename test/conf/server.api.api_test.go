package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

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
