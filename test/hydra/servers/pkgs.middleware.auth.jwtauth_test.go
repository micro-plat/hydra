package hydra

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	octx "github.com/micro-plat/hydra/context/ctx"

	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	wjwt "github.com/micro-plat/lib4go/security/jwt"
	"github.com/micro-plat/lib4go/utility"
)

//author:taoshouyin
//time:2020-11-11
//desc:测试jwt验证中间件逻辑
func TestJWTAuth(t *testing.T) {
	secert := utility.GetGUID()
	requestPath := "/jwt/test"
	// routerObj := router.NewRouter(requestPath, "service", []string{"GET"})
	type testCase struct {
		name    string
		jwtOpts []jwt.Option
		// router      *router.Router
		token       string
		isSource    string //cookie/header
		authURL     string
		isSet       bool
		isSucc      bool
		wantStatus  int
		wantSpecial string
	}
	data := map[string]interface{}{"sdsd": "sdfd", "3ddfs": "gggggg"}
	rawData, _ := wjwt.Encrypt(secert, jwt.ModeHS512, data, 86400)

	tests := []*testCase{
		{name: "jwt-未配置", isSet: false, wantStatus: 200, wantSpecial: "", jwtOpts: []jwt.Option{}},
		{name: "jwt-配置未启动", isSet: true, wantStatus: 200, wantSpecial: "", jwtOpts: []jwt.Option{jwt.WithDisable()}},
		{name: "jwt-配置启动-被排除", isSet: true, wantStatus: 200, wantSpecial: "jwt", jwtOpts: []jwt.Option{jwt.WithExcludes("/jwt/test")}},
		{name: "jwt-配置启动-token不存在", isSet: true, wantStatus: 401, wantSpecial: "jwt", jwtOpts: []jwt.Option{}},
		{name: "jwt-配置启动-token在header中,失败", isSet: true, isSource: "header", token: "errorToken", wantStatus: 403, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithHeader(), jwt.WithExcludes("/jwt/test1")}},
		{name: "jwt-配置启动-token在cookie中,失败authurl不为空", isSet: true, authURL: "www.baidu.com", isSource: "cookie", token: "errorToken", wantStatus: 302, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithCookie(), jwt.WithAuthURL("www.baidu.com"), jwt.WithExcludes("/jwt/test1")}},
		{name: "jwt-配置启动-token在header中,成功", isSucc: true, isSet: true, isSource: "header", token: rawData, wantStatus: 200, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithHeader(), jwt.WithSecret(secert), jwt.WithExcludes("/jwt/test1")}},
		{name: "jwt-配置启动-token在cookie中,成功authurl不为空", isSucc: true, isSet: true, authURL: "www.baidu.com", isSource: "cookie", token: rawData, wantStatus: 200, wantSpecial: "jwt",
			jwtOpts: []jwt.Option{jwt.WithCookie(), jwt.WithSecret(secert), jwt.WithAuthURL("www.baidu.com"), jwt.WithExcludes("/jwt/test1")}},
	}

	for _, tt := range tests {
		mockConf := mocks.NewConf()
		//初始化测试用例参数
		confB := mockConf.GetAPI()
		if tt.isSet {
			confB.Jwt(tt.jwtOpts...)
		}
		headerMap := map[string][]string{}
		cookieMap := map[string]string{}
		if tt.isSource == "header" {
			headerMap[jwt.JWTName] = []string{tt.token}
		} else {
			cookieMap[jwt.JWTName] = tt.token
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:     conf.NewMeta(),
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1", MockAuth: &octx.Auth{}},
			MockResponse: &mocks.MockResponse{MockStatus: 200, MockHeader: map[string][]string{}},
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockHeader:  headerMap,
					MockCookies: cookieMap,
					// MockRouter:      tt.router,
					MockRequestPath: requestPath,
				},
			},
			MockAPPConf: serverConf,
		}

		//获取中间件
		handler := middleware.JwtAuth()
		//调用中间件
		handler(ctx)
		//断言结果
		gotStatus, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		if tt.isSucc {
			val := (ctx.User().Auth().Request()).(map[string]interface{})
			assert.Equalf(t, data, val, tt.name, data, val)
		} else if tt.authURL != "" {
			header := ctx.Response().GetHeaders()
			assert.Equalf(t, []string{tt.authURL}, header["Location"], tt.name, tt.authURL, header["Location"])
		}
	}
}
