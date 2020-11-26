package servers

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:liujinyin
//time:2020-10-15 15:57
//desc:测试黑名单中间件逻辑
func TestBlackList(t *testing.T) {
	type testCase struct {
		name        string
		isBool      bool
		blackOpts   []blacklist.Option
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "1.1 黑名单-配置不存在", isBool: false, blackOpts: []blacklist.Option{}, wantStatus: 200, wantContent: "", wantSpecial: ""},

		{name: "2.1 黑名单-配置存在-未启用-列表为空", isBool: true, blackOpts: []blacklist.Option{blacklist.WithDisable()}, wantStatus: 200, wantContent: "", wantSpecial: ""},
		{name: "2.2 黑名单-配置存在-未启用-黑名单IP", isBool: true, blackOpts: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.1")}, wantStatus: 200, wantContent: "", wantSpecial: ""},
		{name: "2.3 黑名单-配置存在-未启用-IP不在列表单内", isBool: true, blackOpts: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.3", "192.168.0.2")}, wantStatus: 200, wantContent: "", wantSpecial: ""},

		{name: "3.1 黑名单-配置存在-启用-列表为空", isBool: true, blackOpts: []blacklist.Option{blacklist.WithEnable()}, wantStatus: 200, wantContent: "", wantSpecial: "black"},
		{name: "3.2 黑名单-配置存在-启用-IP不在列表单内", isBool: true, blackOpts: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.2", "192.168.0.3")}, wantStatus: 200, wantContent: "", wantSpecial: "black"},
		{name: "3.3 黑名单-配置存在-启用-黑名单IP", isBool: true, blackOpts: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.1", "192.168.0.2")}, wantStatus: http.StatusForbidden, wantContent: "黑名单限制[192.168.0.1]不允许访问", wantSpecial: "black"},
	}
	for _, tt := range tests {

		mockConf := mocks.NewConfBy("middleware_black_test", "black")
		//初始化测试用例参数
		confB := mockConf.GetAPI()
		if tt.isBool {
			confB.BlackList(tt.blackOpts...)
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}

		//获取中间件
		handler := middleware.BlackList()

		//调用中间件
		handler(ctx)
		//断言结果
		gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()

		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}
