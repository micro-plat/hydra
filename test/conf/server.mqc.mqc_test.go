package conf

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

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
