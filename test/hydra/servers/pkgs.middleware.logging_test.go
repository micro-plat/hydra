package servers

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestLogging(t *testing.T) {

	tests := []struct {
		name           string
		clientIP       string
		requstID       string
		requestURL     string
		method         string
		responseStatus int
		wantReq        string
		wantRsp        string
	}{
		{name: "1.1 GET请求-返回101", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 101, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 101  "},
		{name: "1.2 GET请求-返回200", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 200, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 200  "},
		{name: "1.3 GET请求-返回300", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 300, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 300  "},
		{name: "1.4 GET请求返回400", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 400, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 400  "},
		{name: "1.5 GET请求返回401", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 401, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 401  "},
		{name: "1.6 GET请求返回510", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "GET", responseStatus: 510, wantReq: "api.request: GET /URL from 127.0.0.1", wantRsp: "api.response: GET /URL 510  "},

		{name: "2.1 POST请求-返回101", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 101, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 101  "},
		{name: "2.2 POST请求-返回200", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 200, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 200  "},
		{name: "2.3 POST请求-返回300", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 300, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 300  "},
		{name: "2.4 POST请求返回400", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 400, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 400  "},
		{name: "2.5 POST请求返回401", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 401, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 401  "},
		{name: "2.6 POST请求返回510", clientIP: "127.0.0.1", requstID: "06c6fb24c", requestURL: "/URL", method: "POST", responseStatus: 510, wantReq: "api.request: POST /URL from 127.0.0.1", wantRsp: "api.response: POST /URL 510  "},
	}

	for _, tt := range tests {
		time.Sleep(time.Second)
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{MockClientIP: tt.clientIP, MockRequestID: tt.requstID},
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockURL: tt.requestURL, MockMethod: tt.method}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus},
			MockAPPConf:  mocks.NewConfBy("middleware_logging_test", "logging").GetAPIConf(),
		}

		//构建的新的os.Stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//调用中间件
		handler := middleware.Logging()
		handler(ctx)
		time.Sleep(time.Second * 1)

		//获取输出
		w.Close()
		out, err := ioutil.ReadAll(r)
		assert.Equalf(t, false, err != nil, tt.name)

		//还原os.Stdout
		os.Stdout = rescueStdout

		gotStatus, _, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.responseStatus, gotStatus, tt.name)
		assert.Equalf(t, true, strings.Contains(string(out), tt.wantReq), tt.name)
		assert.Equalf(t, true, strings.Contains(string(out), tt.wantRsp), tt.name)
	}
}
