package context

import (
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	c "github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_rpath_GetRouter_WithPanic(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")       //初始化参数
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
		{name: "http非正确路径和方法", ctx: &mocks.TestContxt{
			Routerpath: "/api",
			Method:     "DELETE",
		}, serverConf: httpConf, meta: conf.NewMeta(), wantError: "未找到与[/api][DELETE]匹配的路由"},
	}

	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, tt.serverConf, tt.meta)
		assert.PanicError(t, tt.wantError, func() {
			c.GetRouter()
		}, tt.name)
	}
}

func Test_rpath_GetRouter(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")       //初始化参数
	confObj.CRON(c.WithMasterSlave(), c.WithTrace())
	confObj.Service.API.Add("/api", "/api", []string{"GET"})
	httpConf := confObj.GetAPIConf()  //获取配置
	cronConf := confObj.GetCronConf() //获取配置

	tests := []struct {
		name       string
		ctx        context.IInnerContext
		serverConf app.IAPPConf
		meta       conf.IMeta
		want       *router.Router
		wantError  string
	}{
		{name: "http正确路径和正确方法", ctx: &mocks.TestContxt{
			Routerpath: "/api",
			Method:     "GET",
		}, serverConf: httpConf, meta: conf.NewMeta(), want: &router.Router{
			Path:    "/api",
			Action:  []string{"GET"},
			Service: "/api",
		}},
		{name: "非http的路径和的方法", ctx: &mocks.TestContxt{
			Routerpath: "/cron",
		}, serverConf: cronConf, meta: conf.NewMeta(), want: &router.Router{
			Path:     "/cron",
			Encoding: "utf-8",
			Action:   []string{},
			Service:  "/cron",
		}},
	}

	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, tt.serverConf, tt.meta)
		got, err := c.GetRouter()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_rpath_GetCookies(t *testing.T) {

	tests := []struct {
		name string
		ctx  context.IInnerContext
		want map[string]string
	}{
		{name: "获取全部cookies", ctx: &mocks.TestContxt{
			Cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
		}, want: map[string]string{"cookie1": "value1", "cookie2": "value2"}},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	for _, tt := range tests {
		c := ctx.NewRpath(tt.ctx, serverConf, conf.NewMeta())
		got := c.GetCookies()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_rpath_GetCookie(t *testing.T) {

	tests := []struct {
		name       string
		cookieName string
		want       string
		want1      bool
	}{
		{name: "获取存在cookies", cookieName: "cookie2", want: "value2", want1: true},
		{name: "获取不存在cookies", cookieName: "cookie3", want: "", want1: false},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	rpath := ctx.NewRpath(&mocks.TestContxt{
		Cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
	}, serverConf, conf.NewMeta())

	for _, tt := range tests {
		got, got1 := rpath.GetCookie(tt.cookieName)
		assert.Equal(t, tt.want, got, tt.name)
		assert.Equal(t, tt.want1, got1, tt.name)
	}
}
