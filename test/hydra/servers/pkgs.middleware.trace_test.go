package servers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestTrace(t *testing.T) {

	tests := []struct {
		name           string
		option         []api.Option
		requestMap     map[string]interface{}
		responseStatus int
		wantSpecial    string
		wantDebug1     string
		wantDebug2     string
	}{
		{name: "api未配置trace", responseStatus: 200, wantSpecial: ""},
		{name: "api配置trace", option: []api.Option{api.WithTrace()}, responseStatus: 200, wantSpecial: "trace", wantDebug1: "> trace.request: map[]", wantDebug2: "> trace.response: 200"},
	}

	for _, tt := range tests {
		//初始化测试用例参数
		c := mocks.NewConfBy("middleware_trace_test", "trace")
		c.API(":9090", tt.option...)
		ctx := &mocks.MiddleContext{
			MockNext:     func() { fmt.Println("output") },
			MockUser:     &mocks.MockUser{MockClientIP: "127.0.0.1", MockRequestID: "06c6fb24c"},
			MockRequest:  &mocks.MockRequest{MockQueryMap: tt.requestMap},
			MockResponse: &mocks.MockResponse{MockStatus: tt.responseStatus},
			MockAPPConf:  c.GetAPIConf(),
		}

		//构建的新的os.Stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//调用中间件
		handler := middleware.Trace()
		handler(ctx)
		time.Sleep(time.Second * 3)

		//获取输出
		w.Close()
		out, err := ioutil.ReadAll(r)
		assert.Equalf(t, false, err != nil, tt.name)

		//还原os.Stdout
		os.Stdout = rescueStdout
		assert.Equalf(t, true, strings.Contains(string(out), "output"), tt.name)

		if tt.wantSpecial != "" {
			gotSpecial := ctx.Response().GetSpecials()
			assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name)
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantDebug1), tt.name)
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantDebug2), tt.name)
		}

	}

}
