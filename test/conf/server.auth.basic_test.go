/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/auth/basic"

	"github.com/micro-plat/hydra/conf"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestBasicGetConf(t *testing.T) {

	tests := []struct {
		name string
		args func() conf.IMainConf
		want *basic.BasicAuth
	}{
		{name: "未设置basic节点", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081")
			return conf.GetAPIConf().GetMainConf()
		}, want: &basic.BasicAuth{Disable: true}},
		{name: "配置参数正确", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081").Basic(basic.WithUP("t1", "123"), basic.WithExcludes("/t1/t12"))
			return conf.GetAPIConf().GetMainConf()
		}, want: basic.NewBasic(basic.WithUP("t1", "123"), basic.WithExcludes("/t1/t12"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("apiKeyGetConf 获取配置对对象失败,err: %v", err)
				}
			}()
			got := basic.GetConf(tt.args())
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}
