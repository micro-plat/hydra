package hydra

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:liujinyin
//time:2020-10-15 15:57
//desc:测试黑名单中间件逻辑
func TestBlackList(t *testing.T) {
	type testCase struct {
		name   string
		ctx    middleware.IMiddleContext
		init   func(tc *testCase)
		assert func(t *testing.T, ctx middleware.IMiddleContext)
	}

	tests := []*testCase{
		{
			name: "黑名单-未启用",
			init: func(tc *testCase) {
				mockConf := mocks.NewConf()
				mockConf.GetAPI().BlackList(blacklist.WithDisable())
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mockMiddleContext{
					MockUser: &MockUser{
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assert: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*MockResponse)
				if !ok {
					t.Errorf("中间件流程执行不是预期的类型")
				}
				if v.MockStatus != 200 {
					t.Errorf("中间件流程执行不是预期的过错;expect:%d,actual:%d", 200, v.MockStatus)
				}
			},
		},
		{
			name: "黑名单-启用-不在黑名单内的IP",
			init: func(tc *testCase) {
				mockConf := mocks.NewConf()
				mockConf.GetAPI().BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.2", "192.168.0.3"))
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mockMiddleContext{
					MockUser: &MockUser{
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assert: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*MockResponse)
				if !ok {
					t.Errorf("中间件流程执行不是预期的类型")
				}
				if v.MockStatus != 200 {
					t.Errorf("中间件流程执行不是预期的过错;expect:%d,actual:%d", 200, v.MockStatus)
				}
			},
		},
		{
			name: "黑名单-启用-黑名单IP",
			init: func(tc *testCase) {
				mockConf := mocks.NewConf()
				mockConf.GetAPI().BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.1", "192.168.0.2", "192.168.0.3"))
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mockMiddleContext{
					MockUser: &MockUser{
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assert: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*MockResponse)
				if !ok {
					t.Errorf("中间件流程执行不是预期的类型")
				}
				if v.MockStatus != http.StatusForbidden {
					t.Errorf("中间件流程执行不是预期的过错;expect:%d,actual:%d", http.StatusForbidden, v.MockStatus)
				}
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init(tt)
			handler := middleware.BlackList()
			handler(tt.ctx)
			tt.assert(t, tt.ctx)
		})
	}
}
