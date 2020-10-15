package ctx

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/net"
)

func Test_user_GetRequestID(t *testing.T) {
	type fields struct {
		ctx       context.IInnerContext
		requestID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "header无X-Request-Id,user存在requestID", fields: fields{
			ctx:       &TestContxt{},
			requestID: "456789"}, want: "456789"},
		{name: "header中有X-Request-Id", fields: fields{
			ctx: &TestContxt{
				header: http.Header{"X-Request-Id": []string{"123456"}}},
			requestID: "456789"}, want: "123456"},
		{name: "header无X-Request-Id,user不存在requestID", fields: fields{
			ctx: &TestContxt{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newUser(tt.fields.ctx, conf.NewMeta())
			c.requestID = tt.fields.requestID
			if got := c.GetRequestID(); got != tt.want {
				if tt.want == "" { //返回随机结果
					if len(got) == 9 {
						return
					}
				}
				t.Errorf("user.GetRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_GetClientIP(t *testing.T) {
	type fields struct {
		ctx context.IInnerContext
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "ctx中ip为127.0.0.1", fields: fields{
			ctx: &TestContxt{
				clientIP: "127.0.0.1",
			}}, want: net.GetLocalIPAddress()},
		{name: "ctx中ip为::1", fields: fields{
			ctx: &TestContxt{
				clientIP: "127.0.0.1",
			}}, want: net.GetLocalIPAddress()},
		{name: "ctx中ip非本地ip", fields: fields{
			ctx: &TestContxt{
				clientIP: "192.168.9.99",
			}}, want: "192.168.9.99"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newUser(tt.fields.ctx, conf.NewMeta())
			if got := c.GetClientIP(); got != tt.want {
				t.Errorf("user.GetClientIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
