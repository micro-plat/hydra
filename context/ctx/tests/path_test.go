package tests

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
)

func Test_rpath_GetRouter(t *testing.T) {

	startServer()
	httpConf, _ := server.Cache.GetServerConf("api")
	cronConf, _ := server.Cache.GetServerConf("cron")

	type fields struct {
		ctx        context.IInnerContext
		serverConf server.IServerConf
		meta       conf.IMeta
		isLimit    bool
		fallback   bool
	}
	tests := []struct {
		name      string
		fields    fields
		want      *router.Router
		wantError string
	}{
		{name: "http正确路径和正确方法", fields: fields{ctx: &TestContxt{
			routerpath: "/api",
			method:     "GET",
		}, serverConf: httpConf, meta: conf.NewMeta()}, want: &router.Router{
			Path:    "/api",
			Action:  []string{"GET", "POST"},
			Service: "/api",
		}},
		{name: "http非正确路径和方法", fields: fields{ctx: &TestContxt{
			routerpath: "/api",
			method:     "DELETE",
		}, serverConf: httpConf, meta: conf.NewMeta()}, wantError: "未找到与[/api][DELETE]匹配的路由"},
		{name: "非http的路径和的方法", fields: fields{ctx: &TestContxt{
			routerpath: "/cron",
		}, serverConf: cronConf, meta: conf.NewMeta()}, want: &router.Router{
			Path:     "/cron",
			Encoding: "utf-8",
			Action:   []string{},
			Service:  "/cron",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewRpath(tt.fields.ctx, tt.fields.serverConf, tt.fields.meta)
			defer func() {
				if r := recover(); r != nil {
					if tt.wantError == r {
						return
					}
					t.Errorf("rpath.GetRouter() doesn't run")
				}
			}()
			if got := c.GetRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rpath.GetRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rpath_GetCookies(t *testing.T) {
	startServer()
	serverConf, _ := server.Cache.GetServerConf("api")
	type fields struct {
		ctx context.IInnerContext
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{name: "获取全部cookies", fields: fields{ctx: &TestContxt{
			cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
		}}, want: map[string]string{"cookie1": "value1", "cookie2": "value2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewRpath(tt.fields.ctx, serverConf, conf.NewMeta())
			if got := c.GetCookies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rpath.GetCookies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rpath_GetCookie(t *testing.T) {
	startServer()
	serverConf, _ := server.Cache.GetServerConf("api")
	rpath := ctx.NewRpath(&TestContxt{
		cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
	}, serverConf, conf.NewMeta())

	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{name: "获取存在cookies", args: args{name: "cookie2"}, want: "value2", want1: true},
		{name: "获取不存在cookies", args: args{name: "cookie3"}, want: "", want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := rpath.GetCookie(tt.args.name)
			if got != tt.want {
				t.Errorf("rpath.GetCookie() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("rpath.GetCookie() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
