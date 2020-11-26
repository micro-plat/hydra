package servers

import (
	"sync"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:liujinyin
//time:2020-10-21 10:00
//desc:测试限流中间件逻辑
func TestLimit(t *testing.T) {

	type testCase struct {
		name        string
		requestPath string
		opts        []limiter.Option
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "1.1 限流-配置不存在", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{}},
		{name: "1.2 限流-配置存在-未启用-无数据", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{}},
		{name: "1.3 限流-配置存在-未启用-不在限流配置内", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{}},
		{name: "1.4 限流-配置存在-未启用-在限流配置内", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{}},

		{name: "2.1 限流-配置存在-启用-无数据", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{limiter.WithEnable()}},
		{name: "2.2 限流-配置存在-启用-不在限流配置内-不延迟", requestPath: "/limiter-notin", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 1, MaxWait: 0, Fallback: false, Resp: &limiter.Resp{Status: 510, Content: "fallback"}})}},
		{name: "2.3 限流-配置存在-启用-不在限流配置内-延迟-不降级", requestPath: "/limiter-notin", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 0, MaxWait: 1, Fallback: false, Resp: &limiter.Resp{Status: 510, Content: "fallback"}})}},
		{name: "2.4 限流-配置存在-启用-不在限流配置内-延迟-降级", requestPath: "/limiter-notin", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 0, MaxWait: 1, Fallback: true, Resp: &limiter.Resp{Status: 510, Content: "fallback"}})}},

		{name: "3.1 限流-配置存在-启用-在限流配置内-不延迟", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 1, MaxWait: 0, Fallback: false, Resp: &limiter.Resp{Status: 510, Content: "fallback"}})}},
		{name: "3.2 限流-配置存在-启用-在限流配置内-延迟-降级", requestPath: "/limiter", wantStatus: 410, wantContent: "fallback", wantSpecial: "limit",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 0, MaxWait: 1, Fallback: false, Resp: &limiter.Resp{Status: 410, Content: "fallback"}})}},
		{name: "3.3 限流-配置存在-启用-在限流配置内-延迟-不降级", requestPath: "/limiter", wantStatus: 410, wantContent: "fallback", wantSpecial: "limit",
			opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(&limiter.Rule{Path: "/limiter", MaxAllow: 0, MaxWait: 1, Fallback: true, Resp: &limiter.Resp{Status: 410, Content: "fallback"}})},
		},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		apiConf := mocks.NewConfBy("middleware_limit_test1", "limiter1")
		confB := apiConf.API(":51001")
		confB.Limit(tt.opts...)
		serverConf := apiConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockRequestPath: tt.requestPath,
				},
			},
			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}

		//获取中间件
		handler := middleware.Limit()

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

//并发测试限制流量
func TestLimit1(t *testing.T) {
	opts := []limiter.Option{limiter.WithEnable(),
		limiter.WithRuleList(&limiter.Rule{
			Path:     "/limiter",
			MaxAllow: 1,
			MaxWait:  1,
			Fallback: false,
			Resp:     &limiter.Resp{Status: 410, Content: "fallback"},
		}),
	}

	type testCase struct {
		name        string
		count       int
		requestPath string
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "1. 限流-启用-在限流配置内-延迟-不降级", count: 1, requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "limit"},
		// {name: "2. 限流-启用-在限流配置内-延迟-降级", count: 10, requestPath: "/limiter", wantStatus: 410, wantContent: "fallback", wantSpecial: "limit"},
	}

	global.Def.ServerTypes = []string{http.API}
	apiConf := mocks.NewConfBy("middleware_limit_test", "limiter")
	confB := apiConf.API(":51001")
	confB.Limit(opts...)
	serverConf := apiConf.GetAPIConf()
	//获取中间件
	handler := middleware.Limit()
	ctx := &mocks.MiddleContext{
		MockRequest: &mocks.MockRequest{
			MockPath: &mocks.MockPath{
				MockRequestPath: "/limiter",
			},
		},
		MockResponse: &mocks.MockResponse{MockStatus: 200},
		MockAPPConf:  serverConf,
	}
	//调用中间件
	handler(ctx)

	for _, tt := range tests {

		var wg sync.WaitGroup
		for j := 0; j < tt.count; j++ {

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ctx := &mocks.MiddleContext{
					MockRequest: &mocks.MockRequest{
						MockPath: &mocks.MockPath{
							MockRequestPath: tt.requestPath,
						},
					},
					MockResponse: &mocks.MockResponse{MockStatus: 200},
					MockAPPConf:  serverConf,
				}
				//调用中间件
				handler(ctx)
				// if i == 0 || i == tt.count-1 {
				// 	return
				// }

				//断言结果
				gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
				gotSpecial := ctx.Response().GetSpecials()
				assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
				assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
				assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
			}(j)
		}
		wg.Wait()
	}
}
