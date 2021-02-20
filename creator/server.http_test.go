package creator

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/static"
)

func Test_newHTTP(t *testing.T) {
	type args struct {
		tp      string
		address string
		opts    []api.Option
	}
	tests := []struct {
		name   string
		args   args
		repeat *args
		want   *httpBuilder
	}{
		{name: "1. 无option配置,初始化", args: args{tp: "tp1", address: ":1122", opts: []api.Option{}},
			want: &httpBuilder{tp: "tp1", BaseBuilder: map[string]interface{}{"main": api.New(":1122")}}},
		{name: "2. 设置option配置,初始化", args: args{tp: "tp1", address: ":1122", opts: []api.Option{api.WithDisable()}},
			want: &httpBuilder{tp: "tp1", BaseBuilder: map[string]interface{}{"main": api.New(":1122", api.WithDisable())}}},
		{name: "3. 重复设置节点", args: args{tp: "tp1", address: ":1122", opts: []api.Option{api.WithDisable()}},
			repeat: &args{tp: "tp2", address: ":1123", opts: []api.Option{api.WithDisable()}},
			want:   &httpBuilder{tp: "tp2", BaseBuilder: map[string]interface{}{"main": api.New(":1123", api.WithDisable())}}},
	}
	for _, tt := range tests {
		got := newHTTP(tt.args.tp, tt.args.address, tt.args.opts...)
		if tt.repeat != nil {
			got = newHTTP(tt.repeat.tp, tt.repeat.address, tt.repeat.opts...)
		}
		assert.Equal(t, tt.want.tp, got.tp, tt.name+",tp")
		assert.Equal(t, tt.want.BaseBuilder, got.BaseBuilder, tt.name+",BaseBuilder")
	}
}

func Test_httpBuilder_Jwt(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []jwt.Option
		repeat []jwt.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认jwt配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456")}, want: BaseBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("123456"))}},
		{name: "2. 初始化自定义jwt配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456"), jwt.WithExcludes("/taews/ssss")},
			want: BaseBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("123456"), jwt.WithExcludes("/taews/ssss"))}},
		{name: "3. 重复初始化jwt配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []jwt.Option{jwt.WithSecret("123456")},
			repeat: []jwt.Option{jwt.WithSecret("3232323"), jwt.WithExcludes("/taews/xxx")}, want: BaseBuilder{"auth/jwt": jwt.NewJWT(jwt.WithSecret("3232323"), jwt.WithExcludes("/taews/xxx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Jwt(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Jwt(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_APIKEY(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		fields *httpBuilder
		args   []apikey.Option
		repeat []apikey.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认apikey对象", secret: "", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []apikey.Option{}, want: BaseBuilder{"auth/apikey": apikey.New("")}},
		{name: "2. 初始化自定义apikey对象", secret: "123456", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []apikey.Option{apikey.WithSHA256Mode(), apikey.WithDisable()}, want: BaseBuilder{"auth/apikey": apikey.New("123456", apikey.WithSHA256Mode(), apikey.WithDisable())}},
		{name: "3. 重复初始化apikey对象", secret: "123456", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []apikey.Option{apikey.WithSHA256Mode(), apikey.WithDisable()},
			repeat: []apikey.Option{apikey.WithSHA1Mode(), apikey.WithSecret("xxxxxx")},
			want:   BaseBuilder{"auth/apikey": apikey.New("123456", apikey.WithSHA1Mode(), apikey.WithSecret("xxxxxx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.APIKEY(tt.secret, tt.args...)
		if tt.repeat != nil {
			got = tt.fields.APIKEY(tt.secret, tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Basic(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []basic.Option
		repeat []basic.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认basic配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []basic.Option{}, want: BaseBuilder{"auth/basic": basic.NewBasic()}},
		{name: "2. 初始化自定义basic配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []basic.Option{basic.WithDisable(), basic.WithExcludes("11s")}, want: BaseBuilder{"auth/basic": basic.NewBasic(basic.WithDisable(), basic.WithExcludes("11s"))}},
		{name: "3. 重复初始化自定义basic配置对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []basic.Option{basic.WithDisable(), basic.WithExcludes("11s")},
			repeat: []basic.Option{basic.WithUP("ssss", "bbbbb"), basic.WithExcludes("11xx")},
			want:   BaseBuilder{"auth/basic": basic.NewBasic(basic.WithUP("ssss", "bbbbb"), basic.WithExcludes("11xx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Basic(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Basic(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_WhiteList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []whitelist.Option
		repeat []whitelist.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认whitelist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []whitelist.Option{}, want: BaseBuilder{"acl/white.list": whitelist.New()}},
		{name: "2. 初始化自定义whitelist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []whitelist.Option{whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"}))},
			want: BaseBuilder{"acl/white.list": whitelist.New(whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"})))}},
		{name: "3. 重复初始化自定义whitelist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []whitelist.Option{whitelist.WithDisable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request"}))},
			repeat: []whitelist.Option{whitelist.WithEnable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request1"}))},
			want:   BaseBuilder{"acl/white.list": whitelist.New(whitelist.WithEnable(), whitelist.WithIPList(whitelist.NewIPList([]string{"request1"})))}},
	}
	for _, tt := range tests {
		got := tt.fields.WhiteList(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.WhiteList(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_BlackList(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []blacklist.Option
		repeat []blacklist.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认blacklist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []blacklist.Option{}, want: BaseBuilder{"acl/black.list": blacklist.New()}},
		{name: "2. 初始化自定义blacklist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.101")},
			want: BaseBuilder{"acl/black.list": blacklist.New(blacklist.WithDisable(), blacklist.WithIP("192.168.0.101"))}},
		{name: "3. 重复初始化自定义blacklist对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []blacklist.Option{blacklist.WithDisable(), blacklist.WithIP("192.168.0.101")},
			repeat: []blacklist.Option{blacklist.WithEnable(), blacklist.WithIP("192.168.0.111")},
			want:   BaseBuilder{"acl/black.list": blacklist.New(blacklist.WithEnable(), blacklist.WithIP("192.168.0.111"))}},
	}
	for _, tt := range tests {
		got := tt.fields.BlackList(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.BlackList(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Ras(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []ras.Option
		repeat []ras.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认ras对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []ras.Option{}, want: BaseBuilder{"auth/ras": ras.NewRASAuth()}},
		{name: "2. 初始化自定义ras对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1")))},
			want: BaseBuilder{"auth/ras": ras.NewRASAuth(ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1"))))}},
		{name: "3. 重复初始化自定义ras对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []ras.Option{ras.WithDisable(), ras.WithAuths(ras.New("server1", ras.WithAuthDisable(), ras.WithRequest("patch1")))},
			repeat: []ras.Option{ras.WithEnable(), ras.WithAuths(ras.New("server2", ras.WithAuthDisable(), ras.WithRequest("patch2")))},
			want:   BaseBuilder{"auth/ras": ras.NewRASAuth(ras.WithEnable(), ras.WithAuths(ras.New("server2", ras.WithAuthDisable(), ras.WithRequest("patch2"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Ras(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Ras(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Header(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []header.Option
		repeat []header.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认header对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []header.Option{}, want: BaseBuilder{"header": header.New()}},
		{name: "2. 初始化自定义header对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []header.Option{header.WithAllowMethods("get", "put")}, want: BaseBuilder{"header": header.New(header.WithAllowMethods("get", "put"))}},
		{name: "3. 重复初始化自定义header对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []header.Option{header.WithAllowMethods("get", "put")},
			repeat: []header.Option{header.WithCrossDomain("www.baidu.com"), header.WithAllowMethods("get", "put")},
			want:   BaseBuilder{"header": header.New(header.WithCrossDomain("www.baidu.com"), header.WithAllowMethods("get", "put"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Header(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Header(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Static(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []static.Option
		repeat []static.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认static对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []static.Option{}, want: BaseBuilder{"static": static.New()}},
		{name: "2. 初始化自定义static对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []static.Option{static.WithDisable(), static.WithAssetsPath("./sssss")},
			want: BaseBuilder{"static": static.New(static.WithDisable(), static.WithAssetsPath("./sssss"))}},
		{name: "3. 重复初始化自定义static对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []static.Option{static.WithDisable(), static.WithAssetsPath("./sssss")},
			repeat: []static.Option{static.WithEnable(), static.WithAssetsPath("./xxxx")},
			want:   BaseBuilder{"static": static.New(static.WithEnable(), static.WithAssetsPath("./xxxx"))}},
	}
	for _, tt := range tests {
		got := tt.fields.Static(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Static(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Limit(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		args   []limiter.Option
		repeat []limiter.Option
		want   BaseBuilder
	}{
		{name: "1. 初始化默认limit对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, args: []limiter.Option{}, want: BaseBuilder{"acl/limit": limiter.New()}},
		{name: "2. 初始化自定义limit对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args: []limiter.Option{limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success")))},
			want: BaseBuilder{"acl/limit": limiter.New(limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success"))))}},
		{name: "3. 重复初始化自定义limit对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})},
			args:   []limiter.Option{limiter.WithDisable(), limiter.WithRuleList(limiter.NewRule("patch1", 1, limiter.WithReponse(100, "success")))},
			repeat: []limiter.Option{limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("asasas", 1, limiter.WithReponse(500, "fail")))},
			want:   BaseBuilder{"acl/limit": limiter.New(limiter.WithEnable(), limiter.WithRuleList(limiter.NewRule("asasas", 1, limiter.WithReponse(500, "fail"))))}},
	}
	for _, tt := range tests {
		got := tt.fields.Limit(tt.args...)
		if tt.repeat != nil {
			got = tt.fields.Limit(tt.repeat...)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Proxy(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		script string
		repeat string
		want   BaseBuilder
	}{
		{name: "1. 初始化空proxy对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "", want: BaseBuilder{"acl/proxy": ""}},
		{name: "2. 初始化自定义proxy对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "xxxxx", want: BaseBuilder{"acl/proxy": "xxxxx"}},
		{name: "3. 重复初始化自定义proxy对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "xxxxx", repeat: "yyyyyyyy", want: BaseBuilder{"acl/proxy": "yyyyyyyy"}},
	}
	for _, tt := range tests {
		got := tt.fields.Proxy(tt.script)
		if tt.repeat != "" {
			got = tt.fields.Proxy(tt.repeat)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}

func Test_httpBuilder_Render(t *testing.T) {
	tests := []struct {
		name   string
		fields *httpBuilder
		script string
		repeat string
		want   BaseBuilder
	}{
		{name: "1. 初始化默认render对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "", want: BaseBuilder{"render": ""}},
		{name: "2. 初始化自定义render对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "xxxx", want: BaseBuilder{"render": "xxxx"}},
		{name: "3. 重复初始化自定义render对象", fields: &httpBuilder{tp: "x1", BaseBuilder: make(map[string]interface{})}, script: "xxxx", repeat: "mmmmm", want: BaseBuilder{"render": "mmmmm"}},
	}
	for _, tt := range tests {
		got := tt.fields.Render(tt.script)
		if tt.repeat != "" {
			got = tt.fields.Render(tt.repeat)
		}
		assert.Equal(t, tt.want, got.BaseBuilder, tt.name)
	}
}
