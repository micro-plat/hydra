package servers

import (
	"testing"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name           string
		requestURL     string
		method         string
		responseStatus int
		wantStatus     int
		wantSpecial    string
	}{
		{name: "OPTIONS请求返回200", method: "OPTIONS", responseStatus: 400, wantStatus: 200, wantSpecial: "opt"},
		{name: "GET请求返回400", method: "GET", responseStatus: 400, wantStatus: 400, wantSpecial: ""},
		{name: "POST请求返回400", method: "POST", responseStatus: 400, wantStatus: 400, wantSpecial: ""},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockURL: tt.requestURL, MockMethod: tt.method}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus},
			MockAPPConf:  mocks.NewConf().GetAPIConf(),
		}

		//调用中间件
		handler := middleware.Options()
		handler(ctx)

		gotStatus, _, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name)

		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
	}
}
