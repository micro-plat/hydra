package rpcserver

import (
	"encoding/json"
	"fmt"
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
		mockConf := mocks.NewConf()
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
		mockConf := mocks.NewConf()
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
