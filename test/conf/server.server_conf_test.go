package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/proxy"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/conf/vars"
	 "github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewEmptyServerConf(t *testing.T) {
	platName := "platName1"
	sysName := "sysName1"
	serverType := global.API
	clusterName := "cluster1"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, false, err == nil, "测试conf初始化,没有设置主节点")

	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	varConf, err := vars.NewVarConf(platName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建var节点对象")
	assert.Equal(t, varConf, gotS.GetVarConf(), "测试conf初始化,判断Var节点对象")

	taskConf, err := gotS.GetCRONTaskConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败")
	assert.Equal(t, &task.Tasks{}, taskConf, "测试conf初始化,判断task节点对象")

	headerConf, err := gotS.GetHeaderConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取header对象失败")
	assert.Equal(t, header.Headers{}, headerConf, "测试conf初始化,判断header节点对象")

	jwtConf, err := gotS.GetJWTConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取jwt对象失败")
	assert.Equal(t, &jwt.JWTAuth{Disable: true}, jwtConf, "测试conf初始化,判断jwt节点对象")

	metricConf, err := gotS.GetMetricConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取metric对象失败")
	assert.Equal(t, &metric.Metric{Disable: true}, metricConf, "测试conf初始化,判断metric节点对象")

	staticConf, err := gotS.GetStaticConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取static对象失败")
	assert.Equal(t, &static.Static{FileMap: map[string]static.FileInfo{}, Disable: true}, staticConf, "测试conf初始化,判断static节点对象")

	routerConf, err := gotS.GetRouterConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取router对象失败")
	assert.Equal(t, router.NewRouters(), routerConf, "测试conf初始化,判断router节点对象")

	apikeyConf, err := gotS.GetAPIKeyConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取apikey对象失败")
	assert.Equal(t, &apikey.APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}, apikeyConf, "测试conf初始化,判断apikey节点对象")

	authRasConf, err := gotS.GetRASConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取authRas对象失败")
	assert.Equal(t, &ras.RASAuth{Disable: true}, authRasConf, "测试conf初始化,判断authRas节点对象")

	basicConf, err := gotS.GetBasicConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取basic对象失败")
	assert.Equal(t, &basic.BasicAuth{Disable: true}, basicConf, "测试conf初始化,判断basic节点对象")

	renderConf, err := gotS.GetRenderConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取render对象失败")
	assert.Equal(t, &render.Render{Disable: true}, renderConf, "测试conf初始化,判断render节点对象")

	whiteListConf, err := gotS.GetWhiteListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取whiteList对象失败")
	assert.Equal(t, &whitelist.WhiteList{Disable: true}, whiteListConf, "测试conf初始化,判断whiteList节点对象")

	blackListConf, err := gotS.GetBlackListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取blackList对象失败")
	assert.Equal(t, &blacklist.BlackList{Disable: true}, blackListConf, "测试conf初始化,判断blackList节点对象")

	limiterConf, err := gotS.GetLimiterConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取limiter对象失败")
	assert.Equal(t, &limiter.Limiter{Disable: true}, limiterConf, "测试conf初始化,判断limiter节点对象")

	garyConf, err := gotS.GetProxyConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取gary对象失败")
	assert.Equal(t, &proxy.Proxy{Disable: true}, garyConf, "测试conf初始化,判断gary节点对象")

	_, err = gotS.GetMQCMainConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取mqc对象失败")

	queuesObj, err := gotS.GetMQCQueueConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败")
	assert.Equal(t, &queue.Queues{}, queuesObj, "测试conf初始化,判断queues节点对象")

	layoutObj, err := gotS.GetRLogConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取layout对象失败")
	assert.Equal(t, &rlog.Layout{Layout: rlog.DefaultLayout, Disable: true}, layoutObj, "测试conf初始化,判断layout节点对象")
}

func TestNewAPIServerConf(t *testing.T) {
	platName := "platName2"
	sysName := "sysName2"
	serverType := global.API
	clusterName := "cluster2"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	confN := confM.API(":8080", api.WithDisable(), api.WithTrace(), api.WithDNS("ip1"), api.WithHeaderReadTimeout(10), api.WithTimeout(11, 11))
	confN.APIKEY("123456", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))
	confN.Basic(basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1"))
	confN.BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.121"))
	//confN.Proxy(proxy.WithDisable(), proxy.WithFilter("Filter"), proxy.WithUPCluster("UPCluster"))
	confN.Header(header.WithCrossDomain("localhost"))
	confN.Jwt(jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithAuthURL("1111"))
	confN.Limit(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithFallback(), limiter.WithReponse(200, "success"))))
	confN.Metric("http://192.168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))
	confN.Ras(ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
		ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))
	//confN.Render(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1")))
	confN.Static(static.WithRoot("./test"), static.WithHomePage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	confN.WhiteList(whitelist.WithIPList(whitelist.NewIPList([]string{"/t1/t2/*"}, []string{"192.168.0.101"}...)))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	varConf, err := vars.NewVarConf(platName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建var节点对象")
	assert.Equal(t, varConf, gotS.GetVarConf(), "测试conf初始化,判断Var节点对象")

	taskConf, err := gotS.GetCRONTaskConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败")
	assert.Equal(t, &task.Tasks{}, taskConf, "测试conf初始化,判断task节点对象")

	headerConf, err := gotS.GetHeaderConf()
	headerC := header.New(header.WithCrossDomain("localhost"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取header对象失败")
	assert.Equal(t, headerC, headerConf, "测试conf初始化,判断header节点对象")

	jwtConf, err := gotS.GetJWTConf()
	jwtC := jwt.NewJWT(jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithAuthURL("1111"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取jwt对象失败")
	assert.Equal(t, jwtC, jwtConf, "测试conf初始化,判断jwt节点对象")

	metricConf, err := gotS.GetMetricConf()
	metricC := metric.New("http://192.168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取metric对象失败")
	assert.Equal(t, metricC, metricConf, "测试conf初始化,判断metric节点对象")

	staticConf, err := gotS.GetStaticConf()
	staticC := static.New(static.WithRoot("./test"), static.WithHomePage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取static对象失败")
	assert.Equal(t, staticC, staticConf, "测试conf初始化,判断static节点对象")

	routerConf, err := gotS.GetRouterConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取router对象失败")
	assert.Equal(t, router.NewRouters(), routerConf, "测试conf初始化,判断router节点对象")

	apikeyConf, err := gotS.GetAPIKeyConf()
	apikeyC := apikey.New("123456", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取apikey对象失败")
	assert.Equal(t, apikeyC, apikeyConf, "测试conf初始化,判断apikey节点对象")

	authRasConf, err := gotS.GetRASConf()
	rasC := ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
		ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取authRas对象失败")
	assert.Equal(t, rasC, authRasConf, "测试conf初始化,判断authRas节点对象")

	basicConf, err := gotS.GetBasicConf()
	basicC := basic.NewBasic(basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取basic对象失败")
	assert.Equal(t, basicC, basicConf, "测试conf初始化,判断basic节点对象")

	/*
		renderConf, err := gotS.GetRenderConf()
		renderC := render.NewRender(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1")))
		assert.Equal(t, true, err == nil, "测试conf初始化,获取render对象失败")
		assert.Equal(t, renderC, renderConf, "测试conf初始化,判断render节点对象")
	*/

	whiteListConf, err := gotS.GetWhiteListConf()
	whiteC := whitelist.New(whitelist.WithIPList(whitelist.NewIPList([]string{"/t1/t2/*"}, []string{"192.168.0.101"}...)))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取whiteList对象失败")
	assert.Equal(t, whiteC, whiteListConf, "测试conf初始化,判断whiteList节点对象")

	blackListConf, err := gotS.GetBlackListConf()
	blackC := blacklist.New(blacklist.WithEnable(), blacklist.WithIP("192.168.0.121"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取blackList对象失败")
	assert.Equal(t, blackC, blackListConf, "测试conf初始化,判断blackList节点对象")

	limiterConf, err := gotS.GetLimiterConf()
	limiterC := limiter.New(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithFallback(), limiter.WithReponse(200, "success"))))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取limiter对象失败")
	assert.Equal(t, limiterC, limiterConf, "测试conf初始化,判断limiter节点对象")

	/*
		garyConf, err := gotS.GetProxyConf()
		garyC := proxy.New(proxy.WithDisable(), proxy.WithFilter("Filter"), proxy.WithUPCluster("UPCluster"))
		assert.Equal(t, true, err == nil, "测试conf初始化,获取gary对象失败")
		assert.Equal(t, garyC.Disable, garyConf.Disable, "测试conf初始化,判断gary.Disable节点对象")
		assert.Equal(t, garyC.Filter, garyConf.Filter, "测试conf初始化,判断gary.Filter节点对象")
		assert.Equal(t, garyC.UPCluster, garyConf.UPCluster, "测试conf初始化,判断gary.UPCluster节点对象")

	*/

	_, err = gotS.GetMQCMainConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取mqc对象失败")

	queuesObj, err := gotS.GetMQCQueueConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败")
	assert.Equal(t, &queue.Queues{}, queuesObj, "测试conf初始化,判断queues节点对象")

	layoutObj, err := gotS.GetRLogConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取layout对象失败")
	assert.Equal(t, &rlog.Layout{Layout: rlog.DefaultLayout, Disable: true}, layoutObj, "测试conf初始化,判断layout节点对象")
}

func TestNewRPCServerConf(t *testing.T) {
	platName := "platName3"
	sysName := "sysName3"
	serverType := global.RPC
	clusterName := "cluster3"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, false, err == nil, "测试conf初始化,没有设置主节点")

	confM.RPC(":8081")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	varConf, err := vars.NewVarConf(platName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建var节点对象")
	assert.Equal(t, varConf, gotS.GetVarConf(), "测试conf初始化,判断Var节点对象")
}

func TestNewMQCServerConf(t *testing.T) {
	platName := "platName4"
	sysName := "sysName4"
	serverType := global.MQC
	clusterName := "cluster4"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, false, err == nil, "测试conf初始化,没有设置主节点")

	confN := confM.MQC("redis://11")
	confN.Queue(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	mqcConf, err := gotS.GetMQCMainConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取mqc对象失败")
	mqcC := mqc.New("redis://11")
	assert.Equal(t, mqcC, mqcConf, "测试conf初始化,判断mqc节点对象")

	//@todo queue 需要在应用启动后才会发布到注册中心
	// queuesObj, err := gotS.GetMQCQueueConf()
	// queueC := queue.NewQueues(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	// assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败")
	// assert.Equal(t, queueC, queuesObj, "测试conf初始化,判断queues节点对象")
}

func TestNewCRONServerConf(t *testing.T) {
	platName := "platName5"
	sysName := "sysName5"
	serverType := global.CRON
	clusterName := "cluster5"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, false, err == nil, "测试conf初始化,没有设置主节点")

	confN := confM.CRON(cron.WithTrace(), cron.WithDisable(), cron.WithSharding(1))
	confN.Task(task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable()))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	varConf, err := vars.NewVarConf(platName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建var节点对象")
	assert.Equal(t, varConf, gotS.GetVarConf(), "测试conf初始化,判断Var节点对象")

	taskConf, err := gotS.GetCRONTaskConf()
	taskC := task.NewTasks(task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable()))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败")
	assert.Equal(t, taskC, taskConf, "测试conf初始化,判断task节点对象")
}

func TestNewVARServerConf(t *testing.T) {
	platName := "platName6"
	sysName := "sysName6"
	serverType := global.MQC
	clusterName := "cluster6"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, false, err == nil, "测试conf初始化,没有设置主节点")

	confM.MQC("redis://11").Queue(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	confM.Vars().Queue().Redis("redis", queueredis.New(queueredis.WithAddrs("192.168.0.1")))

	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	mainConf, err := server.NewServerConf(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,新建主节点对象")
	assert.Equal(t, mainConf.GetMainConf(), gotS.GetServerConf().GetMainConf(), "测试conf初始化,判断主节点对象")

	mqcConf, err := gotS.GetMQCMainConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取mqc对象失败")
	mqcC := mqc.New("redis://11")
	assert.Equal(t, mqcC, mqcConf, "测试conf初始化,判断mqc节点对象")

	//@todo OnReady无法调用  不能注册queue
	// queuesObj, err := gotS.GetMQCQueueConf()
	// queueC := queue.NewQueues(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	// assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败")
	// assert.Equal(t, queueC, queuesObj, "测试conf初始化,判断queues节点对象")
}
