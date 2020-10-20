package context

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/net"
)

func Test_user_GetRequestID(t *testing.T) {
	//测试requestID通过uuid获取
	c := ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta())
	got1 := c.GetRequestID()
	assert.Equal(t, 9, len(got1), "X-Request-Id不存在,requestID不存在,通过uuid生成requestID")

	//测试X-Request-Id不存在,requestID 存在
	got2 := c.GetRequestID()
	assert.Equal(t, got1, got2, "获取存在的requestID")

	//X-Request-Id存在
	c1 := ctx.NewUser(&mocks.TestContxt{HttpHeader: http.Header{"X-Request-Id": []string{"123456"}}}, conf.NewMeta())
	got := c1.GetRequestID()
	assert.Equal(t, "123456", got, "通过X-Request-Id获取requestID")
}

func Test_user_GetClientIP(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.IInnerContext
		want string
	}{
		{name: "ctx中ip为127.0.0.1", ctx: &mocks.TestContxt{ClientIp: "127.0.0.1"}, want: net.GetLocalIPAddress()},
		{name: "ctx中ip为::1", ctx: &mocks.TestContxt{ClientIp: "127.0.0.1"}, want: net.GetLocalIPAddress()},
		{name: "ctx中ip非本地ip", ctx: &mocks.TestContxt{ClientIp: "192.168.9.99"}, want: "192.168.9.99"},
	}
	for _, tt := range tests {
		c := ctx.NewUser(tt.ctx, conf.NewMeta())
		got := c.GetClientIP()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
