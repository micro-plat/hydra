package hydra

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
		{name: "限流-未启用-未配置", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "", opts: []limiter.Option{}},
		{name: "限流-未启用-Disable=true,但是rule=0", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable()}},
		{name: "限流-启用-不在限流配置内", requestPath: "/limiter-notin", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(),
				limiter.WithRuleList(&limiter.Rule{
					Path:     "/limiter",
					MaxAllow: 1,
					MaxWait:  0,
					Fallback: false,
					Resp:     &limiter.Resp{Status: 510, Content: "fallback"},
				}),
			},
		},
		{name: "限流-启用-在限流配置内-不延迟", requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "",
			opts: []limiter.Option{limiter.WithEnable(),
				limiter.WithRuleList(&limiter.Rule{
					Path:     "/limiter",
					MaxAllow: 1,
					MaxWait:  0,
					Fallback: false,
					Resp:     &limiter.Resp{Status: 510, Content: "fallback"},
				}),
			},
		},
		{name: "限流-启用-在限流配置内-延迟-降级", requestPath: "/limiter", wantStatus: 410, wantContent: "fallback", wantSpecial: "limit",
			opts: []limiter.Option{limiter.WithEnable(),
				limiter.WithRuleList(&limiter.Rule{
					Path:     "/limiter",
					MaxAllow: 0,
					MaxWait:  1,
					Fallback: false,
					Resp:     &limiter.Resp{Status: 410, Content: "fallback"},
				}),
			},
		},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		apiConf := mocks.NewConf()
		confB := apiConf.API(":51001")
		confB.Limit(tt.opts...)
		serverConf := apiConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockTFuncs: map[string]interface{}{},
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
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()

		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)

	}
}

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
		{name: "限流-启用-在限流配置内-延迟-不降级", count: 1, requestPath: "/limiter", wantStatus: 200, wantContent: "", wantSpecial: "limit"},
		{name: "限流-启用-在限流配置内-延迟-降级", count: 10, requestPath: "/limiter", wantStatus: 410, wantContent: "fallback", wantSpecial: "limit"},
	}

	global.Def.ServerTypes = []string{http.API}
	apiConf := mocks.NewConf()
	confB := apiConf.API(":51001")
	confB.Limit(opts...)
	serverConf := apiConf.GetAPIConf()
	//获取中间件
	handler := middleware.Limit()
	ctx := &mocks.MiddleContext{
		MockTFuncs: map[string]interface{}{},
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
			go func() {
				defer wg.Done()
				ctx := &mocks.MiddleContext{
					MockTFuncs: map[string]interface{}{},
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
				//断言结果
				gotStatus, gotContent := ctx.Response().GetFinalResponse()
				gotSpecial := ctx.Response().GetSpecials()
				assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
				assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
				assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
			}()
		}
		wg.Wait()
	}
}
