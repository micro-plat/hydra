package hydra

import (
	"net/http"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

const (
	HeadeAllowCredentials = "Access-Control-Allow-Credentials"
	HeadeAllowMethods     = "Access-Control-Allow-Methods"
	HeadeAllowOrigin      = "Access-Control-Allow-Origin"
	HeadeAllowHeaders     = "Access-Control-Allow-Headers"
	HeadeExposeHeaders    = "Access-Control-Expose-Headers"
)

var allowHeader = []string{"X-Add-Delay", "X-Request-Id", "X-Requested-With", "Content-Type", "Authorization", "Authorization-Jwt", "Origin", "Accept"}
var exposeHeader = []string{"Authorization-Jwt", "WWW-Authenticate", "Authorization"}
var allMethods = []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

//author:taoshouyin
//time:2020-11-12
//desc:测试header设置中间件逻辑
func TestHeader(t *testing.T) {
	type testCase struct {
		name        string
		isSet       bool
		headerOpts  []header.Option
		wantStatus  int
		wantheader  map[string][]string
		wantSpecial string
	}

	tests := []*testCase{
		{name: "header-不设置延头", isSet: false, wantheader: map[string][]string{}, wantStatus: 200, wantSpecial: ""},
		{name: "header-设置但没有数据", isSet: true, wantheader: map[string][]string{}, wantStatus: 200, wantSpecial: ""},
		{name: "header-设置默认跨域头,无其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain()},
			wantheader: map[string][]string{
				HeadeAllowCredentials: []string{"true"},
				HeadeAllowMethods:     []string{strings.Join(allMethods, ",")},
				HeadeAllowOrigin:      []string{"www.baidu.com"},
				HeadeAllowHeaders:     []string{strings.Join(allowHeader, ",")},
				HeadeExposeHeaders:    []string{strings.Join(exposeHeader, ",")},
			}},
		{name: "header-设置默认跨域头,有其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain(), header.WithHeader("test1", "testttt")},
			wantheader: map[string][]string{
				HeadeAllowCredentials: []string{"true"},
				HeadeAllowMethods:     []string{strings.Join(allMethods, ",")},
				HeadeAllowOrigin:      []string{"www.baidu.com"},
				HeadeAllowHeaders:     []string{strings.Join(allowHeader, ",")},
				HeadeExposeHeaders:    []string{strings.Join(exposeHeader, ",")},
				"test1":               []string{"testttt"},
			}},
		{name: "header-设置指定跨域头,不允许跨域,无其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain("www.sdswds.com")},
			wantheader: map[string][]string{
				HeadeExposeHeaders: []string{strings.Join(exposeHeader, ",")},
			}},
		{name: "header-设置指定跨域头,不允许跨域,有其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain("www.sdswds.com"), header.WithHeader("test1", "testttt")},
			wantheader: map[string][]string{
				HeadeExposeHeaders: []string{strings.Join(exposeHeader, ",")},
				"test1":            []string{"testttt"},
			}},
		{name: "header-设置指定跨域头,允许跨域,无其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain("www.baidu.com")},
			wantheader: map[string][]string{
				HeadeAllowCredentials: []string{"true"},
				HeadeAllowMethods:     []string{strings.Join(allMethods, ",")},
				HeadeAllowOrigin:      []string{"www.baidu.com"},
				HeadeAllowHeaders:     []string{strings.Join(allowHeader, ",")},
				HeadeExposeHeaders:    []string{strings.Join(exposeHeader, ",")},
			}},
		{name: "header-设置指定跨域头,允许跨域,有其他头", isSet: true, wantStatus: 200, wantSpecial: "hdr",
			headerOpts: []header.Option{header.WithCrossDomain("www.baidu.com"), header.WithHeader("test1", "testttt")},
			wantheader: map[string][]string{
				HeadeAllowCredentials: []string{"true"},
				HeadeAllowMethods:     []string{strings.Join(allMethods, ",")},
				HeadeAllowOrigin:      []string{"www.baidu.com"},
				HeadeAllowHeaders:     []string{strings.Join(allowHeader, ",")},
				HeadeExposeHeaders:    []string{strings.Join(exposeHeader, ",")},
				"test1":               []string{"testttt"},
			}},
	}

	for _, tt := range tests {
		mockConf := mocks.NewConf()
		//初始化测试用例参数
		confb := mockConf.GetAPI()
		if tt.isSet {
			confb.Header(tt.headerOpts...)
		}
		serverConf := mockConf.GetAPIConf()
		ctx := &mocks.MiddleContext{
			MockMeta:     conf.NewMeta(),
			MockUser:     &mocks.MockUser{MockClientIP: "192.168.0.1"},
			MockResponse: &mocks.MockResponse{MockStatus: 200, MockHeader: map[string][]string{}},
			MockRequest: &mocks.MockRequest{
				MockPath: &mocks.MockPath{
					MockHeader:      map[string][]string{"Origin": []string{"www.baidu.com"}},
					MockRequestPath: "/header/test",
				},
			},
			MockAPPConf: serverConf,
		}

		//获取中间件
		handler := middleware.Header()
		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, _ := ctx.Response().GetFinalResponse()
		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		gotSpecial := ctx.Response().GetSpecials()
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)
		headers := ctx.Response().GetHeaders()
		assert.Equalf(t, tt.wantheader, headers, tt.name, tt.wantheader, headers)
	}
}
