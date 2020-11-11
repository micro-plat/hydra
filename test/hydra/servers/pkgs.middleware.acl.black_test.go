package hydra

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
		blackOpts   []blacklist.Option
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{
			name:        "黑名单-未启用",
			blackOpts:   []blacklist.Option{blacklist.WithDisable()},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name:        "黑名单-未启用-黑名单IP",
			blackOpts:   []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.1")},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name:        "黑名单-启用-不在黑名单内的IP",
			blackOpts:   []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.2", "192.168.0.3")},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "black",
		},
		{
			name:        "黑名单-启用-黑名单IP",
			blackOpts:   []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.1")},
			wantStatus:  http.StatusForbidden,
			wantContent: "",
			wantSpecial: "black",
		},
	}
	for _, tt := range tests {

		mockConf := mocks.NewConf()
		//初始化测试用例参数
		mockConf.GetAPI().BlackList(tt.blackOpts...)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockUser:       &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse:   &mocks.MockResponse{MockStatus: 200},
			MockAPPConf: serverConf,
		}

		//获取中间件
		handler := middleware.BlackList()

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
