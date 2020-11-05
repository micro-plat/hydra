package creator

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
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
		name string
		args args
		want *httpBuilder
	}{
		{name: "没有option初始化", args: args{tp: "tp1", address: ":1122", f: func(string) *services.ORouter { return nil }, opts: []api.Option{}},
			want: &httpBuilder{tp: "tp1", fnGetRouter: func(string) *services.ORouter { return nil }, CustomerBuilder: map[string]interface{}{"main": api.New(":1122")}}},
		{name: "option初始化", args: args{tp: "tp1", address: ":1122", f: func(string) *services.ORouter { return nil }, opts: []api.Option{api.WithDisable()}},
			want: &httpBuilder{tp: "tp1", fnGetRouter: func(string) *services.ORouter { return nil }, CustomerBuilder: map[string]interface{}{"main": api.New(":1122", api.WithDisable())}}},
	}
	for _, tt := range tests {
		got := newHTTP(tt.args.tp, tt.args.address, tt.args.f, tt.args.opts...)
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
		{name: "空路由", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
			return services.NewORouter()
		}, CustomerBuilder: make(map[string]interface{})}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "有重复路由", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
			r := services.NewORouter()
			r.Add("path1", "service1", []string{"get"})
			r.Add("path1", "service1", []string{"get"})
			return r
		}, CustomerBuilder: make(map[string]interface{})}, want: CustomerBuilder{"router": router.NewRouters()}},
		{name: "正常路由", c: httpBuilder{tp: "api", fnGetRouter: func(string) *services.ORouter {
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
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []jwt.Option{jwt.WithSecret("123456")},
			want: CustomerBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("123456"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Jwt(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_APIKEY(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		fields *httpBuilder
		args   []apikey.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", secret: "", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []apikey.Option{},
			want: CustomerBuilder{"auth/apikey": apikey.New("")}},
		{name: " 初始化实体对象", secret: "123456", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []apikey.Option{apikey.WithSHA256Mode(), apikey.WithDisable()},
			want: CustomerBuilder{"auth/apikey": apikey.New("123456", apikey.WithSHA256Mode(), apikey.WithDisable())}},
	}
	for _, tt := range tests {
		got := tt.fields.APIKEY(tt.secret, tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Basic(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []basic.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []basic.Option{},
			want: CustomerBuilder{"auth/basic": basic.NewBasic()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []basic.Option{basic.WithDisable(), basic.WithExcludes("11s")},
			want: CustomerBuilder{"auth/basic": basic.NewBasic(basic.WithDisable(), basic.WithExcludes("11s"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Basic(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_WhiteList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []whitelist.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []whitelist.Option{},
			want: CustomerBuilder{"acl/white.list": whitelist.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []whitelist.Option{whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList("request"))},
			want: CustomerBuilder{"acl/white.list": whitelist.New(whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList("request")))}},
	}
	for _, tt := range tests {
		got := tt.fields.WhiteList(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_BlackList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []blacklist.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []blacklist.Option{},
			want: CustomerBuilder{"acl/black.list": blacklist.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.101")},
			want: CustomerBuilder{"acl/black.list": blacklist.New(blacklist.WithDisable(), blacklist.WithIP("192.168.0.101"))}},
	}
	for _, tt := range tests {
		got := tt.fields.BlackList(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Ras(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []ras.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []ras.Option{},
			want: CustomerBuilder{"auth/ras": ras.NewRASAuth()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1")))},
			want: CustomerBuilder{"auth/ras": ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Ras(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Header(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []header.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []header.Option{},
			want: CustomerBuilder{"header": header.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []header.Option{header.WithAllowMethods("get", "put")},
			want: CustomerBuilder{"header": header.New(header.WithAllowMethods("get", "put"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Header(tt.args...)
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
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: args{host: "host1", db: "db1", cron: "cron1", opts: []metric.Option{}},
			want: CustomerBuilder{"metric": metric.New("host1", "db1", "cron1")}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: args{host: "host1", db: "db1", cron: "cron1", opts: []metric.Option{metric.WithDisable(), metric.WithUPName("name", "pwd")}},
			want: CustomerBuilder{"metric": metric.New("host1", "db1", "cron1", metric.WithDisable(), metric.WithUPName("name", "pwd"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Metric(tt.args.host, tt.args.db, tt.args.cron, tt.args.opts...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Static(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []static.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []static.Option{},
			want: CustomerBuilder{"static": static.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []static.Option{static.WithDisable(), static.WithArchive("./sssss")},
			want: CustomerBuilder{"static": static.New(static.WithDisable(), static.WithArchive("./sssss"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Static(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Limit(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []limiter.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []limiter.Option{},
			want: CustomerBuilder{"acl/limit": limiter.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []limiter.Option{limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success")))},
			want: CustomerBuilder{"acl/limit": limiter.New(limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Limit(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Gray(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []gray.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []gray.Option{},
			want: CustomerBuilder{"acl/gray": gray.New()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []gray.Option{gray.WithDisable(), gray.WithFilter("sss")},
			want: CustomerBuilder{"acl/gray": gray.New(gray.WithDisable(), gray.WithFilter("sss"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Gray(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}

func Test_httpBuilder_Render(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []render.Option
		want   CustomerBuilder
	}{
		{name: " 初始化默认对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []render.Option{},
			want: CustomerBuilder{"render": render.NewRender()}},
		{name: " 初始化实体对象", fields: &httpBuilder{tp: "x1", fnGetRouter: nil, CustomerBuilder: make(map[string]interface{})},
			args: []render.Option{render.WithDisable(), render.WithTmplt("patch1", "content1", render.WithStatus("templt"))},
			want: CustomerBuilder{"render": render.NewRender(render.WithDisable(), render.WithTmplt("patch1", "content1", render.WithStatus("templt")))}},
	}
	for _, tt := range tests {
		got := tt.fields.Render(tt.args...)
		assert.Equal(t, tt.want, got.CustomerBuilder, tt.name)
	}
}
