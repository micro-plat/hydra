package hydra

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:liujinyin
//time:2020-10-21 17:00
//desc:测试RAS授权
func TestRASAuth(t *testing.T) {

	type testCase struct {
		name        string
		opts        []ras.Option
		queryMap    map[string]interface{}
		router      *router.Router
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{
			name:        "远程服务验证-未启用-未配置",
			opts:        []ras.Option{},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name: "远程服务验证-未启用-Disable=true",
			opts: []ras.Option{
				ras.WithDisable(),
			},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name: "远程服务验证-启用-不在远程服务验证配置内",
			opts: []ras.Option{
				ras.WithEnable(),
				ras.WithAuths(&ras.Auth{
					//远程验证服务名
					Service: "ras",
				}),
			},
			router: &router.Router{
				Service: "/ras-notin",
			},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name: "远程服务验证-启用-在远程服务验证配置内",
			opts: []ras.Option{
				ras.WithEnable(),
				ras.WithAuths(&ras.Auth{
					//远程验证服务名
					Service: "/ras",
				}),
			},
			router: &router.Router{
				Service: "/ras",
			},
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		fmt.Println("---------------------------", tt.name)

		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		//mockConf.Service.API.Add()
		//初始化测试用例参数
		mockConf.GetAPI().Ras(tt.opts...)
		serverConf := mockConf.GetAPIConf()

		ctx := &mocks.MiddleContext{
			MockTFuncs: map[string]interface{}{},
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockRouter: tt.router,
				},
				MockQueryMap: tt.queryMap,
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}

		//获取中间件
		handler := middleware.RASAuth()

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
