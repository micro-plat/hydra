package servers

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

var golHandledErrFunc context.Handler = func(ctx context.IContext) interface{} {
	return errs.NewError(671, fmt.Errorf("全局后处理异常"))
}

var golHandledOKFunc context.Handler = func(ctx context.IContext) interface{} {
	return "golhandledsuccess"
}

var golHandlingErrFunc context.Handler = func(ctx context.IContext) interface{} {
	return errs.NewError(668, fmt.Errorf("全局预处理异常"))
}

var golHandlingOKFunc context.Handler = func(ctx context.IContext) interface{} {
	return "golhandlingsuccess"
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//无全局预处理和后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-限流,不降级", service: "/path/test1", isLimited: true, fallback: false, handleObj: &testObj1{}, wantStatus: 200, wantSpecial: ""}, //保留限流组件设置的status和content
		{name: "Handler-限流,降级,无函数", service: "/path/test2", isLimited: true, fallback: true, handleObj: &testObj2{}, wantStatus: 200, wantSpecial: "fallback"},
		{name: "Handler-限流,降级,有函数,异常输出", service: "/path/test3", isLimited: true, fallback: true, handleObj: &testObj3{}, wantStatus: 611, wantSpecial: "fallback"},
		{name: "Handler-限流,降级,有函数,异常输出1", service: "/path/test4", isLimited: true, fallback: true, handleObj: &testObj4{}, wantStatus: 500, wantSpecial: "fallback"},
		{name: "Handler-限流,降级,有函数,正常输出", service: "/path/test5", isLimited: true, fallback: true, handleObj: &testObj5{}, wantStatus: 200, wantSpecial: "fallback", wantContent: "fallsuccess"},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,无ALL后处理", service: "/path/test6", wantStatus: 666, wantSpecial: "", handleObj: &testObj6{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,无全局后处理", service: "/path/test9", wantStatus: 666, wantSpecial: "", handleObj: &testObj7{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,无全局后处理", service: "/path/test12", wantStatus: 666, wantSpecial: "", handleObj: &testObj8{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无ALL后处理", service: "/path/test15", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,无全局后处理", service: "/path/test18", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,无全局后处理", service: "/path/test21", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无ALL后处理", service: "/path/test42", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,无全局后处理", service: "/path/test45", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,无全局后处理", service: "/path/test48", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,无ALL后处理", service: "/path/test69", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj15{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,无全局后处理", service: "/path/test72", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj16{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,无全局后处理", service: "/path/test75", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj17{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无ALL后处理", service: "/path/test78", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无ALL后处理", service: "/path/test105", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,无全局后处理", service: "/path/test108", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,无全局后处理", service: "/path/test111", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,无全局后处理", service: "/path/test81", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: nil, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,无全局后处理", service: "/path/test84", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: nil, golHandlingFunc: nil},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//无全局预处理,有正确后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler1(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,无ALL预处理,无后处理,正确全局后处理", service: "/path/test7", wantStatus: 666, wantSpecial: "", handleObj: &testObj6{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,正确全局后处理", service: "/path/test10", wantStatus: 666, wantSpecial: "", handleObj: &testObj7{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,正确全局后处理", service: "/path/test13", wantStatus: 666, wantSpecial: "", handleObj: &testObj8{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无后处理,正确全局后处理", service: "/path/test16", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,正确全局后处理", service: "/path/test19", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,正确全局后处理", service: "/path/test22", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无后处理,正确全局后处理", service: "/path/test43", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,正确全局后处理", service: "/path/test46", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,正确全局后处理", service: "/path/test49", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,无后处理,正确全局后处理", service: "/path/test70", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj15{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,正确全局后处理", service: "/path/test73", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj16{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,正确全局后处理", service: "/path/test76", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj17{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无后处理,正确全局后处理", service: "/path/test106", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,正确全局后处理", service: "/path/test112", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,正确全局后处理", service: "/path/test109", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,正确全局后处理", service: "/path/test82", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无后处理,正确全局后处理", service: "/path/test79", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,正确全局后处理", service: "/path/test85", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: nil},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//无全局预处理,有错误后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler2(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,无ALL预处理,无后处理,错误全局后处理", service: "/path/test8", wantStatus: 666, wantSpecial: "", handleObj: &testObj6{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,错误全局后处理", service: "/path/test11", wantStatus: 666, wantSpecial: "", handleObj: &testObj7{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,错误全局后处理", service: "/path/test14", wantStatus: 666, wantSpecial: "", handleObj: &testObj8{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无后处理,错误全局后处理", service: "/path/test17", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,错误全局后处理", service: "/path/test20", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理", service: "/path/test23", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无后处理,错误全局后处理", service: "/path/test44", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,错误全局后处理", service: "/path/test47", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,错误全局后处理", service: "/path/test50", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,无后处理,错误全局后处理", service: "/path/test71", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj15{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,错误全局后处理", service: "/path/test74", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj16{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,错误全局后处理", service: "/path/test77", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj17{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理", service: "/path/test86", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,错误全局后处理", service: "/path/test83", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无后处理,错误全局后处理", service: "/path/test80", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无后处理,错误全局后处理", service: "/path/test107", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,错误全局后处理", service: "/path/test110", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
		{name: "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,错误全局后处理", service: "/path/test113", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: nil},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有正确全局预处理,无全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler3(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无ALL后处理", service: "/path/test24", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,无全局后处理", service: "/path/test27", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,无全局后处理", service: "/path/test30", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无ALL后处理", service: "/path/test51", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,无全局后处理", service: "/path/test54", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,无全局后处理", service: "/path/test57", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无ALL后处理", service: "/path/test87", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,无全局后处理", service: "/path/test90", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无ALL后处理", service: "/path/test114", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,无全局后处理", service: "/path/test117", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,无全局后处理", service: "/path/test120", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,无全局后处理", service: "/path/test93", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: nil, golHandlingFunc: golHandlingOKFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有正确全局预处理,有正确全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler4(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无后处理,正确全局后处理", service: "/path/test25", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,正确全局后处理", service: "/path/test31", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,正确全局后处理", service: "/path/test28", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无后处理,正确全局后处理", service: "/path/test52", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,正确全局后处理", service: "/path/test55", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,正确全局后处理", service: "/path/test58", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,正确全局后处理", service: "/path/test118", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无后处理,正确全局后处理", service: "/path/test115", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,正确全局后处理", service: "/path/test121", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无后处理,正确全局后处理", service: "/path/test88", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,正确全局后处理", service: "/path/test91", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,正确全局后处理", service: "/path/test94", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingOKFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有正确全局预处理,有错误全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler5(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无后处理,错误全局后处理", service: "/path/test26", wantStatus: 666, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,错误全局后处理", service: "/path/test29", wantStatus: 666, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,错误全局后处理", service: "/path/test32", wantStatus: 666, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无后处理,错误全局后处理", service: "/path/test53", wantStatus: 667, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,错误全局后处理", service: "/path/test56", wantStatus: 667, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,错误全局后处理", service: "/path/test59", wantStatus: 667, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无后处理,错误全局后处理", service: "/path/test116", wantStatus: 667, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,错误全局后处理", service: "/path/test119", wantStatus: 667, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,错误全局后处理", service: "/path/test122", wantStatus: 667, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无后处理,错误全局后处理", service: "/path/test89", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj18{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,错误全局后处理", service: "/path/test92", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj19{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,错误全局后处理", service: "/path/test95", wantStatus: 200, wantSpecial: "", wantContent: "mainsuccess", handleObj: &testObj20{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingOKFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有错误全局预处理,无全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler6(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{

		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无ALL后处理", service: "/path/test33", wantStatus: 668, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,无全局后处理", service: "/path/test36", wantStatus: 668, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,无全局后处理", service: "/path/test39", wantStatus: 668, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无ALL后处理", service: "/path/test60", wantStatus: 668, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,无全局后处理", service: "/path/test63", wantStatus: 668, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,无全局后处理", service: "/path/test66", wantStatus: 668, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无ALL后处理", service: "/path/test96", wantStatus: 668, wantSpecial: "", handleObj: &testObj18{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,无全局后处理", service: "/path/test99", wantStatus: 668, wantSpecial: "", handleObj: &testObj19{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,无全局后处理", service: "/path/test102", wantStatus: 668, wantSpecial: "", handleObj: &testObj20{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无ALL后处理", service: "/path/test123", wantStatus: 668, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,无全局后处理", service: "/path/test126", wantStatus: 668, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,无全局后处理", service: "/path/test129", wantStatus: 668, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: nil, golHandlingFunc: golHandlingErrFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有错误全局预处理,有正确全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler7(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无后处理,正确全局后处理", service: "/path/test34", wantStatus: 668, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,正确全局后处理", service: "/path/test37", wantStatus: 668, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,正确全局后处理", service: "/path/test40", wantStatus: 668, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无后处理,正确全局后处理", service: "/path/test61", wantStatus: 668, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,正确全局后处理", service: "/path/test64", wantStatus: 668, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,正确全局后处理", service: "/path/test67", wantStatus: 668, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无后处理,正确全局后处理", service: "/path/test97", wantStatus: 668, wantSpecial: "", handleObj: &testObj18{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,正确全局后处理", service: "/path/test100", wantStatus: 668, wantSpecial: "", handleObj: &testObj19{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,正确全局后处理", service: "/path/test103", wantStatus: 668, wantSpecial: "", handleObj: &testObj20{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无后处理,正确全局后处理", service: "/path/test124", wantStatus: 668, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,正确全局后处理", service: "/path/test127", wantStatus: 668, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,正确全局后处理", service: "/path/test130", wantStatus: 668, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledOKFunc, golHandlingFunc: golHandlingErrFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
//有错误全局预处理,有错误全局后处理函数
//该方法的预处理函数都是循环执行  context的返回数据都会被最后以此执行覆盖   所以所有预处理我们都之放一个
func TestHandler8(t *testing.T) {

	type testCase struct {
		name            string
		service         string
		isLimited       bool
		fallback        bool
		handleObj       interface{}
		golHandlingFunc context.Handler
		golHandledFunc  context.Handler
		wantStatus      int
		wantSpecial     string
		wantContent     string
	}
	tests := []*testCase{
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无后处理,错误全局后处理", service: "/path/test35", wantStatus: 668, wantSpecial: "", handleObj: &testObj9{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,错误全局后处理", service: "/path/test38", wantStatus: 668, wantSpecial: "", handleObj: &testObj10{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,错误全局后处理", service: "/path/test41", wantStatus: 668, wantSpecial: "", handleObj: &testObj11{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无后处理,错误全局后处理", service: "/path/test62", wantStatus: 668, wantSpecial: "", handleObj: &testObj12{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,错误全局后处理", service: "/path/test65", wantStatus: 668, wantSpecial: "", handleObj: &testObj13{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,错误全局后处理", service: "/path/test68", wantStatus: 668, wantSpecial: "", handleObj: &testObj14{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无后处理,错误全局后处理", service: "/path/test98", wantStatus: 668, wantSpecial: "", handleObj: &testObj18{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,错误全局后处理", service: "/path/test101", wantStatus: 668, wantSpecial: "", handleObj: &testObj19{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,错误全局后处理", service: "/path/test104", wantStatus: 668, wantSpecial: "", handleObj: &testObj20{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无后处理,错误全局后处理", service: "/path/test125", wantStatus: 668, wantSpecial: "", handleObj: &testObj21{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,错误全局后处理", service: "/path/test128", wantStatus: 668, wantSpecial: "", handleObj: &testObj22{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
		{name: "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,错误全局后处理", service: "/path/test131", wantStatus: 668, wantSpecial: "", handleObj: &testObj23{}, golHandledFunc: golHandledErrFunc, golHandlingFunc: golHandlingErrFunc},
	}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
		hydra.WithRegistry("lm://."),
	)
	for _, tt := range tests {

		global.Def.ServerTypes = []string{http.API}
		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		mockConf.GetAPI()
		app.API(tt.service, tt.handleObj)
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:   types.XMap{},
			MockTFuncs: map[string]interface{}{},
			MockUser:   &mocks.MockUser{MockRequestID: utility.GetGUID()},
			MockRequest: &mocks.MockRequest{

				MockPath: &mocks.MockPath{
					MockRequestPath:   tt.service,
					MockIsLimit:       tt.isLimited,
					MockAllowFallback: tt.fallback,
				},
			},

			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockAPPConf:  serverConf,
		}
		if tt.golHandlingFunc != nil {
			services.Def.OnHandleExecuting(tt.golHandlingFunc, http.API)
		}

		if tt.golHandledFunc != nil {
			services.Def.OnHandleExecuted(tt.golHandledFunc, http.API)
		}

		//获取中间件
		handler := middleware.ExecuteHandler(tt.service)

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantSpecial, gotSpecial)
	}
}

//author:taoshouyin
//time:2020-11-13
//desc:测试handle执行函数
func TestRPCHandler(t *testing.T) {
	// reqUrl := "/single/hydra/newversion/md5/auth@authserver.sas_debug"
	global.Def.RegistryAddr = "zk://192.168.0.101"
}

// "Handler-限流,不降级 testObj1
type testObj1 struct{}

func (n *testObj1) Handle(ctx context.IContext) interface{} { return nil }
func (n *testObj1) Close() error                            { return nil }

// "Handler-限流,降级,无函数 testObj2
type testObj2 struct{}

func (n *testObj2) Handle(ctx context.IContext) interface{} { return nil }
func (n *testObj2) Close() error                            { return nil }

// "Handler-限流,降级,有函数,异常输出 testObj3
type testObj3 struct{}

func (n *testObj3) Handle(ctx context.IContext) interface{} { return nil }
func (n *testObj3) Fallback(ctx context.IContext) interface{} {
	return errs.NewError(611, "限流,降级,有函数,异常输出")
}
func (n *testObj3) Close() error { return nil }

// "Handler-限流,降级,有函数,异常输出1 testObj4
type testObj4 struct{}

func (n *testObj4) Handle(ctx context.IContext) interface{} { return nil }
func (n *testObj4) Fallback(ctx context.IContext) interface{} {
	return fmt.Errorf("限流,降级,有函数,异常输出")
}
func (n *testObj4) Close() error { return nil }

// "Handler-限流,降级,有函数,正常输出 testObj5
type testObj5 struct{}

func (n *testObj5) Handle(ctx context.IContext) interface{}   { return nil }
func (n *testObj5) Fallback(ctx context.IContext) interface{} { return "fallsuccess" }
func (n *testObj5) Close() error                              { return nil }

// "Handler-不限流,有错误主服务,无ALL预处理,无ALL后处理" testObj6
// "Handler-不限流,有错误主服务,无ALL预处理,无后处理,正确全局后处理" testObj6
// "Handler-不限流,有错误主服务,无ALL预处理,无后处理,错误全局后处理" testObj6
type testObj6 struct{}

func (n *testObj6) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,无ALL预处理,无ALL后处理")
}
func (n *testObj6) Close() error { return nil }

// "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,无全局后处理" testObj7
// "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,正确全局后处理" testObj7
// "Handler-不限流,有错误主服务,无ALL预处理,正确后处理,错误全局后处理" testObj7
type testObj7 struct{}

func (n *testObj7) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,无ALL预处理,正确后处理")
}
func (n *testObj7) Handled(ctx context.IContext) interface{} { return "handledsuccess" }
func (n *testObj7) Close() error                             { return nil }

// "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,无全局后处理" testObj8
// "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,正确全局后处理" testObj8
// "Handler-不限流,有错误主服务,无ALL预处理,错误后处理,错误全局后处理" testObj8
type testObj8 struct{}

func (n *testObj8) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,无ALL预处理,错误后处理")
}
func (n *testObj8) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有错误主服务,无ALL预处理,错误后处理")
}
func (n *testObj8) Close() error { return nil }

// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无ALL后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无后处理,正确全局后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,无后处理,错误全局后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无ALL后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无后处理,正确全局后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,无后处理,错误全局后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无ALL后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无后处理,正确全局后处理" testObj9
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,无后处理,错误全局后处理" testObj9
type testObj9 struct{}

func (n *testObj9) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,正确预处理,无全局预处理,无ALL后处理")
}
func (n *testObj9) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj9) Close() error                              { return nil }

// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,无全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,正确全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,错误全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,无全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,正确全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,正确后处理,错误全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,无全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,正确全局后处理" testObj10
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,正确后处理,错误全局后处理"  testObj10
type testObj10 struct{}

func (n *testObj10) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,正确预处理,无全局预处理,正确后处理,无全局后处理")
}
func (n *testObj10) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj10) Handled(ctx context.IContext) interface{}  { return "handledsuccess" }
func (n *testObj10) Close() error                              { return nil }

// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,无全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,正确全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,无全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,正确全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,正确全局预处理,错误后处理,错误全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,无全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,正确全局后处理" testObj11
// "Handler-不限流,有错误主服务,正确预处理,错误全局预处理,错误后处理,错误全局后处理" testObj11
type testObj11 struct{}

func (n *testObj11) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理")
}
func (n *testObj11) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj11) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有错误主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理")
}
func (n *testObj11) Close() error { return nil }

// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无ALL后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无后处理,正确全局后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,无后处理,错误全局后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无ALL后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无后处理,正确全局后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,无后处理,错误全局后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无ALL后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无后处理,正确全局后处理" testObj12
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,无后处理,错误全局后处理"testObj12
type testObj12 struct{}

func (n *testObj12) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,错误预处理,无全局预处理,无ALL后处理")
}
func (n *testObj12) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有错误主服务,错误预处理,无全局预处理,无ALL后处理")
}
func (n *testObj12) Close() error { return nil }

// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,无全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,正确全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,错误全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,无全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,正确全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,正确后处理,错误全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,无全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,正确全局后处理" testObj13
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,正确后处理,错误全局后处理" testObj13
type testObj13 struct{}

func (n *testObj13) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,无全局后处理")
}
func (n *testObj13) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有错误主服务,错误预处理,无全局预处理,正确后处理,无全局后处理")
}
func (n *testObj13) Handled(ctx context.IContext) interface{} { return "handledsuccess" }
func (n *testObj13) Close() error                             { return nil }

// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,无全局后处理"testObj14
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,正确全局后处理" testObj14
// "Handler-不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,错误全局后处理" testObj14
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,无全局后处理" testObj14
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,正确全局后处理" testObj14
// "Handler-不限流,有错误主服务,错误预处理,正确全局预处理,错误后处理,错误全局后处理"  testObj14
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,无全局后处理"testObj14
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,正确全局后处理" testObj14
// "Handler-不限流,有错误主服务,错误预处理,错误全局预处理,错误后处理,错误全局后处理"  testObj14
type testObj14 struct{}

func (n *testObj14) Handle(ctx context.IContext) interface{} {
	return errs.NewError(666, "不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj14) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj14) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有错误主服务,错误预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj14) Close() error { return nil }

// "Handler-不限流,有正确主服务,无ALL预处理,无ALL后处理", testObj15
// "Handler-不限流,有正确主服务,无ALL预处理,无后处理,正确全局后处理" testObj15
// "Handler-不限流,有正确主服务,无ALL预处理,无后处理,错误全局后处理" testObj15
type testObj15 struct{}

func (n *testObj15) Handle(ctx context.IContext) interface{} { return "mainsuccess" }
func (n *testObj15) Close() error                            { return nil }

// "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,无全局后处理" testObj16
// "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,正确全局后处理" testObj16
// "Handler-不限流,有正确主服务,无ALL预处理,正确后处理,错误全局后处理" testObj16
type testObj16 struct{}

func (n *testObj16) Handle(ctx context.IContext) interface{}  { return "mainsuccess" }
func (n *testObj16) Handled(ctx context.IContext) interface{} { return "handledsuccess" }
func (n *testObj16) Close() error                             { return nil }

// "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,无全局后处理" testObj17
// "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,正确全局后处理" testObj17
// "Handler-不限流,有正确主服务,无ALL预处理,错误后处理,错误全局后处理" testObj17
type testObj17 struct{}

func (n *testObj17) Handle(ctx context.IContext) interface{} { return "mainsuccess" }
func (n *testObj17) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有正确主服务,无ALL预处理,错误后处理")
}
func (n *testObj17) Close() error { return nil }

// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无ALL后处理" testObj18
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无后处理,正确全局后处理" testObj18
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,无后处理,错误全局后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无ALL后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无后处理,正确全局后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,无后处理,错误全局后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无ALL后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无后处理,正确全局后处理"testObj18
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,无后处理,错误全局后处理"testObj18
type testObj18 struct{}

func (n *testObj18) Handle(ctx context.IContext) interface{}   { return "mainsuccess" }
func (n *testObj18) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj18) Close() error                              { return nil }

// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,无全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,正确全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,正确后处理,错误全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,无全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,正确全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,正确后处理,错误全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,无全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,正确全局后处理" testObj19
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,正确后处理,错误全局后处理" testObj19
type testObj19 struct{}

func (n *testObj19) Handle(ctx context.IContext) interface{}   { return "mainsuccess" }
func (n *testObj19) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj19) Handled(ctx context.IContext) interface{}  { return "handledsuccess" }
func (n *testObj19) Close() error                              { return nil }

// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,无全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,正确全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,错误全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,无全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,正确全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,正确全局预处理,错误后处理,错误全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,无全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,正确全局后处理" testObj20
// "Handler-不限流,有正确主服务,正确预处理,错误全局预处理,错误后处理,错误全局后处理" testObj20
type testObj20 struct{}

func (n *testObj20) Handle(ctx context.IContext) interface{}   { return "mainsuccess" }
func (n *testObj20) Handling(ctx context.IContext) interface{} { return "handlingsuccess" }
func (n *testObj20) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有正确主服务,正确预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj20) Close() error { return nil }

// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无ALL后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无后处理,正确全局后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,无后处理,错误全局后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无ALL后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无后处理,正确全局后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,无后处理,错误全局后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无ALL后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无后处理,正确全局后处理" testObj21
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,无后处理,错误全局后处理" testObj21
type testObj21 struct{}

func (n *testObj21) Handle(ctx context.IContext) interface{} { return "mainsuccess" }
func (n *testObj21) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有正确主服务,错误预处理,无全局预处理,无ALL后处理")
}
func (n *testObj21) Close() error { return nil }

// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,无全局后处理"  testObj22
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,正确全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,错误全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,无全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,正确全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,正确后处理,错误全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,无全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,正确全局后处理" testObj22
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,正确后处理,错误全局后处理" testObj22
type testObj22 struct{}

func (n *testObj22) Handle(ctx context.IContext) interface{} { return "mainsuccess" }
func (n *testObj22) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有正确主服务,错误预处理,无全局预处理,正确后处理,无全局后处理")
}
func (n *testObj22) Handled(ctx context.IContext) interface{} { return "handledsuccess" }
func (n *testObj22) Close() error                             { return nil }

// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,无全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,正确全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,错误全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,无全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,正确全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,正确全局预处理,错误后处理,错误全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,无全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,正确全局后处理" testObj23
// "Handler-不限流,有正确主服务,错误预处理,错误全局预处理,错误后处理,错误全局后 testObj23
type testObj23 struct{}

func (n *testObj23) Handle(ctx context.IContext) interface{} { return "mainsuccess" }
func (n *testObj23) Handling(ctx context.IContext) interface{} {
	return errs.NewError(667, "不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj23) Handled(ctx context.IContext) interface{} {
	return errs.NewError(670, "不限流,有正确主服务,错误预处理,无全局预处理,错误后处理,无全局后处理")
}
func (n *testObj23) Close() error { return nil }
