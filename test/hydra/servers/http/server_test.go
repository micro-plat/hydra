package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	xhttp "net/http"
	"os"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/types"
)

func TestServer_GetAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
		h    []string
		want string
	}{
		{name: "参数为空", addr: "127.0.0.1:58080", want: "ws://127.0.0.1:58080"},
		{name: "参数为ip", addr: "127.0.0.1:58080", h: []string{"192.168.0.1"}, want: "ws://192.168.0.1:58080"},
		{name: "参数为ip", addr: "0.0.0.0:58080", h: []string{}, want: fmt.Sprintf("ws://%s:58080", global.LocalIP())},
	}
	for _, tt := range tests {
		s, _ := http.NewWSServer("ws", tt.addr, []*router.Router{})
		got := s.GetAddress(tt.h...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func setAPICacheConf() {
	confObj := mocks.NewConfBy("hydra_server_test", "servertest") //构建对象
	confObj.API(":58080")                                         //初始化参数
	serverConf := confObj.GetAPIConf()                            //获取配置
	_, _ = http.NewResponsive(serverConf)
}

func testRequest(addr string) {
	xhttp.NewRequest("GET", addr, nil)
}

func TestServer_Start_WithErr(t *testing.T) {
	tests := []struct {
		name             string
		serverName       string
		addr             string
		opts             []http.Option
		wantErr          string
		wantRequestPanic string
	}{
		{name: "错误的ssl", serverName: "", addr: "127.0.0.1:58081", opts: []http.Option{http.WithTLS([]string{"pem", "key"})}, wantErr: "open pem: no such file or directory"},
		{name: "未设置serverType", serverName: "", addr: "127.0.0.1:58081", opts: []http.Option{}, wantRequestPanic: "未找到的缓存配置信息"},
		{name: "没有保存api对应的缓存配置", serverName: "", addr: "127.0.0.1:58082", opts: []http.Option{http.WithServerType(global.API)}, wantRequestPanic: "未找到api的缓存配置信息"},
	}

	for _, tt := range tests {
		s, err := http.NewServer("api", tt.addr, []*router.Router{}, tt.opts...)
		assert.Equal(t, nil, err, tt.name)
		err = s.Start()
		defer s.Shutdown()
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		if tt.wantRequestPanic != "" {
			//构建的新的os.Stderr
			rescueStderr := os.Stderr
			r, w, _ := os.Pipe()
			*os.Stderr = *w

			req, _ := xhttp.NewRequest("GET", fmt.Sprintf("http://%s", tt.addr), nil)
			client := &xhttp.Client{}
			client.Do(req)

			//获取输出
			w.Close()
			out, err := ioutil.ReadAll(r)
			assert.Equalf(t, false, err != nil, tt.name)
			fmt.Println(string(out))
			//	assert.Equalf(t, true, strings.Contains(string(out), tt.wantRequestPanic), tt.name)
			//还原os.Stderr
			*os.Stderr = *rescueStderr
		}
	}
}

func doTestRequest(addr string, isSSL bool) (*xhttp.Response, error) {
	req, _ := xhttp.NewRequest("GET", addr, nil)

	client := &xhttp.Client{}

	// 设置跳过不安全的 HTTPS
	if isSSL {
		tls11Transport := &xhttp.Transport{
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		}
		client.Transport = tls11Transport
	}

	return client.Do(req)
}

func TestServer_Start_WithSSL(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
		addr       string
		routers    []*router.Router
		opts       []http.Option
		isSSL      bool
	}{
		{name: "启动不带有ssl证书的服务", serverName: "", addr: "127.0.0.1:58083", opts: []http.Option{http.WithServerType(global.API)}},
		{name: "启动带有ssl证书的服务", serverName: "", addr: "127.0.0.1:58084", isSSL: true, opts: []http.Option{http.WithServerType(global.API), http.WithTLS([]string{"server_test_crt.txt", "server_test_key.txt"})}},
	}

	confObj := mocks.NewConfBy("hydra_server_test1", "servertest1") //构建对象
	confObj.API(":58081")                                           //初始化参数
	serverConf := confObj.GetAPIConf()                              //获取配置
	_, _ = http.NewResponsive(serverConf)

	for _, tt := range tests {
		s, err := http.NewServer("api", tt.addr, []*router.Router{}, tt.opts...)
		assert.Equal(t, nil, err, tt.name)
		err = s.Start()
		assert.Equal(t, nil, err, tt.name)

		//对启动服务进行访问
		proto := types.DecodeString(tt.isSSL, true, "https://", "http://")
		resp, err := doTestRequest(fmt.Sprintf("%s%s", proto, tt.addr), tt.isSSL)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, "404 Not Found", resp.Status, tt.name)
		assert.Equal(t, 404, resp.StatusCode, tt.name)
		if tt.isSSL {
			assert.Equal(t, true, resp.TLS != nil, tt.name)
		}

		err = s.Shutdown()
		assert.Equal(t, false, err != nil, tt.name)
	}
}
