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
		name          string
		ctx           middleware.IMiddleContext
		initCaseParam func(tc *testCase)
		assertResult  func(t *testing.T, ctx middleware.IMiddleContext)
	}

	tests := []*testCase{
		{
			name: "黑名单-未启用",
			initCaseParam: func(tc *testCase) {
				mockConf := mocks.NewConf()
				//设置不启用
				mockConf.GetAPI().BlackList(blacklist.WithDisable())
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mocks.MiddleContext{
					MockUser: &mocks.MockUser{
						//设置客户端IP
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &mocks.MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assertResult: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*mocks.MockResponse)
				assert.Equal(t, true, ok, "中间件流程执行不是预期的类型")
				assert.Equalf(t, 200, v.MockStatus, "中间件流程执行不是预期的过错;expect:%d,actual:%d", 200, v.MockStatus)
			},
		},
		{
			name: "黑名单-启用-不在黑名单内的IP",
			initCaseParam: func(tc *testCase) {
				mockConf := mocks.NewConf()
				//设置启用，设置限定的IP
				mockConf.GetAPI().BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.2", "192.168.0.3"))
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mocks.MiddleContext{
					MockUser: &mocks.MockUser{
						//设置不在黑名单的客户端IP
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &mocks.MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assertResult: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*mocks.MockResponse)
				assert.Equal(t, true, ok, "中间件流程执行不是预期的类型")
				assert.Equalf(t, 200, v.MockStatus, "中间件流程执行不是预期的过错;expect:%d,actual:%d", 200, v.MockStatus)
			},
		},
		{
			name: "黑名单-启用-黑名单IP",
			initCaseParam: func(tc *testCase) {
				mockConf := mocks.NewConf()
				//设置启用，设置限定的IP
				mockConf.GetAPI().BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.1", "192.168.0.2", "192.168.0.3"))
				serverConf := mockConf.GetAPIConf()
				tc.ctx = &mocks.MiddleContext{
					MockUser: &mocks.MockUser{
						//设置客户端IP在黑名单内
						MockClientIP: "192.168.0.1",
					},
					MockResponse:   &mocks.MockResponse{MockStatus: 200},
					MockServerConf: serverConf,
				}
			},
			assertResult: func(t *testing.T, ctx middleware.IMiddleContext) {
				resp := ctx.Response()
				v, ok := resp.(*mocks.MockResponse)
				assert.Equal(t, true, ok, "中间件流程执行不是预期的类型")
				assert.Equalf(t, http.StatusForbidden, v.MockStatus, "中间件流程执行不是预期的过错;expect:%d,actual:%d", http.StatusForbidden, v.MockStatus)
			},
		},
	}
	for _, tt := range tests {

		//初始化测试用例参数
		tt.initCaseParam(tt)
		//获取中间件
		handler := middleware.BlackList()
		//调用中间件
		handler(tt.ctx)
		//断言结果
		tt.assertResult(t, tt.ctx)
	}
}
