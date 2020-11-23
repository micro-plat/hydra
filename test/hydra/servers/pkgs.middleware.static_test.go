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
		requestPath    string
		method         string
		responseStatus int
		wantStatus     int
		wantContent    string
		wantSpecial    string
	}{
		{name: "1.配置Static出错", opts: []static.Option{static.WithArchive("pkgs.middleware.middle.test.txt")}, responseStatus: 200, wantStatus: 510, wantContent: "pkgs.middleware.middle.test.txt获取失败:指定的文件不是归档文件:pkgs.middleware.middle.test.txt"},
		{name: "2.配置Static.disable为true", opts: []static.Option{static.WithDisable()}, responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "3.请求不按静态文件处理", requestPath: "/path", method: "OPTIONS", responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "4.请求是静态文件,但不存在", requestPath: "/robots.txt", method: "GET", responseStatus: 200, wantStatus: 404, wantContent: "找不到文件:src/robots.txt stat src/robots.txt: no such file or directory", wantSpecial: "static"},
		{name: "5.请求是静态文件夹", opts: []static.Option{static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test", method: "GET", responseStatus: 200, wantStatus: 404, wantContent: "找不到文件:pkgs.middleware.static_test", wantSpecial: "static"},
		{name: "6.请求是正确的静态文件", opts: []static.Option{static.WithRoot("./")}, requestPath: "pkgs.middleware.static_test.txt", method: "GET", wantStatus: 200, wantContent: "pkgs.middleware.static_test.txt", wantSpecial: "static"},
		//{name: "请求路径是静态文件夹", opts: []static.Option{static.WithRoot("./pkgs.middleware.static_test"), static.WithPrefix("pkgs."), static.AppendExts(".js")}, requestPath: "pkgs.test1.gz", method: "GET", responseStatus: 200, wantStatus: 200, wantContent: "", wantSpecial: "static"},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		c := mocks.NewConfBy("middleware_static_test", "static")
		apiConf := c.API(":9090")
		apiConf.Static(tt.opts...)
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
