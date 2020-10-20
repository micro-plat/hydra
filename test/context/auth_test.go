package context

import (
	"testing"

	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/errs"
)

func Test_Auth_Response(t *testing.T) {

	c := &ctx.Auth{} //构建对象

	//测试参数为空
	gotNil := c.Response()
	assert.Equal(t, nil, gotNil, "参数为空")

	//测试参数不为空
	response := c.Response(1)
	assert.Equal(t, 1, response, "参数不为空")

	//测试获取response
	gotResponse := c.Response()
	assert.Equal(t, 1, gotResponse, "获取response")
}

func Test_Auth_Request(t *testing.T) {

	c := &ctx.Auth{} //构建对象

	//测试参数为空
	gotNil := c.Request()
	assert.Equal(t, nil, gotNil, "参数为空")

	//测试参数不为空
	request := c.Request(1)
	assert.Equal(t, 1, request, "参数不为空")

	//获取request
	gotRequest := c.Request()
	assert.Equal(t, 1, gotRequest, "获取request")

}

func Test_Auth_Bind(t *testing.T) {
	type result struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	tests := []struct {
		name           string
		request        interface{}
		want           interface{}
		wantPanicError bool
		def            bool
	}{
		{name: "request为空", wantPanicError: true},
		{name: "request为func返回空值", request: func() interface{} {
			return nil
		}, wantPanicError: true},
		{name: "request为func返回非空值", request: func() interface{} {
			return result{Key: "1", Value: "1"}
		}, want: result{Key: "1", Value: "1"}, wantPanicError: false},
		{name: "request为错误的json字符串", request: `{"key":"1",v}`, wantPanicError: true},
		{name: "request为json字符串", request: `{"key":"1","value":"1"}`, want: result{Key: "1", Value: "1"}, wantPanicError: false},
		{name: "默认情况", request: map[string]string{"key": "value"}, def: true, want: map[string]string{"key": "value"}, wantPanicError: false},
	}
	for _, tt := range tests {
		c := &ctx.Auth{}
		c.Request(tt.request)
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(*errs.Error); ok {
					assert.Equal(t, 401, e.GetCode(), tt.name)
					assert.Equal(t, "请求中未包含用户信息,用户未登录", e.Error(), tt.name)
				}
				assert.Equal(t, tt.wantPanicError, r != nil, tt.name)
			}
		}()

		if !tt.def {
			out := result{}
			c.Bind(&out)
			assert.Equal(t, tt.want, out, tt.name)
			return
		}

		out := map[string]string{}
		c.Bind(&out)
		assert.Equal(t, tt.want, out, tt.name)
	}
}
