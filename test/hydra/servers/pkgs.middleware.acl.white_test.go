package hydra

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:taoshouyin
//time:2020-11-11
//desc:测试白名单中间件逻辑
func TestWhiteList(t *testing.T) {
	type testCase struct {
		name        string
		whiteOpts   []whitelist.Option
		isSet       bool
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "白名单-未配置", isSet: false, wantStatus: 200, wantContent: "", wantSpecial: "", whiteOpts: []whitelist.Option{}},
		{name: "白名单-配置未启动", isSet: true, wantStatus: 200, wantContent: "", wantSpecial: "", whiteOpts: []whitelist.Option{whitelist.WithDisable()}},
		{name: "白名单-配置启动-不存在路径,不存在ip", isSet: true, wantStatus: 510, wantContent: "", wantSpecial: "",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{}))}},
		{name: "白名单-配置启动-存在路径,不存在ip", isSet: true, wantStatus: 510, wantContent: "", wantSpecial: "",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{"/whitelist/test"}))}},
		{name: "白名单-配置启动-不匹配路径,不匹配ip", isSet: true, wantStatus: 200, wantContent: "", wantSpecial: "white",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{"/whitelist/test1"}, "192.168.0.11"))}},
		{name: "白名单-配置启动-匹配路径,不匹配ip", isSet: true, wantStatus: 403, wantContent: "", wantSpecial: "white",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{"/whitelist/test"}, "192.168.0.11"))}},
		{name: "白名单-配置启动-不匹配路径,匹配ip", isSet: true, wantStatus: 200, wantContent: "", wantSpecial: "white",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{"/whitelist/test1"}, "192.168.0.1"))}},
		{name: "白名单-配置启动-匹配路径,匹配ip", isSet: true, wantStatus: 200, wantContent: "", wantSpecial: "white",
			whiteOpts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList([]string{"/whitelist/test"}, "192.168.0.1"))}},
	}
	for _, tt := range tests {
		mockConf := mocks.NewConf()
		//初始化测试用例参数
		confB := mockConf.GetAPI()
		if tt.isSet {
			confB.WhiteList(tt.whiteOpts...)
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockRequestPath: "/whitelist/test",
				},
			},
			MockServerConf: serverConf,
		}

		//获取中间件
		handler := middleware.WhiteList()
		//调用中间件
		handler(ctx)
		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}
