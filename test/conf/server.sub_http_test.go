package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/proxy"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func newReady(t *testing.T, platName, sysName, serverType, clusterName string) (string, string, string, string, registry.IRegistry) {
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")
	return platName, sysName, serverType, clusterName, rgst
}

func Test_httpSub_GetHeaderConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName1", "sysName1", global.API, "cluster1")
	tests := []struct {
		name     string
		opts     []header.Option
		wantErr  bool
		wantConf header.Headers
	}{
		{name: "空task获取对象", opts: []header.Option{}, wantErr: true, wantConf: header.Headers{}},
		{name: "设置正确的task对象", opts: []header.Option{header.WithCrossDomain("localhost")}, wantErr: true, wantConf: header.New(header.WithCrossDomain("localhost"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Header(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		headerConf, err := gotS.GetHeaderConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, headerConf, tt.name+",conf")
	}
	//无法设置错误header数据
}

func Test_httpSub_GetJWTConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName2", "sysName2", global.API, "cluster2")
	var nilJwt *jwt.JWTAuth
	tests := []struct {
		name     string
		opts     []jwt.Option
		wantErr  bool
		wantConf *jwt.JWTAuth
	}{
		{name: "不设置jwt节点", opts: []jwt.Option{}, wantErr: true, wantConf: &jwt.JWTAuth{Disable: true}},
		{name: "设置错误jwt数据", opts: []jwt.Option{jwt.WithMode("错误数据")}, wantErr: false, wantConf: nilJwt},
		{name: "设置正确的jwt对象", opts: []jwt.Option{jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithAuthURL("1111")}, wantErr: true,
			wantConf: jwt.NewJWT(jwt.WithDisable(), jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithAuthURL("1111"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Jwt(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		jwtConf, err := gotS.GetJWTConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, jwtConf, tt.name+",conf")
	}
}

func Test_httpSub_GetMetricConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName3", "sysName3", global.API, "cluster3")
	var nilMetric *metric.Metric
	tests := []struct {
		name     string
		host     string
		db       string
		cron     string
		opts     []metric.Option
		wantErr  bool
		wantConf *metric.Metric
	}{
		{name: "不设置metric节点", opts: []metric.Option{}, wantErr: true, wantConf: &metric.Metric{Disable: true}},
		{name: "设置错误的metric节点", host: "168.0.111:8080", db: "1", cron: "cron1", opts: []metric.Option{metric.WithEnable(), metric.WithUPName("upnem", "1223456")}, wantErr: false, wantConf: nilMetric},
		{name: "设置正确的metric节点", host: "http://192.168.0.111:8080", db: "1", cron: "cron1", opts: []metric.Option{metric.WithEnable(), metric.WithUPName("upnem", "1223456")}, wantErr: true,
			wantConf: metric.New("http://192.168.0.111:8080", "1", "cron1", metric.WithEnable(), metric.WithUPName("upnem", "1223456"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Metric(tt.host, tt.db, tt.cron, tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		metricConf, err := gotS.GetMetricConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, metricConf, tt.name+",conf")
	}
}

func Test_httpSub_GetStaticConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName4", "sysName4", global.API, "cluster4")
	var nilstatic *static.Static
	tests := []struct {
		name     string
		opts     []static.Option
		wantErr  bool
		wantConf *static.Static
	}{
		{name: "不设置static节点", opts: []static.Option{}, wantErr: true, wantConf: &static.Static{FileMap: map[string]static.FileInfo{}, Disable: true}},
		{name: "设置错误的static节点", opts: []static.Option{static.WithRoot("错误的数据"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
			static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip")}, wantErr: false, wantConf: nilstatic},
		{name: "设置正确的static节点", opts: []static.Option{static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
			static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip")}, wantErr: true,
			wantConf: static.New(static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
				static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Static(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		staticConf, err := gotS.GetStaticConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, staticConf, tt.name+",conf")
	}
}

func Test_httpSub_GetRouterConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName5", "sysName5", global.API, "cluster5")
	tests := []struct {
		name     string
		wantErr  bool
		wantConf *router.Routers
	}{
		{name: "不设置router节点", wantErr: true, wantConf: router.NewRouters()},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080")
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		routerConf, err := gotS.GetRouterConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, routerConf, tt.name+",conf")
	}
}

func Test_httpSub_GetAPIKeyConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName6", "sysName6", global.API, "cluster6")
	var nilapikey *apikey.APIKeyAuth
	tests := []struct {
		name     string
		secert   string
		opts     []apikey.Option
		wantErr  bool
		wantConf *apikey.APIKeyAuth
	}{
		{name: "不设置apikey节点", opts: []apikey.Option{}, wantErr: true, wantConf: &apikey.APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}},
		{name: "设置错误的apikey节点", secert: "错误的数据", opts: []apikey.Option{apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2")}, wantErr: false, wantConf: nilapikey},
		{name: "设置正确的apikey节点", secert: "123456", opts: []apikey.Option{apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2")}, wantErr: true,
			wantConf: apikey.New("123456", apikey.WithDisable(), apikey.WithSHA256Mode(), apikey.WithExcludes("/p1/p2"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.APIKEY(tt.secert, tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		apikeyConf, err := gotS.GetAPIKeyConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, apikeyConf, tt.name+",conf")
	}
}

func Test_httpSub_GetRASConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName7", "sysName7", global.API, "cluster7")
	var nilRas *ras.RASAuth
	tests := []struct {
		name     string
		secert   string
		opts     []ras.Option
		wantErr  bool
		wantConf *ras.RASAuth
	}{
		{name: "不设置ras节点", opts: []ras.Option{}, wantErr: true, wantConf: &ras.RASAuth{Disable: true}},
		{name: "设置错误的ras节点", secert: "错误的数据", opts: []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
			ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable()))}, wantErr: false, wantConf: nilRas},
		{name: "设置正确的ras节点", secert: "123456", opts: []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
			ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable()))}, wantErr: true,
			wantConf: ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("service1", ras.WithRequest("/t1/t2"), ras.WithRequired("taofield"), ras.WithUIDAlias("userID"), ras.WithTimestampAlias("timespan"), ras.WithSignAlias("signname"),
				ras.WithCheckTimestamp(false), ras.WithDecryptName("duser"), ras.WithParam("key1", "v1"), ras.WithParam("key2", "v2"), ras.WithAuthDisable())))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Ras(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		rasConf, err := gotS.GetRASConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, rasConf, tt.name+",conf")
	}
}

func Test_httpSub_GetBasicConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName8", "sysName8", global.API, "cluster8")
	tests := []struct {
		name     string
		secert   string
		opts     []basic.Option
		wantErr  bool
		wantConf *basic.BasicAuth
	}{
		{name: "不设置basic节点", opts: []basic.Option{}, wantErr: true, wantConf: &basic.BasicAuth{Disable: true}},
		{name: "设置正确的basic节点", secert: "123456", opts: []basic.Option{basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1")}, wantErr: true,
			wantConf: basic.NewBasic(basic.WithDisable(), basic.WithUP("basicName", "basicPwd"), basic.WithExcludes("/basic/basic1"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Basic(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		basicConf, err := gotS.GetBasicConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, basicConf, tt.name+",conf")
	}
	//不能设置错误的basic节点
}

func Test_httpSub_GetRenderConf(t *testing.T) {

	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName9", "sysName9", global.API, "cluster9")
	var nilRender *render.Render
	tests := []struct {
		name     string
		opts     []render.Option
		wantErr  bool
		wantConf *render.Render
	}{
		{name: "不设置render节点", opts: []render.Option{}, wantErr: true, wantConf: &render.Render{Disable: true}},
		{name: "设置错误的render节点", opts: []render.Option{render.WithDisable(), render.WithTmplt("/path1", "success", render.WithContentType("tpltm1"))}, wantErr: false,
			wantConf: nilRender},
		{name: "设置正确的render节点", opts: []render.Option{render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1"))}, wantErr: true,
			wantConf: render.NewRender(render.WithDisable(), render.WithTmplt("/path1", "success", render.WithStatus("500"), render.WithContentType("tpltm1")))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Render(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		renderConf, err := gotS.GetRenderConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, renderConf, tt.name+",conf")
	}
}

func Test_httpSub_GetWhiteListConf(t *testing.T) {

	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName10", "sysName10", global.API, "cluster10")
	var nilWhitelist *whitelist.WhiteList
	tests := []struct {
		name     string
		opts     []whitelist.Option
		wantErr  bool
		wantConf *whitelist.WhiteList
	}{
		{name: "不设置whitelist节点", opts: []whitelist.Option{}, wantErr: true, wantConf: &whitelist.WhiteList{Disable: true}},
		{name: "设置错误的whitelist节点", opts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList("", []string{"192.168.0.101"}...))}, wantErr: false,
			wantConf: nilWhitelist},
		{name: "设置正确的whitelist节点", opts: []whitelist.Option{whitelist.WithIPList(whitelist.NewIPList("/t1/t2/*", []string{"192.168.0.101"}...))}, wantErr: true,
			wantConf: whitelist.New(whitelist.WithIPList(whitelist.NewIPList("/t1/t2/*", []string{"192.168.0.101"}...)))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.WhiteList(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		whitelistConf, err := gotS.GetWhiteListConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, whitelistConf, tt.name+",conf")
	}
}

func Test_httpSub_GetBlackListConf(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName11", "sysName11", global.API, "cluster11")
	var nilBlacklist *blacklist.BlackList
	tests := []struct {
		name     string
		opts     []blacklist.Option
		wantErr  bool
		wantConf *blacklist.BlackList
	}{
		{name: "不设置blacklist节点", opts: []blacklist.Option{}, wantErr: true, wantConf: &blacklist.BlackList{Disable: true}},
		{name: "设置错误的blacklist节点", opts: []blacklist.Option{blacklist.WithEnable()}, wantErr: true,
			wantConf: nilBlacklist},
		{name: "设置正确的blacklist节点", opts: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.121")}, wantErr: true,
			wantConf: blacklist.New(blacklist.WithEnable(), blacklist.WithIP("192.168.0.121"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.BlackList(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		_, err = gotS.GetBlackListConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
	}
}

func Test_httpSub_GetLimiter(t *testing.T) {
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName12", "sysName12", global.API, "cluster12")
	var nilLImiter *limiter.Limiter
	tests := []struct {
		name     string
		opts     []limiter.Option
		wantErr  bool
		wantConf *limiter.Limiter
	}{
		{name: "不设置limiter节点", opts: []limiter.Option{}, wantErr: true, wantConf: &limiter.Limiter{Disable: true}},
		{name: "设置错误的limiter节点", opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("错误数据1", "错误数据"), limiter.WithFallback(), limiter.WithReponse(200, "success")))}, wantErr: false,
			wantConf: nilLImiter},
		{name: "设置正确的limiter节点", opts: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("GET", "POST"), limiter.WithFallback(), limiter.WithReponse(200, "success")))}, wantErr: true,
			wantConf: limiter.New(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("path1", 1, limiter.WithMaxWait(3), limiter.WithAction("GET", "POST"), limiter.WithFallback(), limiter.WithReponse(200, "success"))))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Limit(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		limiterConf, err := gotS.GetLimiter()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, limiterConf, tt.name+",conf")
	}
}

func Test_httpSub_GetGray(t *testing.T) {
	type test struct {
		name     string
		opts     []proxy.Option
		wantErr  bool
		wantConf *proxy.Proxy
	}
	platName, sysName, serverType, clusterName, rgst := newReady(t, "platName13", "sysName13", global.API, "cluster13")
	var nilgray *proxy.Proxy
	tests := []test{
		{name: "不设置gray节点", opts: []proxy.Option{}, wantErr: true, wantConf: &proxy.Proxy{Disable: true}},
		{name: "设置正确的gray节点", opts: []proxy.Option{proxy.WithDisable(), proxy.WithFilter("Filter"), proxy.WithUPCluster("UPCluster")}, wantErr: true,
			wantConf: proxy.New(proxy.WithDisable(), proxy.WithFilter("Filter"), proxy.WithUPCluster("UPCluster"))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080")
		if len(tt.opts) > 0 {
			confN.Proxy(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		grayConf, err := gotS.GetProxy()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf.Disable, grayConf.Disable, "测试conf初始化,判断gary.Disable节点对象")
		assert.Equal(t, tt.wantConf.Filter, grayConf.Filter, "测试conf初始化,判断gary.Filter节点对象")
		assert.Equal(t, tt.wantConf.UPCluster, grayConf.UPCluster, "测试conf初始化,判断gary.UPCluster节点对象")
	}

	test1 := test{name: "设置错误的gray节点", opts: []proxy.Option{proxy.WithDisable()}, wantErr: false,
		wantConf: nilgray}
	confM := mocks.NewConfBy(platName, clusterName)
	confN := confM.API(":8080")
	confN.Proxy(test1.opts...)
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	grayConf, err := gotS.GetProxy()
	assert.Equal(t, test1.wantErr, err == nil, test1.name+",err")
	assert.Equal(t, test1.wantConf, grayConf, test1.name+",conf")
}
