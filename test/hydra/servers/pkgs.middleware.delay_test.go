package servers

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:taoshouyin
//time:2020-11-12
//desc:测试basic验证中间件逻辑
func TestDelay(t *testing.T) {
	type testCase struct {
		name        string
		isSet       bool
		delayTime   string
		wantStatus  int
		header      map[string]interface{}
		wantSpecial string
	}

	tests := []*testCase{
		{name: "1. delay-配置不存在", isSet: false, header: map[string]interface{}{}, delayTime: "0", wantStatus: 200, wantSpecial: ""},
		{name: "2. delay-配置错误的头名称", isSet: true, header: map[string]interface{}{"errorname": []string{"111"}}, delayTime: "0", wantStatus: 200, wantSpecial: ""},
		{name: "3. delay-配置错误的数据", isSet: true, header: map[string]interface{}{"X-Add-Delay": []string{"errdata"}}, delayTime: "0", wantStatus: 200, wantSpecial: "delay"},
		{name: "4. delay-配置延迟0s", isSet: true, header: map[string]interface{}{"X-Add-Delay": []string{"0"}}, delayTime: "0", wantStatus: 200, wantSpecial: "delay"},
		{name: "5. delay-配置延迟1s", isSet: true, header: map[string]interface{}{"X-Add-Delay": []string{"1s"}}, delayTime: "1s", wantStatus: 200, wantSpecial: "delay"},
		{name: "6. delay-配置延迟3秒", isSet: true, header: map[string]interface{}{"X-Add-Delay": []string{"3s"}}, delayTime: "3s", wantStatus: 200, wantSpecial: "delay"},
	}

	for _, tt := range tests {
		mockConf := mocks.NewConfBy("middleware_delay_test", "delay")
		//初始化测试用例参数
		mockConf.GetAPI()
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:     conf.NewMeta(),
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse: &mocks.MockResponse{MockStatus: 200},
			MockRequest: &mocks.MockRequest{
				MockHeader: tt.header,
				MockPath: &mocks.MockPath{
					MockRequestPath: "/delay/test",
				},
			},
			MockAPPConf: serverConf,
		}

		//获取中间件
		handler := middleware.Delay()
		//调用中间件
		start := time.Now()
		handler(ctx)
		restime := time.Now().Sub(start)
		delayDuration, _ := time.ParseDuration(tt.delayTime)
		t1 := decimal.NewFromFloat(restime.Seconds())
		t2 := decimal.NewFromFloat(delayDuration.Seconds())
		assert.Equalf(t, t1.IntPart(), t2.IntPart(), tt.name, t1.IntPart(), t2.IntPart())
		//断言结果
		gotStatus, _, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
	}
}
