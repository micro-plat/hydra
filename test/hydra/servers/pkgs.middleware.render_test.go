package servers

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name            string
		requestURL      string
		requstID        string
		method          string
		options         []render.Option
		responseStatus  int
		wantStatus      int
		wantContent     string
		wantSpecial     string
		wantContentType string
	}{
		{name: "render TMPLT配置格式错误返回510", options: []render.Option{render.WithTmplt("/path1", "success")}, responseStatus: 200, wantStatus: http.StatusNotExtended, wantContent: "render Tmplt配置数据有误:status: non zero value required"},
		{name: "render配置Disable为true", options: []render.Option{render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1"))}, responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "render获取渲染响应结果路径不匹配", requestURL: "/URL", method: "GET", options: []render.Option{render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1"))}, responseStatus: 200, wantStatus: 200, wantContent: ""},
		{name: "render获取渲染响应出错", requstID: "123456", requestURL: "/path1", method: "GET", options: []render.Option{render.WithTmplt("/path1", "success", render.WithStatus("0"), render.WithContentType("tpltm1"))}, responseStatus: 200,
			wantStatus: 200, wantContent: ""},
		{name: "render获取渲染响应成功", requestURL: "/path1", method: "GET", options: []render.Option{render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1"))}, responseStatus: 200,
			wantSpecial: "render", wantStatus: 500, wantContent: "success", wantContentType: "tpltm1"},
	}

	for _, tt := range tests {
		conf := mocks.NewConfBy("middleware_test", "render")
		confN := conf.API(":8080")
		confN.Render(tt.options...)
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{MockRequestID: tt.requstID},
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockRequestPath: tt.requestURL, MockURL: tt.requestURL, MockMethod: tt.method}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus, MockHeader: map[string][]string{}},
			MockAPPConf:  conf.GetAPIConf(),
		}

		//调用中间件
		handler := middleware.Render()
		handler(ctx)

		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name)

		if tt.wantSpecial != "" {
			gotSpecial := ctx.Response().GetSpecials()
			assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
			gotHeaders := ctx.Response().GetHeaders()
			assert.Equalf(t, tt.wantContentType, gotHeaders["Content-Type"][0], tt.name)
		}
	}
}
