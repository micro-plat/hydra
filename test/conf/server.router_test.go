package conf

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/test/assert"

	"github.com/micro-plat/hydra/conf/server/router"
)

func TestGetWSHomeRouter(t *testing.T) {
	tests := []struct {
		name string
		want *router.Router
	}{
		{name: "获取默认对象WS", want: &router.Router{Path: "/", Action: router.Methods, Service: "/"}},
	}
	for _, tt := range tests {
		got := router.GetWSHomeRouter()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewRouters(t *testing.T) {
	tests := []struct {
		name string
		want *router.Routers
	}{
		{name: "获取默认对象", want: &router.Routers{Routers: make([]*router.Router, 0, 1)}},
	}
	for _, tt := range tests {
		got := router.NewRouters()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewRouter(t *testing.T) {
	type args struct {
		path    string
		service string
		action  []string
		opts    []router.Option
	}
	tests := []struct {
		name string
		args args
		want *router.Router
	}{
		{name: "初始化空对象", args: args{path: "", service: "", action: nil, opts: nil}, want: &router.Router{Path: "", Action: nil, Service: "", Encoding: "", Pages: nil}},
		{name: "初始化可用对象", args: args{path: "/t1", service: "server1", action: []string{"xxx"}, opts: nil},
			want: &router.Router{Path: "/t1", Action: []string{"xxx"}, Service: "server1", Encoding: "", Pages: nil}},
		{name: "初始化全量对象", args: args{path: "/t1", service: "server1", action: []string{"xxx", "yyy"}, opts: []router.Option{router.WithEncoding("gbk2312"), router.WithPages("asas", "sasa")}},
			want: &router.Router{Path: "/t1", Action: []string{"xxx", "yyy"}, Service: "server1", Encoding: "gbk2312", Pages: []string{"asas", "sasa"}}},
	}
	for _, tt := range tests {
		got := router.NewRouter(tt.args.path, tt.args.service, tt.args.action, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouters_String(t *testing.T) {
	f := func(h []*router.Router) string {
		var sb strings.Builder
		for _, v := range h {
			sb.WriteString(fmt.Sprintf("%-16s %-32s %-32s %v\n", v.Path, v.Service, strings.Join(v.Action, " "), v.Pages))
		}
		return sb.String()
	}

	tests := []struct {
		name    string
		routers []*router.Router
		want    string
	}{
		{name: "空对象获取str", routers: []*router.Router{}, want: f([]*router.Router{})},
		{name: "单个对象获取str", routers: []*router.Router{router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk"))},
			want: f([]*router.Router{router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk"))})},
		{name: "多个对象获取str", routers: []*router.Router{router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk")), router.NewRouter("/t2", "s2", []string{"post"}, router.WithEncoding("ggg"))},
			want: f([]*router.Router{router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk")), router.NewRouter("/t2", "s2", []string{"post"}, router.WithEncoding("ggg"))})},
	}
	for _, tt := range tests {
		h := &router.Routers{
			Routers: tt.routers,
		}
		got := h.String()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouter_GetEncoding(t *testing.T) {
	tests := []struct {
		name   string
		fields *router.Router
		want   string
	}{
		{name: "获取默认编码方式", fields: router.NewRouter("/t1", "s1", []string{"get"}), want: "utf-8"},
		{name: "设置编码方式", fields: router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk2312")), want: "gbk2312"},
	}
	for _, tt := range tests {
		got := tt.fields.GetEncoding()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouter_IsUTF8(t *testing.T) {
	tests := []struct {
		name   string
		fields *router.Router
		want   bool
	}{
		{name: "默认编码是utf8", fields: router.NewRouter("/t1", "s1", []string{"get"}), want: true},
		{name: "不是utf8", fields: router.NewRouter("/t1", "s1", []string{"get"}, router.WithEncoding("gbk2312")), want: false},
	}
	for _, tt := range tests {
		got := tt.fields.IsUTF8()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouters_Append(t *testing.T) {
	type args struct {
		path    string
		service string
		action  []string
		opts    []router.Option
	}
	tests := []struct {
		name string
		obj  *router.Routers
		args args
		want *router.Routers
	}{
		{name: "空对象累加", obj: &router.Routers{Routers: []*router.Router{}}, args: args{},
			want: &router.Routers{Routers: []*router.Router{router.NewRouter("", "", nil)}}},
		{name: "空对象累加实体", obj: &router.Routers{Routers: []*router.Router{}}, args: args{path: "/p1", service: "s1", action: []string{"22"}, opts: []router.Option{router.WithEncoding("gbk")}},
			want: &router.Routers{Routers: []*router.Router{router.NewRouter("/p1", "s1", []string{"22"}, []router.Option{router.WithEncoding("gbk")}...)}}},
		{name: "对象累加实体", obj: &router.Routers{Routers: []*router.Router{router.NewRouter("/p1", "s1", []string{"22"}, []router.Option{router.WithEncoding("gbk")}...)}}, args: args{path: "/p1", service: "s1", action: []string{"22"}, opts: []router.Option{router.WithEncoding("gbk")}},
			want: &router.Routers{Routers: []*router.Router{router.NewRouter("/p1", "s1", []string{"22"}, []router.Option{router.WithEncoding("gbk")}...), router.NewRouter("/p1", "s1", []string{"22"}, []router.Option{router.WithEncoding("gbk")}...)}}},
	}
	for _, tt := range tests {
		got := tt.obj.Append(tt.args.path, tt.args.service, tt.args.action, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouters_Match(t *testing.T) {
	type args struct {
		path   string
		method string
	}
	type test struct {
		name   string
		fields *router.Routers
		args   args
		want   *router.Router
	}
	tests := []test{
		{name: "匹配路径为空", fields: router.NewRouters(), args: args{path: "", method: "GET"}, want: &router.Router{Path: "", Action: []string{"GET"}}},
		{name: "匹配方法Options", fields: router.NewRouters(), args: args{path: "/1", method: http.MethodOptions}, want: &router.Router{Path: "/1", Action: []string{http.MethodOptions}}},
		{name: "匹配方法Head", fields: router.NewRouters(), args: args{path: "/1", method: http.MethodHead}, want: &router.Router{Path: "/1", Action: []string{http.MethodHead}}},
		{name: "存在匹配的路由", fields: &router.Routers{Routers: []*router.Router{router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312"))}},
			args: args{path: "/t1/t2", method: http.MethodGet}, want: router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312"))},
	}
	for _, tt := range tests {
		got := tt.fields.Match(tt.args.path, tt.args.method)
		assert.Equal(t, tt.want, got, tt.name)
	}

	defer func() {
		e := recover()
		assert.NotEqual(t, nil, e, "不存在匹配的路由匹配结果不正常")
	}()

	test1 := test{name: "不存在匹配的路由", fields: &router.Routers{Routers: []*router.Router{router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312"))}},
		args: args{path: "/t1/t2/tt", method: http.MethodGet}, want: router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312"))}
	got := test1.fields.Match(test1.args.path, test1.args.method)
	assert.Equal(t, test1.want, got, test1.name)
}

func TestRouters_GetPath(t *testing.T) {
	tests := []struct {
		name   string
		fields *router.Routers
		want   []string
	}{
		{name: "空对象获取", fields: router.NewRouters(), want: make([]string, 0)},
		{name: "单个路径对象获取", fields: &router.Routers{Routers: []*router.Router{router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312"))}}, want: []string{"/t1/t2"}},
		{name: "多个路径对象获取", fields: &router.Routers{Routers: []*router.Router{router.NewRouter("/t1/t2", "s1", []string{http.MethodGet}, router.WithEncoding("gbk2312")), router.NewRouter("/t1/t2/t3", "s2", []string{http.MethodGet}, router.WithEncoding("gbk"))}}, want: []string{"/t1/t2", "/t1/t2/t3"}},
	}
	for _, tt := range tests {
		got := tt.fields.GetPath()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRouterGetConf(t *testing.T) {
	type test struct {
		name    string
		cnf     conf.IMainConf
		want    *router.Routers
		wantErr bool
	}
	//暂时不能测试   由于不能初始化注册路由
	// gotRouter, err := router.GetConf(tt.args.cnf)

}
