package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_httpSub_GetHeaderConf(t *testing.T) {
	platName := "platName1"
	sysName := "sysName1"
	serverType := global.API
	clusterName := "cluster1"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置header节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	headerConf, err := gotS.GetHeaderConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取header对象失败")
	assert.Equal(t, header.Headers{}, headerConf, "测试conf初始化,判断header节点对象")

	//无法设置错误header数据

	//设置正确的header
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Header(header.WithCrossDomain("localhost"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	headerConf, err = gotS.GetHeaderConf()
	headerC := header.New(header.WithCrossDomain("localhost"))
	assert.Equal(t, true, err == nil, "测试conf初始化,获取header对象失败1")
	assert.Equal(t, headerC, headerConf, "测试conf初始化,判断header节点对象1")
}

func Test_httpSub_GetJWTConf(t *testing.T) {
	platName := "platName2"
	sysName := "sysName2"
	serverType := global.API
	clusterName := "cluster2"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置jwt节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	jwtConf, err := gotS.GetJWTConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取jwt对象失败")
	assert.Equal(t, &jwt.JWTAuth{Disable: true}, jwtConf, "测试conf初始化,判断jwt节点对象")

	//设置错误jwt数据
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Jwt(jwt.WithMode("错误数据"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	jwtConf, err = gotS.GetJWTConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取jwt对象失败1")
	var nilJwt *jwt.JWTAuth
	assert.Equal(t, nilJwt, jwtConf, "测试conf初始化,判断jwt节点对象1")

	//设置正确的header
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Jwt(jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithRedirect("1111"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	jwtConf, err = gotS.GetJWTConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取jwt对象失败2")
	jwtC := jwt.NewJWT(jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithRedirect("1111"))
	assert.Equal(t, jwtC, jwtConf, "测试conf初始化,判断jwt节点对象2")
}

func Test_httpSub_GetMetricConf(t *testing.T) {
	platName := "platName3"
	sysName := "sysName3"
	serverType := global.API
	clusterName := "cluster3"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置metric节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	metricConf, err := gotS.GetMetricConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取metric对象失败")
	assert.Equal(t, &metric.Metric{Disable: true}, metricConf, "测试conf初始化,判断metric节点对象")

	//设置错误的metric节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Metric("168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	metricConf, err = gotS.GetMetricConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取metric对象失败1")
	var nilMetric *metric.Metric
	assert.Equal(t, nilMetric, metricConf, "测试conf初始化,判断metric节点对象1")

	//设置正确的metric节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Metric("http://192.168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	metricConf, err = gotS.GetMetricConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取metric对象失败2")
	metricC := metric.New("http://192.168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))
	assert.Equal(t, metricC, metricConf, "测试conf初始化,判断metric节点对象2")

}

func Test_httpSub_GetStaticConf(t *testing.T) {
	platName := "platName4"
	sysName := "sysName4"
	serverType := global.API
	clusterName := "cluster4"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置Static节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	staticConf, err := gotS.GetStaticConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取static对象失败")
	assert.Equal(t, &static.Static{FileMap: map[string]static.FileInfo{}, Disable: true}, staticConf, "测试conf初始化,判断static节点对象")

	//设置错误的Static节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Static(static.WithRoot("错误的数据"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	staticConf, err = gotS.GetStaticConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取static对象失败1")
	var nilstatic *static.Static
	assert.Equal(t, nilstatic, staticConf, "测试conf初始化,判断static节点对象1")

	//设置正确的Static节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Static(static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	staticConf, err = gotS.GetStaticConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取static对象失败2")
	staticC := static.New(static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	assert.Equal(t, staticC, staticConf, "测试conf初始化,判断static节点对象2")
}

func Test_httpSub_GetRouterConf(t *testing.T) {
	platName := "platName5"
	sysName := "sysName5"
	serverType := global.API
	clusterName := "cluster5"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	routerConf, err := gotS.GetRouterConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取router对象失败")
	assert.Equal(t, router.NewRouters(), routerConf, "测试conf初始化,判断router节点对象")

	//不能对router进行进行初始化   所以只能测试未设置节点的情况
}

func Test_httpSub_GetAPIKeyConf(t *testing.T) {
	platName := "platName6"
	sysName := "sysName6"
	serverType := global.API
	clusterName := "cluster6"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	apikeyConf, err := gotS.GetAPIKeyConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取apikey对象失败")
	assert.Equal(t, &apikey.APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}, apikeyConf, "测试conf初始化,判断apikey节点对象")

	//设置错误的apikye节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").APIKEY("错误的数据", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	apikeyConf, err = gotS.GetAPIKeyConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取apikey对象失败")
	var nilapikey *apikey.APIKeyAuth
	assert.Equal(t, nilapikey, apikeyConf, "测试conf初始化,判断apikey节点对象")

	//设置正确的apikey节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").APIKEY("123456", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	apikeyConf, err = gotS.GetAPIKeyConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取apikey对象失败")
	apikeyC := apikey.New("123456", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))
	assert.Equal(t, apikeyC, apikeyConf, "测试conf初始化,判断apikey节点对象")
}

func Test_httpSub_GetRASConf(t *testing.T) {
	platName := "platName7"
	sysName := "sysName7"
	serverType := global.API
	clusterName := "cluster7"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	rasConf, err := gotS.GetRASConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取ras对象失败")
	assert.Equal(t, &ras.RASAuth{Disable: true}, rasConf, "测试conf初始化,判断ras节点对象")

	//设置错误的ras节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Ras(ras.WithDisable(), ras.WithAuths(ras.New("", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
		ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	rasConf, err = gotS.GetRASConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取ras对象失败1")
	var nilRas *ras.RASAuth
	assert.Equal(t, nilRas, rasConf, "测试conf初始化,判断ras节点对象1")

	//设置正确的ras节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Ras(ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
		ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	rasConf, err = gotS.GetRASConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取ras对象失败2")
	rasC := ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
		ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))
	assert.Equal(t, rasC, rasConf, "测试conf初始化,判断ras节点对象2")
}

func Test_httpSub_GetBasicConf(t *testing.T) {
	platName := "platName8"
	sysName := "sysName8"
	serverType := global.API
	clusterName := "cluster8"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	basicConf, err := gotS.GetBasicConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取basic对象失败")
	assert.Equal(t, &basic.BasicAuth{Disable: true}, basicConf, "测试conf初始化,判断basic节点对象")

	//不能设置错误的basic节点

	//设置正确的basic节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Basic(basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	basicConf, err = gotS.GetBasicConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取basic对象失败1")
	basicC := basic.NewBasic(basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1"))
	assert.Equal(t, basicC, basicConf, "测试conf初始化,判断basic节点对象1")

}

func Test_httpSub_GetRenderConf(t *testing.T) {
	platName := "platName9"
	sysName := "sysName9"
	serverType := global.API
	clusterName := "cluster9"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	renderConf, err := gotS.GetRenderConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取render对象失败")
	assert.Equal(t, &render.Render{Disable: true}, renderConf, "测试conf初始化,判断render节点对象")

	//设置错误的render节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Render(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithContentType("tpltm1")))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	renderConf, err = gotS.GetRenderConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取render对象失败1")
	var nilRender *render.Render
	assert.Equal(t, nilRender, renderConf, "测试conf初始化,判断render节点对象1")

	//设置错误的render节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Render(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1")))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	renderConf, err = gotS.GetRenderConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取render对象失败2")
	renderC := render.NewRender(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1")))
	assert.Equal(t, renderC, renderConf, "测试conf初始化,判断render节点对象2")
}

func Test_httpSub_GetWhiteListConf(t *testing.T) {
	platName := "platName10"
	sysName := "sysName10"
	serverType := global.API
	clusterName := "cluster10"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	whitelistConf, err := gotS.GetWhiteListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取whitelist对象失败")
	assert.Equal(t, &whitelist.WhiteList{Disable: true}, whitelistConf, "测试conf初始化,判断whitelist节点对象")

	//设置错误的whitelist节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").WhiteList(whitelist.WithIPList(whitelist.NewIPList("", []string{"192.168.0.101"}...)))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	whitelistConf, err = gotS.GetWhiteListConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取whitelist对象失败1")
	var nilWhitelist *whitelist.WhiteList
	assert.Equal(t, nilWhitelist, whitelistConf, "测试conf初始化,判断whitelist节点对象1")

	//设置正确的whitelist节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").WhiteList(whitelist.WithIPList(whitelist.NewIPList("/t1/t2/*", []string{"192.168.0.101"}...)))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	whitelistConf, err = gotS.GetWhiteListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取whitelist对象失败2")
	whiteC := whitelist.New(whitelist.WithIPList(whitelist.NewIPList("/t1/t2/*", []string{"192.168.0.101"}...)))
	assert.Equal(t, whiteC, whitelistConf, "测试conf初始化,判断whitelist节点对象2")
}

func Test_httpSub_GetBlackListConf(t *testing.T) {
	platName := "platName11"
	sysName := "sysName11"
	serverType := global.API
	clusterName := "cluster11"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	blacklistConf, err := gotS.GetBlackListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取blacklist对象失败")
	assert.Equal(t, &blacklist.BlackList{Disable: true}, blacklistConf, "测试conf初始化,判断blacklist节点对象")

	//设置错误的blacklist节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").BlackList(blacklist.WithEnable())
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	blacklistConf, err = gotS.GetBlackListConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取blacklist对象失败1")
	var nilBlacklist *blacklist.BlackList
	assert.Equal(t, nilBlacklist, blacklistConf, "测试conf初始化,判断blacklist节点对象1")

	//设置正确的blacklist节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.121"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	blacklistConf, err = gotS.GetBlackListConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取blacklist对象失败2")
	blackC := blacklist.New(blacklist.WithEnable(), blacklist.WithIP("192.168.0.121"))
	assert.Equal(t, blackC, blacklistConf, "测试conf初始化,判断blacklist节点对象2")
}

func Test_httpSub_GetLimiter(t *testing.T) {
	platName := "platName12"
	sysName := "sysName12"
	serverType := global.API
	clusterName := "cluster12"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	limiterConf, err := gotS.GetLimiter()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取limiter对象失败")
	assert.Equal(t, &limiter.Limiter{Disable: true}, limiterConf, "测试conf初始化,判断limiter节点对象")

	//设置错误的limiter节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Limit(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("错误数据1", "错误数据"), limiter.WithFallback(), limiter.WithReponse(200, "success"))))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	limiterConf, err = gotS.GetLimiter()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取limiter对象失败1")
	var nilLImiter *limiter.Limiter
	assert.Equal(t, nilLImiter, limiterConf, "测试conf初始化,判断limiter节点对象1")

	//设置正确的limiter节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Limit(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("GET", "POST"), limiter.WithFallback(), limiter.WithReponse(200, "success"))))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	limiterConf, err = gotS.GetLimiter()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取limiter对象失败2")
	limiterC := limiter.New(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("GET", "POST"), limiter.WithFallback(), limiter.WithReponse(200, "success"))))
	assert.Equal(t, limiterC, limiterConf, "测试conf初始化,判断limiter节点对象2")
}

func Test_httpSub_GetGray(t *testing.T) {
	platName := "platName13"
	sysName := "sysName13"
	serverType := global.API
	clusterName := "cluster13"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//不设置节点
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	grayConf, err := gotS.GetGray()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取gray对象失败")
	assert.Equal(t, &gray.Gray{Disable: true}, grayConf, "测试conf初始化,判断gray节点对象")

	//设置错误的gray节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Gray(gray.WithDisable())
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	grayConf, err = gotS.GetGray()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取gray对象失败1")
	var nilgray *gray.Gray
	assert.Equal(t, nilgray, grayConf, "测试conf初始化,判断gray节点对象1")

	//设置正确的gray节点
	confM = mocks.NewConfBy(platName, clusterName)
	confM.API(":8080").Gray(gray.WithDisable(), gray.WithFilter("Filter"), gray.WithUPCluster("UPCluster"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	grayConf, err = gotS.GetGray()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取gray对象失败")
	garyC := gray.New(gray.WithDisable(), gray.WithFilter("Filter"), gray.WithUPCluster("UPCluster"))
	assert.Equal(t, garyC.Disable, grayConf.Disable, "测试conf初始化,判断gary.Disable节点对象")
	assert.Equal(t, garyC.Filter, grayConf.Filter, "测试conf初始化,判断gary.Filter节点对象")
	assert.Equal(t, garyC.UPCluster, grayConf.UPCluster, "测试conf初始化,判断gary.UPCluster节点对象")
}
