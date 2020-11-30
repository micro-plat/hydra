package rpc

import (
	orcontext "context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/test/mocks"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/errs"
)

type rpcServerObj struct{}

func (n *rpcServerObj) Handle(ctx context.IContext) interface{} { return "success" }

type rpcServerObj1 struct{}

func (n *rpcServerObj1) Handle(ctx context.IContext) interface{} {
	return errs.NewError(668, fmt.Errorf("全局预处理异常"))
}

func TestNewProcessor(t *testing.T) {
	services.Def.RPC("/rpcserver/taosy/test", &rpcServerObj{})
	tests := []struct {
		name    string
		routers []*router.Router
	}{
		{name: "1. 添加空路由", routers: []*router.Router{}},
		{name: "2. 添加单条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"Get"})}},
		{name: "3. 添加多条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"Get"}), router.NewRouter("/rpcserver/taosy/test1", "/rpcserver/taosy/test1", []string{"Post"})}},
	}
	for _, tt := range tests {
		gotP := rpc.NewProcessor(tt.routers...)
		assert.Equalf(t, 15, len(gotP.Engine.RouterGroup.Handlers), tt.name+",中间件数量")
		assert.Equalf(t, len(tt.routers), len(gotP.Engine.Routes()), tt.name+",路由数量")
	}
}

func TestProcessor_Request(t *testing.T) {
	services.Def.RPC("/rpcserver/taosy/test1", &rpcServerObj1{})
	services.Def.RPC("/rpcserver/taosy/test2", &rpcServerObj{})
	conf := mocks.NewConfBy("server_rpc_pross_test", "tessRpctes")
	conf.RPC(":41501")
	app.Cache.Save(conf.GetRPCConf())

	tests := []struct {
		name    string
		fields  *rpc.Processor
		context orcontext.Context
		request *pb.RequestContext
		wantP   *pb.ResponseContext
		wantErr string
	}{
		{name: "1. Processor_Request-设置错误的request对象", fields: rpc.NewProcessor([]*router.Router{}...), context: nil, request: &pb.RequestContext{Header: "错误数据"}, wantP: &pb.ResponseContext{Status: http.StatusNotAcceptable, Header: "", Result: "输入参数有误"}, wantErr: ""},
		{name: "2. Processor_Request-设置请求路径不存在", fields: rpc.NewProcessor([]*router.Router{router.NewRouter("/rpcserver/taosy/test2", "/rpcserver/taosy/test2", []string{"GET"})}...), context: nil, request: &pb.RequestContext{Service: "/rpcserver/taosy/testx", Method: "GET", Header: `{"Host":"baidu.com"}`, Input: "{}"}, wantP: &pb.ResponseContext{Status: 404, Header: "", Result: "404 service not found"}, wantErr: ""},
		{name: "3. Processor_Request-设置错误的请求路径", fields: rpc.NewProcessor([]*router.Router{router.NewRouter("/rpcserver/taosy/test1", "/rpcserver/taosy/test1", []string{"GET"})}...), context: nil, request: &pb.RequestContext{Service: "/rpcserver/taosy/test1", Method: "GET", Header: `{"Host":"baidu.com"}`, Input: "{}"}, wantP: &pb.ResponseContext{Status: 668, Header: "", Result: "Internal Server Error"}, wantErr: ""},
		{name: "4. Processor_Request-设置正确请求路径", fields: rpc.NewProcessor([]*router.Router{router.NewRouter("/rpcserver/taosy/test2", "/rpcserver/taosy/test2", []string{"GET"})}...), context: nil, request: &pb.RequestContext{Service: "/rpcserver/taosy/test2", Method: "GET", Header: `{"Host":"baidu.com"}`, Input: "{}"}, wantP: &pb.ResponseContext{Status: 200, Header: "", Result: "success"}, wantErr: ""},
	}
	for _, tt := range tests {
		gotP, err := tt.fields.Request(tt.context, tt.request)
		if tt.wantErr != "" {
			assert.Equalf(t, true, strings.Contains(err.Error(), tt.wantErr), tt.name, err, tt.wantErr)
		}

		assert.Equalf(t, tt.wantP.Status, gotP.Status, tt.name, tt.wantP.Status, gotP.Status)
		assert.Equalf(t, true, strings.Contains(gotP.Result, tt.wantP.Result), tt.name, gotP.Result, tt.wantP.Result)
	}
}
