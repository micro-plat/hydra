package servers

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
)

//author:taoshouyin
//time:2020-11-11
//desc:测试apikey验证中间件逻辑
//由于密钥验证方法已经经过单元测试,现在之需要对中间件处理逻辑进行测试,所以以下用例只使用md5进行测试
func TestAuthAPIKey(t *testing.T) {
	secret := "54545454"
	type testCase struct {
		name        string
		apikeyOpts  []apikey.Option
		params      map[string]string
		isSet       bool
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "apikey-未配置", isSet: false, params: map[string]string{"f1": "21"}, wantStatus: 200, wantContent: "", wantSpecial: "", apikeyOpts: []apikey.Option{}},
		{name: "apikey-配置未启动", isSet: true, params: map[string]string{"f1": "21"}, wantStatus: 200, wantContent: "", wantSpecial: "", apikeyOpts: []apikey.Option{apikey.WithDisable(), apikey.WithSecret(secret)}},
		{name: "apikey-配置启动-配置错误", isSet: true, params: map[string]string{"f1": "21"}, wantStatus: 510, wantContent: "apikey配置数据有误", wantSpecial: "",
			apikeyOpts: []apikey.Option{apikey.WithSecret("错误密钥")}},
		{name: "apikey-配置启动-路径被排除", isSet: true, params: map[string]string{"f1": "21"}, wantStatus: 200, wantContent: "", wantSpecial: "",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test")}},
		{name: "apikey-配置启动-缺少sign字段", isSet: true, params: map[string]string{"f1": "21", "timestamp": "1212121221"}, wantStatus: 401, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
		{name: "apikey-配置启动-缺少timestamp字段", isSet: true, params: map[string]string{"f1": "21", "sign": getSign(map[string]string{"f1": "21"}, secret)}, wantStatus: 401, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
		{name: "apikey-配置启动-验证失败", isSet: true, params: map[string]string{"f1": "21", "sign": "4444444", "timestamp": "5421515"}, wantStatus: 403, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
		{name: "apikey-配置启动-验证通过-utf8", isSet: true, params: map[string]string{"f1": "21", "timestamp": "5421515", "sign": getSign(map[string]string{"f1": "21", "timestamp": "5421515"}, secret)}, wantStatus: 200, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
		{name: "apikey-配置启动-验证通过-gbk", isSet: true, params: map[string]string{"f1": "21", "timestamp": "5421515", "sign": getSign(map[string]string{"f1": "21", "timestamp": "5421515"}, secret)}, wantStatus: 200, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
		{name: "apikey-配置启动-验证通过-gb2312", isSet: true, params: map[string]string{"f1": "21", "timestamp": "5421515", "sign": getSign(map[string]string{"f1": "21", "timestamp": "5421515"}, secret)}, wantStatus: 200, wantContent: "", wantSpecial: "apikey",
			apikeyOpts: []apikey.Option{apikey.WithSecret(secret), apikey.WithMD5Mode(), apikey.WithExcludes("/apikey/test1")}},
	}
	for _, tt := range tests {
		mockConf := mocks.NewConfBy("middleware_apikey_test", "apikey")
		//初始化测试用例参数
		confB := mockConf.GetAPI()
		if tt.isSet {
			confB.APIKEY("", tt.apikeyOpts...)
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockRequest: &mocks.MockRequest{
				MockParamMap: tt.params,
				MockPath: &mocks.MockPath{
					MockRequestPath: "/apikey/test",
				},
			},
			MockAPPConf: serverConf,
		}

		//获取中间件
		handler := middleware.APIKeyAuth()
		//调用中间件
		handler(ctx)
		//断言结果
		gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, true, strings.Contains(gotContent, tt.wantContent), tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}

func getSign(param map[string]string, secret string) string {
	values := net.NewValues()
	for key, v := range param {
		values.Set(key, v)
	}
	values.Sort()
	raw := values.Join("", "")
	return md5.Encrypt(raw + secret)
}
