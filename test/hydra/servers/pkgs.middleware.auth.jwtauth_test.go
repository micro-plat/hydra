package hydra

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	octx "github.com/micro-plat/hydra/context/ctx"

	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:taoshouyin
//time:2020-11-11
//desc:测试basic验证中间件逻辑
func TestJWTAuth(t *testing.T) {
	type testCase struct {
		name        string
		jwtOpts     []jwt.Option
		requestPath string
		isSet       bool
		wantStatus  int
		wantSpecial string
	}

	tests := []*testCase{
		{name: "jwt-未配置", isSet: false, wantStatus: 200, wantSpecial: "", jwtOpts: []jwt.Option{}},
		{name: "jwt-配置未启动", isSet: true, wantStatus: 200, wantSpecial: "", jwtOpts: []jwt.Option{jwt.WithDisable()}},
		{name: "jwt-配置启动-路由不存在", isSet: true, wantStatus: 510, wantSpecial: "",
			jwtOpts: []jwt.Option{}},
		{name: "jwt-配置启动-路由存在,被排除", isSet: true, wantStatus: 200, wantSpecial: "",
			jwtOpts: []jwt.Option{jwt.WithExcludes("/jwt/test")}},
		{name: "jwt-配置启动-不排除,认证信息为空", isSet: true, wantStatus: 200, wantSpecial: "",
			jwtOpts: []jwt.Option{jwt.WithExcludes("/jwt/test1")}},
		{name: "jwt-配置启动-不排除,认证失败", isSet: true, wantStatus: 401, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithExcludes("/jwt/test1")}},
		{name: "jwt-配置启动-不排除,认证成功", isSet: true, wantStatus: 200, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithExcludes("/jwt/test1")}},
	}

	for _, tt := range tests {
		mockConf := mocks.NewConf()
		//初始化测试用例参数
		confB := mockConf.GetAPI()
		if tt.isSet {
			confB.Jwt(tt.jwtOpts...)
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:     conf.NewMeta(),
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1", MockAuth: &octx.Auth{}},
			MockResponse: &mocks.MockResponse{MockStatus: 200, MockHeader: map[string][]string{}},
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockHeader:      nil,
					MockRequestPath: tt.requestPath,
				},
			},
			MockServerConf: serverConf,
		}

		//获取中间件
		handler := middleware.BasicAuth()
		//调用中间件
		handler(ctx)
		//断言结果
		gotStatus, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}
