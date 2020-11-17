package servers

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestTag(t *testing.T) {
	tests := []struct {
		name        string
		requstID    string
		conf        app.IAPPConf
		wantSpecial string
	}{
		{name: "配置serverType为rpc", conf: mocks.NewConf().GetCronConf(), wantSpecial: "cron"},
		{name: "配置serverType为api", conf: mocks.NewConf().GetAPIConf(), wantSpecial: "api"},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockResponse: &mocks.MockResponse{MockStatus: 200, MockHeader: map[string][]string{}},
			MockAPPConf:  tt.conf,
		}

		//调用中间件
		handler := middleware.Tag()
		handler(ctx)

		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
	}
}
