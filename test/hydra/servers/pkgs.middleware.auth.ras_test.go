package servers

import (
	"testing"

	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:taoshouyin
//time:2020-11-12
//desc:测试RAS授权
func TestRASAuth(t *testing.T) {
	reqUrl := "/single/hydra/newversion/md5/auth@authserver.sas_debug"
	global.Def.RegistryAddr = "zk://192.168.0.101"
	type testCase struct {
		name        string
		opts        []ras.Option
		isSet       bool
		queryMap    map[string]interface{}
		wantStatus  int
		wantSpecial string
	}

	tests := []*testCase{
		{name: "rasAuth-配置不存在", isSet: false, wantStatus: 200, wantSpecial: "", queryMap: nil, opts: []ras.Option{}},
		{name: "rasAuth-配置存在数据为空", isSet: true, wantStatus: 200, wantSpecial: "", queryMap: nil, opts: []ras.Option{}},
		{name: "rasAuth-配置存在,未启用", isSet: true, wantStatus: 200, wantSpecial: "", queryMap: nil,
			opts: []ras.Option{ras.WithDisable(),
				ras.WithAuths(ras.New(reqUrl, ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-未在验证范围内", isSet: true, wantStatus: 200, wantSpecial: "", queryMap: nil,
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithRequest("/rasauth/test1"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-请求数据为空", isSet: true, wantStatus: 500, wantSpecial: "ras", queryMap: nil,
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-没有签名字段", isSet: true, wantStatus: 403, wantSpecial: "ras", queryMap: map[string]interface{}{"fied1": "111111"},
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-错误的请求路径", isSet: true, wantStatus: 403, wantSpecial: "ras", queryMap: map[string]interface{}{"fied1": "111111"},
			opts: []ras.Option{ras.WithAuths(ras.New("/error/path/auth@authserver.sas_debug",
				ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-错误的服务器地址", isSet: true, wantStatus: 403, wantSpecial: "ras", queryMap: map[string]interface{}{"fied1": "111111"},
			opts: []ras.Option{ras.WithAuths(ras.New("/single/hydra/newversion/md5/auth@authserver.sas_err",
				ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-错误的签名字段", isSet: true, wantStatus: 403, wantSpecial: "ras", queryMap: map[string]interface{}{"fied1": "111111", "sign": "errordata"},
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithConnect(ras.WithConnectChar("=", "&"), ras.WithConnectSortByData(), ras.WithSecretConnect(ras.WithSecretHeadMode(""))),
				ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-签名认证成功", isSet: true, wantStatus: 200, wantSpecial: "ras", queryMap: map[string]interface{}{"fied1": "111111", "sign": getRasSign(map[string]interface{}{"fied1": "111111"})},
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithConnect(ras.WithConnectChar("=", "&"), ras.WithConnectSortByData(), ras.WithSecretConnect(ras.WithSecretHeadMode(""))),
				ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
		{name: "rasAuth-启用-空数据签名认证成功", isSet: true, wantStatus: 200, wantSpecial: "ras", queryMap: map[string]interface{}{"sign": getRasSign(map[string]interface{}{})},
			opts: []ras.Option{ras.WithAuths(ras.New(reqUrl, ras.WithConnect(ras.WithConnectChar("=", "&"), ras.WithConnectSortByData(), ras.WithSecretConnect(ras.WithSecretHeadMode(""))),
				ras.WithRequest("/rasauth/test"), ras.WithSignAlias("sign")))}},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		confB := mockConf.GetAPI()
		if tt.isSet {
			confB.Ras(tt.opts...)
		}
		mockConf.Vars()["rpc"] = map[string]interface{}{
			"rpc": map[string]interface{}{
				"a": "1",
			},
		}

		serverConf := mockConf.GetAPIConf()

		app.Cache.Save(serverConf)

		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockRequest: &mocks.MockRequest{
				MockPath:     &mocks.MockPath{MockRequestPath: "/rasauth/test"},
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
		gotStatus, _ := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}

func getRasSign(param map[string]interface{}) string {
	secret := "JHfReB4Mn38z6V3npU3AKIvYqXI8b3VT"
	values := net.NewValues()
	for key, v := range param {
		values.Set(key, v.(string))
	}
	values.Sort()
	raw := values.Join("=", "&")
	return md5.Encrypt(secret + raw)
}
