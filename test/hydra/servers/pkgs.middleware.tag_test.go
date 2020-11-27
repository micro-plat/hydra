package servers

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

// tag中间件只有在websocket在使用    现在websocket服务暂时不开放
func TestTag(t *testing.T) {
	confMock := mocks.NewConfBy("middleware_tag_test1", "tag")
	confMock.API(":5454")
	confMock.RPC(":6541")
	confMock.MQC("redis://redisname")
	confMock.Web(":8541")
	confMock.WS(":5214")
	confMock.CRON()

	tests := []struct {
		name        string
		requstID    string
		conf        app.IAPPConf
		wantSpecial string
	}{
		{name: "1. 配置serverType为mqc", conf: confMock.GetMQCConf(), wantSpecial: "mqc"},
		{name: "2. 配置serverType为cron", conf: confMock.GetCronConf(), wantSpecial: "cron"},
		{name: "3. 配置serverType为api", conf: confMock.GetAPIConf(), wantSpecial: "api"},
		{name: "4. 配置serverType为web", conf: confMock.GetWebConf(), wantSpecial: "web"},
		{name: "5. 配置serverType为ws", conf: confMock.GetWSConf(), wantSpecial: "ws"},
		{name: "6. 配置serverType为rpc", conf: confMock.GetRPCConf(), wantSpecial: "rpc"},
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
