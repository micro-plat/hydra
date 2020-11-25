package creator

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/services"
)

func Test_newHTTP(t *testing.T) {
	type args struct {
		tp      string
		address string
		f       func(string) *services.ORouter
		opts    []api.Option
	}
	tests := []struct {
		name   string
		args   args
		repeat *args
		want   *httpBuilder
	}{
		{name: "1. 无option配置,初始化", args: args{tp: "tp1", address: ":1122", f: func(string) *services.ORouter { return nil }, opts: []api.Option{}},
			want: &httpBuilder{tp: "tp1", fnGetRouter: func(string) *services.ORouter { return nil }, CustomerBuilder: map[string]interface{}{"main": api.New(":1122")}}},
		{name: "2. 设置option配置,初始化", args: args{tp: "tp1", address: ":1122", f: func(string) *services.ORouter { return nil }, opts: []api.Option{api.WithDisable()}},
			want: &httpBuilder{tp: "tp1", fnGetRouter: func(string) *services.ORouter { return nil }, CustomerBuilder: map[string]interface{}{"main": api.New(":1122", api.WithDisable())}}},
		{name: "3. 重复设置节点", args: args{tp: "tp1", address: ":1122", f: func(string) *services.ORouter { return nil }, opts: []api.Option{api.WithDisable()}},
			repeat: &args{tp: "tp2", address: ":1123", f: func(string) *services.ORouter { return nil }, opts: []api.Option{api.WithDisable()}},
			want:   &httpBuilder{tp: "tp2", fnGetRouter: func(string) *services.ORouter { return nil }, CustomerBuilder: map[string]interface{}{"main": api.New(":1123", api.WithDisable())}}},
	}
	for _, tt := range tests {
		got := newHTTP(tt.args.tp, tt.args.address, tt.args.f, tt.args.opts...)
		if tt.repeat != nil {
			got = newHTTP(tt.repeat.tp, tt.repeat.address, tt.repeat.f, tt.repeat.opts...)
		}
		assert.Equal(t, tt.want.tp, got.tp, tt.name+",tp")
		assert.Equal(t, tt.want.fnGetRouter(""), got.fnGetRouter(""), tt.name+",fnGetRouter")
		assert.Equal(t, tt.want.CustomerBuilder, got.CustomerBuilder, tt.name+",CustomerBuilder")
	}
}

func Test_httpBuilder_Load(t *testing.T) {
	tests := []struct {
		name string
		c    httpBuilder
		want CustomerBuilder
	}{
		{name: "1. 空路由,加载http路由配置", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
			return services.NewORouter()
		}, CustomerBuilder: make(map[string]interface{})}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "2. 重复路由,加载http路由配置", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
			r := services.NewORouter()
			r.Add("path1", "service1", []string{"get"})
			r.Add("path1", "service1", []string{"get"})
			return r
		}, CustomerBuilder: make(map[string]interface{})}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "3. 正常路由,加载http路由配置", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
			r := services.NewORouter()
			r.Add("path1", "service1", []string{"get"})
			return r
		}, CustomerBuilder: make(map[string]interface{})}, want: CustomerBuilder{"router": router.NewRouters()}},
	}
	for _, tt := range tests {
		defer func() {
			if e := recover(); e != nil {
				assert.Equal(t, true, strings.Contains(types.GetString(e), "重复注册的服务"), tt.name)
			}
		}()
		tt.c.Load()
		assert.Equal(t, tt.want, tt.c.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Jwt(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []jwt.Option
		repeat []jwt.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认jwt配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456")}, want: CustomerBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("123456"))}},
		{name: "2. 初始化自定义jwt配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456"), jwt.WithExcludes("/taews/ssss")},
			want: CustomerBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("123456"), jwt.WithExcludes("/taews/ssss"))}},
		{name: "3. 重复初始化jwt配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456")},
			repeat: []jwt.Option{jwt.WithSecret("3232323"), jwt.WithExcludes("/taews/xxx")}, want: CustomerBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("3232323"), jwt.WithExcludes("/taews/xxx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Jwt(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Jwt(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_APIKEY(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		fields *httpBuilder
		args   []apikey.Option
		repeat []apikey.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认apikey对象", secret: "", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []apikey.Option{}, want: CustomerBuilder{"auth/apikey": apikey.New("")}},
		{name: "2. 初始化自定义apikey对象", secret: "123456", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []apikey.Option{apikey.WithSHA256Mode(), apikey.WithDisable()}, want: CustomerBuilder{"auth/apikey": apikey.New("123456", apikey.WithSHA256Mode(), apikey.WithDisable())}},
		{name: "3. 重复初始化apikey对象", secret: "123456", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []apikey.Option{apikey.WithSHA256Mode(), apikey.WithDisable()},
			repeat: []apikey.Option{apikey.WithSHA1Mode(), apikey.WithSecret("xxxxxx")},
			want:   CustomerBuilder{"auth/apikey": apikey.New("123456", apikey.WithSHA1Mode(), apikey.WithSecret("xxxxxx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.APIKEY(tt.secret, tt.args...)
		if tt.repeat != nil {
			got = tt.fields.APIKEY(tt.secret, tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Basic(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []basic.Option
		repeat []basic.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认basic配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []basic.Option{}, want: CustomerBuilder{"auth/basic": basic.NewBasic()}},
		{name: "2. 初始化自定义basic配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []basic.Option{basic.WithDisable(), basic.WithExcludes("11s")}, want: CustomerBuilder{"auth/basic": basic.NewBasic(basic.WithDisable(), basic.WithExcludes("11s"))}},
		{name: "3. 重复初始化自定义basic配置对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []basic.Option{basic.WithDisable(), basic.WithExcludes("11s")},
			repeat: []basic.Option{basic.WithUP("ssss", "bbbbb"), basic.WithExcludes("11xx")},
			want:   CustomerBuilder{"auth/basic": basic.NewBasic(basic.WithUP("ssss", "bbbbb"), basic.WithExcludes("11xx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Basic(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Basic(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_WhiteList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []whitelist.Option
		repeat []whitelist.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认whitelist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []whitelist.Option{}, want: CustomerBuilder{"acl/white.list": whitelist.New()}},
		{name: "2. 初始化自定义whitelist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []whitelist.Option{whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"}))},
			want: CustomerBuilder{"acl/white.list": whitelist.New(whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"})))}},
		{name: "3. 重复初始化自定义whitelist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []whitelist.Option{whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"}))},
			repeat: []whitelist.Option{whitelist.WithEnable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request1"}))},
			want:   CustomerBuilder{"acl/white.list": whitelist.New(whitelist.WithEnable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request1"})))}},
	}
	for _, tt := range tests {
		got := tt.fields.WhiteList(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.WhiteList(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_BlackList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []blacklist.Option
		repeat []blacklist.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认blacklist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []blacklist.Option{}, want: CustomerBuilder{"acl/black.list": blacklist.New()}},
		{name: "2. 初始化自定义blacklist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.101")},
			want: CustomerBuilder{"acl/black.list": blacklist.New(blacklist.WithDisable(), blacklist.WithIP("192.168.0.101"))}},
		{name: "3. 重复初始化自定义blacklist对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.101")},
			repeat: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.111")},
			want:   CustomerBuilder{"acl/black.list": blacklist.New(blacklist.WithEnable(), blacklist.WithIP("192.168.0.111"))}},
	}
	for _, tt := range tests {
		got := tt.fields.BlackList(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.BlackList(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Ras(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []ras.Option
		repeat []ras.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认ras对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []ras.Option{}, want: CustomerBuilder{"auth/ras": ras.NewRASAuth()}},
		{name: "2. 初始化自定义ras对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1")))},
			want: CustomerBuilder{"auth/ras": ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1"))))}},
		{name: "3. 重复初始化自定义ras对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1")))},
			repeat: []ras.Option{ras.WithEnable(), ras.WithAuths(ras.New("server2", ras.WithAuthDisable(), ras.WithRequest("patch2")))},
			want:   CustomerBuilder{"auth/ras": ras.NewRASAuth(ras.WithEnable(), ras.WithAuths(ras.New("server2", ras.WithAuthDisable(), ras.WithRequest("patch2"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Ras(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Ras(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Header(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []header.Option
		repeat []header.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认header对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []header.Option{}, want: CustomerBuilder{"header": header.New()}},
		{name: "2. 初始化自定义header对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []header.Option{header.WithAllowMethods("get", "put")}, want: CustomerBuilder{"header": header.New(header.WithAllowMethods("get", "put"))}},
		{name: "3. 重复初始化自定义header对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []header.Option{header.WithAllowMethods("get", "put")},
			repeat: []header.Option{header.WithCrossDomain("www.baidu.com"), header.WithAllowMethods("get", "put")},
			want:   CustomerBuilder{"header": header.New(header.WithCrossDomain("www.baidu.com"), header.WithAllowMethods("get", "put"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Header(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Header(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Metric(t *testing.T) {
	type args struct {
		host string
		db   string
		cron string
		opts []metric.Option
	}
	tests := []struct {
		name   string
		fields *httpBuilder
		args   args
		repeat *args
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认metric对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: args{host: "host1", db: "db1", cron: "cron1", opts: []metric.Option{}}, want: CustomerBuilder{"metric": metric.New("host1", "db1", "cron1")}},
		{name: "2. 初始化自定义metric对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: args{host: "host1", db: "db1", cron: "cron1", opts: []metric.Option{metric.WithDisable(), metric.WithUPName("name", "pwd")}},
			want: CustomerBuilder{"metric": metric.New("host1", "db1", "cron1", metric.WithDisable(), metric.WithUPName("name", "pwd"))}},
		{name: "3. 重复初始化自定义metric对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   args{host: "host1", db: "db1", cron: "cron1", opts: []metric.Option{metric.WithDisable(), metric.WithUPName("name", "pwd")}},
			repeat: &args{host: "host2", db: "db2", cron: "cron2", opts: []metric.Option{metric.WithEnable(), metric.WithUPName("xxxx", "pppp")}},
			want:   CustomerBuilder{"metric": metric.New("host2", "db2", "cron2", metric.WithEnable(), metric.WithUPName("xxxx", "pppp"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Metric(tt.args.host, tt.args.db, tt.args.cron, tt.args.opts...)
		if tt.repeat != nil {
			got = tt.fields.Metric(tt.repeat.host, tt.repeat.db, tt.repeat.cron, tt.repeat.opts...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Static(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []static.Option
		repeat []static.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认static对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []static.Option{}, want: CustomerBuilder{"static": static.New()}},
		{name: "2. 初始化自定义static对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []static.Option{static.WithDisable(), static.WithArchive("./sssss")},
			want: CustomerBuilder{"static": static.New(static.WithDisable(), static.WithArchive("./sssss"))}},
		{name: "3. 重复初始化自定义static对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []static.Option{static.WithDisable(), static.WithArchive("./sssss")},
			repeat: []static.Option{static.WithEnable(), static.WithArchive("./xxxx"), static.WithExts(".ss", ".dic")},
			want:   CustomerBuilder{"static": static.New(static.WithEnable(), static.WithArchive("./xxxx"), static.WithExts(".ss", ".dic"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Static(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Static(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Limit(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []limiter.Option
		repeat []limiter.Option
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认limit对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, args: []limiter.Option{}, want: CustomerBuilder{"acl/limit": limiter.New()}},
		{name: "2. 初始化自定义limit对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []limiter.Option{limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success")))},
			want: CustomerBuilder{"acl/limit": limiter.New(limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success"))))}},
		{name: "3. 重复初始化自定义limit对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args:   []limiter.Option{limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success")))},
			repeat: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("asasas", 1, limiter.WithReponse(500, "fail")))},
			want:   CustomerBuilder{"acl/limit": limiter.New(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("asasas", 1, limiter.WithReponse(500, "fail"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Limit(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Limit(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Proxy(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		script string
		repeat string
		want   CustomerBuilder
	}{
		{name: "1. 初始化空proxy对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "", want: CustomerBuilder{"acl/proxy": ""}},
		{name: "2. 初始化自定义proxy对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "xxxxx", want: CustomerBuilder{"acl/proxy": "xxxxx"}},
		{name: "3. 重复初始化自定义proxy对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "xxxxx", repeat: "yyyyyyyy", want: CustomerBuilder{"acl/proxy": "yyyyyyyy"}},
	}
	for _, tt := range tests {
		got := tt.fields.Proxy(tt.script)
		if tt.repeat != "" {
			got = tt.fields.Proxy(tt.repeat)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Render(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		script string
		repeat string
		want   CustomerBuilder
	}{
		{name: "1. 初始化默认render对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "", want: CustomerBuilder{"render": ""}},
		{name: "2. 初始化自定义render对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "xxxx", want: CustomerBuilder{"render": "xxxx"}},
		{name: "3. 重复初始化自定义render对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})}, script: "xxxx", repeat: "mmmmm", want: CustomerBuilder{"render": "mmmmm"}},
	}
	for _, tt := range tests {
		got := tt.fields.Render(tt.script)
		if tt.repeat != "" {
			got = tt.fields.Render(tt.repeat)
		}
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}
