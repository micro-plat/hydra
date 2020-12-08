package http

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestWithServerType(t *testing.T) {
	tests := []struct {
		name string
		t    string
	}{
		{name: "1. httpserver-设置api类型", t: "api"},
		{name: "2. httpserver-设置cron类型", t: "cron"},
	}
	for _, tt := range tests {
		f := WithServerType(tt.t)
		o := &option{}
		f(o)
		assert.Equal(t, tt.t, o.serverType, tt.name)
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name              string
		readTimeout       int
		writeTimeout      int
		readHeaderTimeout int
	}{
		{name: "1. httpserver-设置读取超时时间", readTimeout: 30},
		{name: "2. httpserver-设置写入超时时间", writeTimeout: 30},
		{name: "3. httpserver-设置头部读取超时时间", readHeaderTimeout: 30},
	}
	for _, tt := range tests {
		f := WithTimeout(tt.readTimeout, tt.writeTimeout, tt.readHeaderTimeout)
		o := &option{}
		f(o)
		assert.Equal(t, tt.readTimeout, o.readTimeout, tt.name)
		assert.Equal(t, tt.writeTimeout, o.writeTimeout, tt.name)
		assert.Equal(t, tt.readHeaderTimeout, o.readHeaderTimeout, tt.name)
	}
}

func TestWithTLS(t *testing.T) {
	tests := []struct {
		name    string
		tls     []string
		wantTLS []string
	}{
		{name: "1. httpserver-设置空的TLS", tls: []string{}},
		{name: "2. httpserver-设置错误的TLS", tls: []string{"pem"}},
		{name: "3. httpserver-设置争取的TLS", tls: []string{"pem", "key"}, wantTLS: []string{"pem", "key"}},
	}
	for _, tt := range tests {
		f := WithTLS(tt.tls)
		o := &option{}
		f(o)
		assert.Equal(t, tt.wantTLS, o.tls, tt.name)
	}
}
