package rpc

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	xrpc "github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	crpc "github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

type errorObj struct{}

func (n *errorObj) Handle(ctx context.IContext) interface{} {
	return errs.NewError(670, "返回错误的数据")
}

type okObj struct{}

func (n *okObj) Handle(ctx context.IContext) interface{} { return "success" }

type okObj1 struct{}

func (n *okObj1) Handle(ctx context.IContext) interface{} {
	m, _ := ctx.Request().GetMap()
	b, _ := json.Marshal(m)
	return string(b)
}

//author:taoshouyin
//time:2020-11-18
//desc:测试rpc服务无路由初始化
func TestServer(t *testing.T) {
	localIP := global.LocalIP()
	tests := []struct {
		name    string
		address string
		routers []*router.Router
		handle  interface{}
		reqUrl  string
		path    string
		params  map[string]interface{}

		wantStatus  int
		wantContent string
	}{
		{name: "rpc-server 注册返回错误的服务", address: ":8091", reqUrl: fmt.Sprintf("tcp://%s:8091", localIP), path: "/rpc/server/test", handle: &errorObj{}, params: map[string]interface{}{},
			wantStatus: 670, wantContent: "Internal Server Error", routers: []*router.Router{router.NewRouter("/rpc/server/test", "/rpc/server/test", []string{"Get"})}},
		{name: "rpc-server 注册返回正确的服务", address: "", reqUrl: fmt.Sprintf("tcp://%s:8090", localIP), path: "/rpc/server/test1", handle: &okObj{}, params: map[string]interface{}{},
			wantStatus: 200, wantContent: "success", routers: []*router.Router{router.NewRouter("/rpc/server/test1", "/rpc/server/test1", []string{"Get"})}},
		{name: "rpc-server 注册返回正确的服务,有参数", address: "", reqUrl: fmt.Sprintf("tcp://%s:8090", localIP), path: "/rpc/server/test2", handle: &okObj1{}, params: map[string]interface{}{"taosy": "testrpcserver"},
			wantStatus: 200, wantContent: `{"taosy":"testrpcserver"}`, routers: []*router.Router{router.NewRouter("/rpc/server/test2", "/rpc/server/test2", []string{"Get"})}},
	}
	for _, tt := range tests {
		mockConf := mocks.NewConfBy("rpacserve_test", "testrpcserver")
		mockConf.RPC(":51001")
		serverConf := mockConf.GetRPCConf()
		app.Cache.Save(serverConf)
		services.Def.RPC(tt.path, tt.handle)

		server, err := rpc.NewServer(tt.address, tt.routers, crpc.DefaultMaxRecvMsgSize, crpc.DefaultMaxSendMsgSize)
		assert.Equalf(t, true, err == nil, tt.name+"server error")

		err = server.Start()
		assert.Equalf(t, true, err == nil, tt.name+"start error")
		time.Sleep(1 * time.Second)

		rclient, err := xrpc.NewClient(tt.reqUrl, "", "")
		assert.Equalf(t, true, err == nil, tt.name+"rpc cilent error")

		ctx := &mocks.MiddleContext{
			MockMeta: types.XMap{},
			MockRequest: &mocks.MockRequest{
				MockPath:     &mocks.MockPath{MockRequestPath: tt.path},
				MockQueryMap: tt.params,
			},
			MockAPPConf: serverConf,
		}

		resp, err := rclient.Request(ctx.Context(), tt.path, tt.params)
		assert.Equalf(t, true, err == nil, tt.name+"rpc request error", err)
		assert.Equalf(t, tt.wantStatus, resp.Status, tt.name+"rpc request Status")
		assert.Equalf(t, tt.wantContent, resp.Result, tt.name+"rpc request Result")

		server.Shutdown()
	}
}

var oncelock sync.Once

var rclient *xrpc.Client

var serverConf app.IAPPConf

//并发测试rpc服务器调用性能
func BenchmarkRPCServer(b *testing.B) {
	oncelock.Do(func() {
		mockConf := mocks.NewConfBy("rpacserve", "Benchmarktestserver")
		mockConf.RPC(":51001")
		serverConf = mockConf.GetRPCConf()
		app.Cache.Save(serverConf)
		services.Def.RPC("/rpc/server/test1", &okObj{})

		routers := []*router.Router{router.NewRouter("/rpc/server/test1", "/rpc/server/test1", []string{"Get"})}
		server, err := rpc.NewServer(":8092", routers, crpc.DefaultMaxRecvMsgSize, crpc.DefaultMaxSendMsgSize)
		assert.Equalf(b, true, err == nil, "server 初始化 error")

		err = server.Start()
		assert.Equalf(b, true, err == nil, "server 启动 error")
		time.Sleep(1 * time.Second)

		rclient, err = xrpc.NewClient("http://"+global.LocalIP()+":8092", "", "")
		assert.Equalf(b, true, err == nil, "rpc cilent 初始化 error")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &mocks.MiddleContext{
			MockMeta: types.XMap{},
			MockRequest: &mocks.MockRequest{
				MockPath:     &mocks.MockPath{MockRequestPath: "/rpc/server/test1"},
				MockQueryMap: map[string]interface{}{},
			},
			MockAPPConf: serverConf,
		}

		resp, err := rclient.Request(ctx.Context(), "/rpc/server/test1", map[string]interface{}{})
		assert.Equalf(b, true, err == nil, "rpc request error")
		assert.Equalf(b, 200, resp.Status, "rpc request Status")
		assert.Equalf(b, "success", resp.Result, "rpc request Result")
	}
}

func TestGetAddress(t *testing.T) {

	tests := []struct {
		name    string
		addr    string
		want    string
		wantErr string
	}{
		{name: "1.传入空字符串", addr: "", want: global.LocalIP() + ":8090", wantErr: ""},
		{name: "2.三段地址不合法", addr: "localhost:tttt:ddd", want: "", wantErr: "地址不合法"},
		{name: "3.端口错误", addr: "localhost:54545454", want: "", wantErr: "端口不合法"},
		{name: "4.一段字母", addr: "tsosy", want: "", wantErr: "地址不合法"},
		{name: "5.一段数字,是端口号", addr: "12346", want: global.LocalIP() + ":12346", wantErr: ""},
		{name: "6.一段数字,不是端口号", addr: "55554444", want: "", wantErr: "地址不合法"},
		{name: "7.冒号+一段字母", addr: ":tsosy", want: "", wantErr: ""},
		{name: "8.冒号+一段数字,是端口号", addr: ":12346", want: global.LocalIP() + ":12346", wantErr: ""},
		{name: "9.冒号+一段数字,不是端口号", addr: ":55554444", want: "", wantErr: "端口不合法"},
		{name: "10.两段字母", addr: "taosy:test", want: "", wantErr: "地址不合法"},
		{name: "11.两段数字+字母", addr: "192.168.0.111:tast", want: "", wantErr: "地址不合法"},
		{name: "12.两段数字+字母,本地ip", addr: global.LocalIP() + ":tast", want: "", wantErr: "端口不合法"},
		{name: "13.两段字母+数字,是端口号", addr: "ksdfdksj:5454", want: "", wantErr: "地址不合法"},
		{name: "14.两段字母+数字,不是是端口号", addr: "ksdfdksj:5454333", want: "", wantErr: "地址不合法"},

		{name: "15.localhost单段", addr: "localhost", want: "localhost:8090", wantErr: ""},
		{name: "16.localhost单段+冒号", addr: "localhost:", want: "", wantErr: "端口不合法"},
		{name: "17.localhost单段+前冒号", addr: ":localhost", want: "", wantErr: "端口不合法"},
		{name: "18.localhost双段+错误端口", addr: "localhost:rrrrr", want: "", wantErr: "端口不合法"},
		{name: "19.localhost双段+正确端口", addr: "localhost:5454", want: "localhost:5454", wantErr: ""},

		{name: "20.0.0.0.0单段", addr: "0.0.0.0", want: global.LocalIP() + ":8090", wantErr: ""},
		{name: "21.0.0.0.0单段+冒号", addr: "0.0.0.0:", want: "", wantErr: "端口不合法"},
		{name: "22.0.0.0.0单段+前冒号", addr: ":0.0.0.0", want: "", wantErr: "端口不合法"},
		{name: "23.0.0.0.0双段+错误端口", addr: "0.0.0.0:ffff", want: "", wantErr: "端口不合法"},
		{name: "24.0.0.0.0双段+正确端口", addr: "0.0.0.0:5454", want: global.LocalIP() + ":5454", wantErr: ""},

		{name: "25.127.0.0.1单段", addr: "127.0.0.1", want: "127.0.0.1:8090", wantErr: ""},
		{name: "26.127.0.0.1单段+冒号", addr: "127.0.0.1:", want: "", wantErr: "端口不合法"},
		{name: "27.127.0.0.1单段+前冒号", addr: ":127.0.0.1", want: "", wantErr: "端口不合法"},
		{name: "28.127.0.0.1双段+错误端口", addr: "127.0.0.1:rrrr", want: "", wantErr: "端口不合法"},
		{name: "29.127.0.0.1双段+正确端口", addr: "127.0.0.1:5454", want: "127.0.0.1:5454", wantErr: ""},
	}
	for _, tt := range tests {
		got, err := rpc.GetAddress(tt.addr)
		if err != nil {
			assert.Equalf(t, true, strings.Contains(err.Error(), tt.wantErr), tt.name, err)
		}
		assert.Equalf(t, tt.want, got, tt.name, got)
	}
}
