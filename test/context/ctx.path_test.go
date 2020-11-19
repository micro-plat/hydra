package context

import (
	"io/ioutil"
	"net/http"
	"strings"
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

func Test_rpath_GetCookies_WithHttp(t *testing.T) {

	tests := []struct {
		name            string
		contentType     string
		cookie          http.Cookie
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		//net/http: invalid byte 'Ö' in Cookie.Value; dropping invalid bytes
		{name: "cookie编码为中文GBK,无法提交的cookie", contentType: "application/json;charset=gbk", wantContentType: "application/json; charset=gbk",
			cookie:     http.Cookie{Name: "cname", Value: Utf8ToGbk("中文")},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":""}`},
		//net/http: invalid byte 'ä' in Cookie.Value; dropping invalid bytes
		{name: "cookie编码为中文UTF-8,无法提交的cookie", contentType: "application/json;charset=utf-8", wantContentType: "application/json; charset=utf-8",
			cookie:     http.Cookie{Name: "cname", Value: "中文"},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":""}`},
		{name: "cookie不带中文", contentType: "application/json;charset=utf-8", wantContentType: "application/json; charset=utf-8",
			cookie:     http.Cookie{Name: "cname", Value: "value!@#$%^&*()_+"},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":"value!@#$%^\u0026*()_+"}`},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getcookies/encoding", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.AddCookie(&tt.cookie)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, "200 OK", resp.Status, tt.name)
		assert.Equal(t, 200, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_rpath_GetHeader_Encoding(t *testing.T) {

	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部编码为GBK,内容为GBK", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部编码为GBK,内容为UTF8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码为UTF-8,内容为UTF8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码为UTF-8,内容为GBK", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_rpath_GetHeader_Encoding_GBK(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部为gbk,内容为gbk", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部为gbk,内容为utf-8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为utf-8,内容为gbk", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部为utf-8,内容为utf-8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding/gbk", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}
func Test_rpath_GetHeader_Encoding_UTF8(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部为gbk,内容为gbk", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为gbk,内容为utf-8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部为utf-8,内容为gbk", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为utf-8,内容为utf-8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding/utf8", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
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
