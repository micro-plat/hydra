package context

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/api"
	c "github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/router"
	  "github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_rpath_GetRouter_WithPanic(t *testing.T) {

	confObj := mocks.NewConfBy("context_path_test", "pathctx") //构建对象
	confObj.API(":8080")                                       //初始化参数
	confObj.CRON(c.WithMasterSlave(), c.WithTrace())
	confObj.Service.API.Add("/api", "/api", []string{"GET"})
	httpConf := confObj.GetAPIConf() //获取配置

	tests := []struct {
		name       string
		ctx        context.IInnerContext
		serverConf app.IAPPConf
		meta       conf.IMeta
		want       *router.Router
		wantError  string
	}{
		{name: "1.1 路径正确 请求方法为DELETE", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "DELETE"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api][DELETE]匹配的路由"},
		{name: "1.2 路径正确 请求方法为POST", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "POST"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api][POST]匹配的路由"},
		{name: "1.3 路径正确 请求方法为PUT", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "PUT"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api][PUT]匹配的路由"},
		{name: "1.4 路径正确 请求方法为PATCH", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "PATCH"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api][PATCH]匹配的路由"},
		{name: "1.5 路径不正确 请求方法正确", ctx: &mocks.TestContxt{Routerpath: "/api2", Method: "GET"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api2][GET]匹配的路由"},
		{name: "1.6 路径不正确 请求方法正不确", ctx: &mocks.TestContxt{Routerpath: "/api2", Method: "PATCH"}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api2][PATCH]匹配的路由"},
	}

	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, tt.serverConf, tt.meta)
		assert.PanicError(t, tt.wantError, func() { c.GetRouter() }, tt.name)
	}
}

func Test_rpath_GetRouter(t *testing.T) {
	confObj := mocks.NewConfBy("context_path_test1", "pathctx1") //构建对象
	confObj.API(":8080")
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx").Queue(queue.NewQueue("queue1", "/service1")).Queue(queue.NewQueue("queue2", "/service2"))
	confObj.CRON(c.WithMasterSlave(), c.WithTrace())
	confObj.Service.API.Add("/api", "/api", []string{"GET"}, api.WithEncoding("utf-8"))
	confObj.Service.Web.Add("/web", "/web", []string{"GET"}, api.WithEncoding("utf-8"))
	confObj.Service.WS.Add("/ws", "/ws", []string{"GET"}, api.WithEncoding("utf-8"))
	apiConf := confObj.GetAPIConf()   //获取配置
	webConf := confObj.GetWebConf()   //获取配置
	wsConf := confObj.GetWSConf()     //获取配置
	mqcConf := confObj.GetMQCConf()   //获取配置
	cronConf := confObj.GetCronConf() //获取配置

	tests := []struct {
		name       string
		ctx        context.IInnerContext
		serverConf app.IAPPConf
		meta       conf.IMeta
		want       *router.Router
		wantError  string
	}{
		{name: "1 api类型的router", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "GET"}, serverConf: apiConf, meta: conf.NewMeta(), want: &router.Router{Path: "/api", Encoding: "utf-8", Action: []string{"GET"}, Service: "/api"}},
		{name: "2 web类型的router", ctx: &mocks.TestContxt{Routerpath: "/web", Method: "GET"}, serverConf: webConf, meta: conf.NewMeta(), want: &router.Router{Path: "/web", Encoding: "utf-8", Action: []string{"GET"}, Service: "/web"}},
		{name: "3 ws类型的router", ctx: &mocks.TestContxt{Routerpath: "/ws", Method: "GET"}, serverConf: wsConf, meta: conf.NewMeta(), want: &router.Router{Path: "/ws", Encoding: "utf-8", Action: []string{"GET"}, Service: "/ws"}},
		{name: "4 cron类型的router", ctx: &mocks.TestContxt{Routerpath: "/mqc"}, serverConf: mqcConf, meta: conf.NewMeta(), want: &router.Router{Path: "/mqc", Encoding: "utf-8", Action: []string{}, Service: "/mqc"}},
		{name: "5 mqc类型的router", ctx: &mocks.TestContxt{Routerpath: "/cron"}, serverConf: cronConf, meta: conf.NewMeta(), want: &router.Router{Path: "/cron", Encoding: "utf-8", Action: []string{}, Service: "/cron"}},
	}

	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, tt.serverConf, tt.meta)
		got, err := c.GetRouter()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_rpath_GetEncoding(t *testing.T) {
	confObj := mocks.NewConfBy("context_path_test1", "pathctx1") //构建对象
	confObj.API(":8080")
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx").Queue(queue.NewQueue("queue1", "/service1")).Queue(queue.NewQueue("queue2", "/service2"))
	confObj.CRON(c.WithMasterSlave(), c.WithTrace())
	confObj.Service.API.Add("/api", "/api", []string{"GET"}, api.WithEncoding("utf-8"))
	confObj.Service.API.Add("/api2", "/api2", []string{"GET"}, api.WithEncoding("gbk"))
	confObj.Service.API.Add("/api3", "/api3", []string{"GET"})
	confObj.Service.Web.Add("/web", "/web", []string{"GET"}, api.WithEncoding("utf-8"))
	confObj.Service.Web.Add("/web2", "/web2", []string{"GET"}, api.WithEncoding("gbk"))
	confObj.Service.Web.Add("/web3", "/web3", []string{"GET"})
	confObj.Service.WS.Add("/ws", "/ws", []string{"GET"}, api.WithEncoding("utf-8"))
	confObj.Service.WS.Add("/ws2", "/ws2", []string{"GET"}, api.WithEncoding("gbk"))
	confObj.Service.WS.Add("/ws3", "/ws3", []string{"GET"})
	apiConf := confObj.GetAPIConf()   //获取配置
	webConf := confObj.GetWebConf()   //获取配置
	wsConf := confObj.GetWSConf()     //获取配置
	mqcConf := confObj.GetMQCConf()   //获取配置
	cronConf := confObj.GetCronConf() //获取配置

	tests := []struct {
		name       string
		ctx        context.IInnerContext
		serverConf app.IAPPConf
		meta       conf.IMeta
		want       string
	}{
		{name: "1.1 api类型,注册时设置encoding为utf-8", ctx: &mocks.TestContxt{Routerpath: "/api", Method: "GET"}, serverConf: apiConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "1.2 api类型,注册时设置encoding为gbk", ctx: &mocks.TestContxt{Routerpath: "/api2", Method: "GET"}, serverConf: apiConf, meta: conf.NewMeta(), want: "gbk"},
		{name: "1.3 api类型,注册时未设置encoding", ctx: &mocks.TestContxt{Routerpath: "/api3", Method: "GET"}, serverConf: apiConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "2.1 web类型,注册时设置encoding为utf-8", ctx: &mocks.TestContxt{Routerpath: "/web", Method: "GET"}, serverConf: webConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "2.2 web类型,注册时设置encoding为gbk", ctx: &mocks.TestContxt{Routerpath: "/web2", Method: "GET"}, serverConf: webConf, meta: conf.NewMeta(), want: "gbk"},
		{name: "2.3 web类型,注册时未设置encoding", ctx: &mocks.TestContxt{Routerpath: "/web3", Method: "GET"}, serverConf: webConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "3.1 ws类型,注册时设置encoding为utf-8", ctx: &mocks.TestContxt{Routerpath: "/ws", Method: "GET"}, serverConf: wsConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "3.2 ws类型,注册时设置encoding为gbk", ctx: &mocks.TestContxt{Routerpath: "/ws2", Method: "GET"}, serverConf: wsConf, meta: conf.NewMeta(), want: "gbk"},
		{name: "3.3 ws类型,注册时未设置encoding", ctx: &mocks.TestContxt{Routerpath: "/ws3", Method: "GET"}, serverConf: wsConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "4 cron类型获取默认值", ctx: &mocks.TestContxt{Routerpath: "/mqc", Method: "GET"}, serverConf: mqcConf, meta: conf.NewMeta(), want: "utf-8"},
		{name: "5 mqc类型获取默认值", ctx: &mocks.TestContxt{Routerpath: "/cron", Method: "GET"}, serverConf: cronConf, meta: conf.NewMeta(), want: "utf-8"},
	}
	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, tt.serverConf, tt.meta)
		got := c.GetEncoding()
		assert.Equal(t, tt.want, got, tt.name)

		//再次获取
		got2 := c.GetEncoding()
		assert.Equal(t, got, got2, tt.name)

	}
}
