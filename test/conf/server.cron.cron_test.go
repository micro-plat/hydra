package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

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
