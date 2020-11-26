package servers

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestStatic(t *testing.T) {
	tests := []struct {
		name           string
		opts           []static.Option
		isBool         bool
		requestPath    string
		method         string
		responseStatus int
		wantStatus     int
		wantContent    string
		wantSpecial    string
	}{

		{name: "1.1 static-配置不存在", isBool: false, opts: []static.Option{}, responseStatus: 200, wantStatus: 200, wantContent: ""},

		{name: "2.1 static-配置存在-不启用-无数据", isBool: false, opts: []static.Option{static.WithDisable()}, responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "2.2 static-配置存在-不启用-配置Static出错", isBool: true, opts: []static.Option{static.WithDisable(), static.WithArchive("pkgs.middleware.middle.test.txt")}, responseStatus: 200, wantStatus: 510, wantContent: "pkgs.middleware.middle.test.txt获取失败:指定的文件不是归档文件:pkgs.middleware.middle.test.txt"},
		{name: "2.3 static-配置存在-不启用-请求不按静态文件处理", isBool: true, opts: []static.Option{static.WithDisable()}, requestPath: "/path", method: "OPTIONS", responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "2.4 static-配置存在-不启用-请求是静态文件,但不存在", isBool: true, opts: []static.Option{static.WithDisable()}, requestPath: "/robots.txt", method: "GET", responseStatus: 200, wantStatus: 200, wantContent: "", wantSpecial: ""},
		{name: "2.5 static-配置存在-不启用-请求是静态文件夹", isBool: true, opts: []static.Option{static.WithDisable(), static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test", method: "GET", responseStatus: 200, wantStatus: 200, wantContent: "", wantSpecial: ""},
		{name: "2.6 static-配置存在-不启用-请求是正确的静态文件", isBool: true, opts: []static.Option{static.WithDisable(), static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test.txt", method: "GET", responseStatus: 200, wantStatus: 200, wantContent: "", wantSpecial: ""},

		{name: "3.1 static-配置存在-启用-无数据", isBool: true, opts: []static.Option{static.WithEnable()}, responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "3.2 static-配置存在-启用-配置Static出错", isBool: true, opts: []static.Option{static.WithArchive("pkgs.middleware.middle.test.txt")}, responseStatus: 200, wantStatus: 510, wantContent: "pkgs.middleware.middle.test.txt获取失败:指定的文件不是归档文件:pkgs.middleware.middle.test.txt"},
		{name: "3.3 static-配置存在-启用-请求不按静态文件处理", isBool: true, requestPath: "/path", method: "OPTIONS", responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "3.4 static-配置存在-启用-请求是静态文件,但不存在", isBool: true, requestPath: "/robots.txt", method: "GET", responseStatus: 200, wantStatus: 404, wantContent: "找不到文件:src/robots.txt stat src/robots.txt: no such file or directory", wantSpecial: "static"},
		{name: "3.5 static-配置存在-启用-请求是静态文件夹", isBool: true, opts: []static.Option{static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test", method: "GET", responseStatus: 200, wantStatus: 404, wantContent: "找不到文件:pkgs.middleware.static_test", wantSpecial: "static"},
		{name: "3.6 static-配置存在-启用-请求是正确的静态文件", isBool: true, opts: []static.Option{static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test.txt", method: "GET", wantStatus: 200, wantContent: "pkgs.middleware.static_test.txt", wantSpecial: "static"},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		c := mocks.NewConfBy("middleware_static_test", "static")
		apiConf := c.API(":9090")
		if tt.isBool {
			apiConf.Static(tt.opts...)
		}
		ctx := &mocks.MiddleContext{
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockRequestPath: tt.requestPath, MockMethod: tt.method}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus},
			MockAPPConf:  c.GetAPIConf(),
		}

		//调用中间件
		handler := middleware.Static()
		handler(ctx)

		gotStatus, getContent, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name)
		assert.Equalf(t, tt.wantContent, getContent, tt.name)

		if tt.wantSpecial != "" {
			gotSpecial := ctx.Response().GetSpecials()
			assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
		}
	}
}
