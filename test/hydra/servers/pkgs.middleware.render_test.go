package servers

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

var scriptOK = `
request := import("request")
response := import("response")
text := import("text")
types :=import("types")

rc:="<response><code>{@status}</code><msg>{@content}</msg></response>"

getContent := func(){

	input:={status:response.getStatus(),content:response.getContent()}
	if text.has_prefix(request.getPath(),"/tx/request"){
		return [333,types.translate(rc,input),"application/xml"]
	}
	if text.has_prefix(request.getPath(),"/tx/query"){
		return [444,"<json>","application/json"]
	}
	return [204,response.getContent()]
}

render := getContent()`

var scriptERR = `
getContent := func(){
	return "error"
}

render := getContent()`

var scriptERR1 = `
getContent := func(){
	return [error]
}

render := getContent()`

var scriptERR2 = `
response := import("response")

getContent := func(){
	return response.getContent1()
}

render := getContent()`

func TestRender(t *testing.T) {

	tests := []struct {
		name            string
		script          string
		requestURL      string
		isSet           bool
		responseStatus  int
		responseCType   string
		responseContent string
		wantStatus      int
		wantContent     string
		wantSpecial     string
		wantContentType string
	}{
		{name: "render 未设置节点", isSet: false, script: "", requestURL: "/tx/request", responseStatus: 200, responseContent: "success", responseCType: "dddd", wantStatus: 200, wantContent: "success", wantContentType: "dddd", wantSpecial: ""},
		{name: "render 设置错误的节点,编译报错", isSet: true, script: scriptERR1, requestURL: "/tx/request", responseStatus: 200, wantStatus: 510, responseContent: "success", responseCType: "dddd", wantContent: "render脚本错误", wantContentType: "dddd", wantSpecial: ""},
		{name: "render 设置正确节点,运行报错", isSet: true, script: scriptERR2, requestURL: "/tx/ssss", responseStatus: 200, wantStatus: 200, responseContent: "success", responseCType: "dddd", wantContent: "success", wantContentType: "dddd", wantSpecial: ""},
		{name: "render 设置返回一个参数的节点", isSet: true, script: scriptERR, requestURL: "/tx/request", responseStatus: 200, responseContent: "success", responseCType: "dddd", wantStatus: 200, wantContent: "success", wantContentType: "dddd", wantSpecial: ""},
		{name: "render 设置正确节点,返回xml数据", isSet: true, script: scriptOK, requestURL: "/tx/request", responseStatus: 200, wantStatus: 333, responseContent: "success", responseCType: "application/xml", wantContent: "<response><code>200</code><msg>success</msg></response>", wantContentType: "application/xml", wantSpecial: "render"},
		{name: "render 设置正确节点,返回json数据", isSet: true, script: scriptOK, requestURL: "/tx/query", responseStatus: 200, wantStatus: 444, responseContent: "success", responseCType: "application/json", wantContent: "<json>", wantContentType: "application/json", wantSpecial: "render"},
		{name: "render 设置正确节点,返回两个参数数据", isSet: true, script: scriptOK, requestURL: "/tx/ssss", responseStatus: 200, wantStatus: 204, responseContent: "success", responseCType: "", wantContent: "success", wantContentType: "", wantSpecial: "render"},
	}

	for _, tt := range tests {
		conf := mocks.NewConfBy("middleware_render_test", "render")
		confN := conf.API(":8080")
		if tt.isSet {
			confN.Render(tt.script)
		}
		//初始化测试用例参数
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{},
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockRequestPath: tt.requestURL}},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus, MockContent: tt.responseContent, MockHeader: map[string][]string{"Content-Type": []string{tt.responseCType}}},
			MockAPPConf:  conf.GetAPIConf(),
		}

		//调用中间件
		gid := global.GetGoroutineID()
		context.Del(gid)
		context.Cache(ctx)
		handler := middleware.Render()
		handler(ctx)

		gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name)
		assert.Equalf(t, true, strings.Contains(gotContent, tt.wantContent), tt.name)
		gotHeaders := ctx.Response().GetHeaders()
		assert.Equalf(t, tt.wantContentType, gotHeaders["Content-Type"][0], tt.name)

		if tt.wantSpecial != "" {
			gotSpecial := ctx.Response().GetSpecials()
			assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
		}
	}
}

func BenchmarkRender(b *testing.B) {
	conf := mocks.NewConfBy("middleware_render1_test", "render1")
	confN := conf.API(":8080")
	confN.Render(scriptOK)
	//初始化测试用例参数

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gid := global.GetGoroutineID()
		context.Del(gid)
		ctx := &mocks.MiddleContext{
			MockUser:     &mocks.MockUser{},
			MockRequest:  &mocks.MockRequest{MockPath: &mocks.MockPath{MockRequestPath: "/tx/request"}},
			MockResponse: &mocks.MockResponse{MockStatus: 200, MockContent: "success", MockHeader: map[string][]string{"Content-Type": []string{"application/json"}}},
			MockAPPConf:  conf.GetAPIConf(),
		}
		context.Cache(ctx)
		handler := middleware.Render()
		handler(ctx)

		gotStatus, gotContent, _ := ctx.Response().GetFinalResponse()
		if gotStatus != 333 {
			b.Error("获取的数据有误")
			return
		}

		if gotContent != "<response><code>200</code><msg>success</msg></response>" {
			b.Error("获取的数据有误1")
			return
		}

		gotHeaders := ctx.Response().GetHeaders()
		if "application/json" != gotHeaders["Content-Type"][0] {
			b.Error("获取的数据有误2")
			return
		}
	}
}
