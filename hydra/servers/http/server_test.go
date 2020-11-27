package http

import (
	xnet "net"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/net"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
		addr       string
		routers    []*router.Router
		opts       []Option
		wantT      *Server
		wantErr    string
	}{
		{name: "1. NewServer-地址不合法", serverName: "api", addr: "xxxllsd", wantErr: "端口不合法 xxxllsd"},
		{name: "2. NewServer-构建http服务", serverName: "api", addr: "127.0.0.1:8080", opts: []Option{WithServerType("api"), WithTimeout(30, 30, 30), WithTLS([]string{"pem", "key"})}},
	}
	for _, tt := range tests {
		gotT, err := NewServer(tt.serverName, tt.addr, tt.routers, tt.opts...)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, "http", gotT.proto, tt.name)
		assert.Equal(t, net.GetLocalIPAddress(), gotT.ip, tt.name)
		o := &option{}
		for _, v := range tt.opts {
			v(o)
		}
		assert.Equal(t, o.readTimeout, gotT.option.readTimeout, tt.name)
		assert.Equal(t, o.readHeaderTimeout, gotT.option.readHeaderTimeout, tt.name)
		assert.Equal(t, o.writeTimeout, gotT.option.writeTimeout, tt.name)
		assert.Equal(t, o.tls, gotT.option.tls, tt.name)

		host, port, _ := global.GetHostPort(tt.addr)
		assert.Equal(t, xnet.JoinHostPort(host, port), gotT.server.Addr, tt.name)
		assert.Equal(t, time.Duration(o.readHeaderTimeout)*time.Second, gotT.server.ReadHeaderTimeout, tt.name)
		assert.Equal(t, time.Duration(o.readTimeout)*time.Second, gotT.server.ReadTimeout, tt.name)
		assert.Equal(t, time.Duration(o.writeTimeout)*time.Second, gotT.server.WriteTimeout, tt.name)
	}
}

func TestNewWSServer(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
		addr       string
		routers    []*router.Router
		opts       []Option
		wantT      *Server
		wantErr    string
	}{
		{name: "1. NewWSServer-地址不合法", serverName: "ws", addr: "xxxllsd", wantErr: "端口不合法 xxxllsd"},
		{name: "2. NewWSServer-构建ws服务", serverName: "ws", addr: "127.0.0.1:8080", opts: []Option{WithServerType("api"), WithTimeout(30, 30, 30), WithTLS([]string{"pem", "key"})}},
	}
	for _, tt := range tests {
		gotT, err := NewWSServer(tt.serverName, tt.addr, tt.routers, tt.opts...)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, "ws", gotT.proto, tt.name)
		assert.Equal(t, net.GetLocalIPAddress(), gotT.ip, tt.name)
		o := &option{}
		for _, v := range tt.opts {
			v(o)
		}
		assert.Equal(t, o.readTimeout, gotT.option.readTimeout, tt.name)
		assert.Equal(t, o.readHeaderTimeout, gotT.option.readHeaderTimeout, tt.name)
		assert.Equal(t, o.writeTimeout, gotT.option.writeTimeout, tt.name)
		assert.Equal(t, o.tls, gotT.option.tls, tt.name)

		host, port, _ := global.GetHostPort(tt.addr)
		assert.Equal(t, host, gotT.host, tt.name)
		assert.Equal(t, port, gotT.port, tt.name)
		assert.Equal(t, xnet.JoinHostPort(host, port), gotT.server.Addr, tt.name)
		assert.Equal(t, time.Duration(o.readHeaderTimeout)*time.Second, gotT.server.ReadHeaderTimeout, tt.name)
		assert.Equal(t, time.Duration(o.readTimeout)*time.Second, gotT.server.ReadTimeout, tt.name)
		assert.Equal(t, time.Duration(o.writeTimeout)*time.Second, gotT.server.WriteTimeout, tt.name)
	}
}
