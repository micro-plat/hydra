package servers

import (
	"testing"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name        string
		requestURL  string
		method      string
		status      int
		content     string
		wantStatus  int
		wantContent string
		wantSpecial string
	}{
		{name: "1.1 OPTIONS-请求返回200", method: "OPTIONS", status: 200, content: "result", wantStatus: 200, wantContent: "", wantSpecial: "opt"},
		{name: "1.2 OPTIONS-请求返回400", method: "OPTIONS", status: 400, content: "result1", wantStatus: 200, wantContent: "", wantSpecial: "opt"},
		{name: "1.3 OPTIONS-请求返回500", method: "OPTIONS", status: 500, content: "result2", wantStatus: 200, wantContent: "", wantSpecial: "opt"},

		{name: "2.1 GET-请求返回200", method: "GET", status: 200, content: "result", wantStatus: 200, wantContent: "result", wantSpecial: ""},
		{name: "2.2 GET-请求返回400", method: "GET", status: 400, content: "result1", wantStatus: 400, wantContent: "result1", wantSpecial: ""},
		{name: "2.3 GET-请求返回500", method: "GET", status: 500, content: "result2", wantStatus: 500, wantContent: "result2", wantSpecial: ""},

		{name: "3.1 POST-请求返回200", method: "POST", status: 200, content: "result", wantStatus: 200, wantContent: "result", wantSpecial: ""},
		{name: "3.2 POST-请求返回400", method: "POST", status: 400, content: "result1", wantStatus: 400, wantContent: "result1", wantSpecial: ""},
		{name: "3.3 POST-请求返回500", method: "POST", status: 500, content: "result2", wantStatus: 500, wantContent: "result2", wantSpecial: ""},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockURL: tt.requestURL, MockMethod: tt.method}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.status, MockContent: tt.content},
			MockAPPConf:  mocks.NewConfBy("middleware_options_test", "options").GetAPIConf(),
		}

		//调用中间件
		handler := middleware.Options()
		handler(ctx)

		gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name)

		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
	}
}
