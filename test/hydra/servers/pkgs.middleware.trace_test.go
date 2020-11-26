package servers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestTrace(t *testing.T) {
	confMock := mocks.NewConfBy("middleware_trace_test", "trace")
	confMock.Web(":8541")
	confMock.WS(":5214")
	tests := []struct {
		name           string
		conf           func() app.IAPPConf
		requestMap     map[string]interface{}
		responseStatus int
		serverType     string
		wantSpecial    string
		wantDebug1     string
		wantDebug2     string
	}{
		{name: "1.1 Trace-api-未配置trace", serverType: "api", conf: func() app.IAPPConf { confMock.API(":5454"); return confMock.GetAPIConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "1.2 Trace-api-配置trace", serverType: "api", conf: func() app.IAPPConf { confMock.API(":5454", api.WithTrace()); return confMock.GetAPIConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},

		{name: "2.1 Trace-rpc未配置trace", serverType: "rpc", conf: func() app.IAPPConf { confMock.RPC(":6541"); return confMock.GetRPCConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "2.2 Trace-rpc配置trace", serverType: "rpc", conf: func() app.IAPPConf { confMock.RPC(":5454", rpc.WithTrace()); return confMock.GetRPCConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},

		{name: "3.1 Trace-mqc未配置trace", serverType: "mqc", conf: func() app.IAPPConf { confMock.MQC("redis://redisname"); return confMock.GetMQCConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "3.2 Trace-mqc配置trace", serverType: "mqc", conf: func() app.IAPPConf { confMock.MQC("redis://redisname", mqc.WithTrace()); return confMock.GetMQCConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},

		{name: "4.1 Trace-cron未配置trace", serverType: "cron", conf: func() app.IAPPConf { confMock.CRON(); return confMock.GetCronConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "4.2 Trace-cron配置trace", serverType: "cron", conf: func() app.IAPPConf { confMock.CRON(cron.WithTrace()); return confMock.GetCronConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},

		{name: "5.1 Trace-web未配置trace", serverType: "web", conf: func() app.IAPPConf { confMock.Web(":8541"); return confMock.GetWebConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "5.2 Trace-web配置trace", serverType: "web", conf: func() app.IAPPConf { confMock.Web(":8541", api.WithTrace()); return confMock.GetWebConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},

		{name: "6.1 Trace-ws未配置trace", serverType: "ws", conf: func() app.IAPPConf { confMock.WS(":5214"); return confMock.GetWSConf() }, responseStatus: 200, wantSpecial: ""},
		{name: "6.2 Trace-ws配置trace", serverType: "ws", conf: func() app.IAPPConf { confMock.WS(":5214", api.WithTrace()); return confMock.GetWSConf() }, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockNext:     func() { fmt.Println("output") },
			MockUser:     &mocks.MockUser{MockClientIP: "127.0.0.1", MockRequestID: "06c6fb24c"},
			MockRequest:  &mocks.MockRequest{MockQueryMap: tt.requestMap},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus},
			MockAPPConf:  tt.conf(),
		}

		//构建的新的os.Stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//调用中间件
		handler := middleware.Trace()
		handler(ctx)
		time.Sleep(time.Second * 3)

		//获取输出
		w.Close()
		out, err := ioutil.ReadAll(r)
		assert.Equalf(t, false, err != nil, tt.name)

		//还原os.Stdout
		os.Stdout = rescueStdout
		assert.Equalf(t, true, strings.Contains(string(out), "output"), tt.name)

		if tt.wantSpecial != "" {
			gotSpecial := ctx.Response().GetSpecials()
			assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantDebug1), tt.name)
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantDebug2), tt.name)
		}
	}
}
