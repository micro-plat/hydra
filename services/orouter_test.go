/*
处理对外提供服务的路由管理，包括注册、获取路由列表
*/

package services

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
)

func Test_pathRouter_Add_WithPanic(t *testing.T) {
	tests := []struct {
		name      string
		service   string
		action    []string
		opts      []interface{}
		wantPanic error
	}{
		{name: "1.1 opts不是router.Option", service: "service", action: []string{"POST"}, opts: []interface{}{"opt1", "opt2"}, wantPanic: fmt.Errorf("%s注册的服务类型必须是router.Option", "service")},
		//{name: "1.2 opts为空", service: "service", action: []string{"POST"}, opts: []interface{}{}, wantPanic: nil},
	}
	p := newPathRouter("path")
	for _, tt := range tests {
		assert.Panics(t, func() {
			p.Add(tt.service, tt.action, tt.opts...)
		}, tt.name)
	}
}

func Test_pathRouter_String(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		routers []*router.Router
	}{
		{name: "1.1 routers为空", want: ""},
		{name: "1.2 routers不为空", want: fmt.Sprintf("%-32s %-32s\n", "service", strings.Join([]string{"POST", "GET"}, " ")), routers: []*router.Router{&router.Router{Path: "path", Service: "service", Action: []string{"POST", "GET"}, Encoding: "utf-8", Pages: []string{"pages"}}}},
	}
	p := newPathRouter("path")
	for _, tt := range tests {
		p.routers = tt.routers
		got := p.String()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

type testAdd struct {
	name    string
	path    string
	service string
	action  []string
	opts    []interface{}
	wantErr bool
}

func getTests(path string, h interface{}, tests *[]*testAdd) {
	//构建测试用例
	g, err := reflectHandle(path, h)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range g.Services {
		*tests = append(*tests, &testAdd{
			name:    fmt.Sprintf("路径%s,service:%s", v.Path, v.Service),
			path:    v.Path,
			service: v.Service,
			action:  v.Actions,
			opts:    []interface{}{router.WithPages("pages"), router.WithEncoding("utf-8")},
			wantErr: false,
		})
	}
}

func TestORouter_Add(t *testing.T) {
	tests := []*testAdd{}

	//添加错误用例
	tests = append(tests, &testAdd{name: "1. action重复", path: "/path", service: "service",
		action: []string{"GET", "GET"}, opts: []interface{}{}, wantErr: true})

	//构建测试用例
	getTests("/path", &testHandler2{}, &tests)
	getTests("/path2", func(ctx context.IContext) (r interface{}) {
		return "success"
	}, &tests)
	getTests("/path3", newTestHandler9, &tests)

	//测试添加
	s := NewORouter("TEST")
	rightTestCount := 0
	for _, tt := range tests {
		err := s.Add(tt.path, tt.service, tt.action, tt.opts...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			s = NewORouter("TEST")
			continue
		}
		rightTestCount++
	}
	s.BuildRouters("")
	//获取routers
	a, err := s.GetRouters()
	assert.Equal(t, false, err != nil, "获取routers")
	assert.Equal(t, rightTestCount, len(a.Routers), "获取routers")
}

func Test_pathRouter_GetRouters(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		testHandler  []interface{}
		wantErr      bool
		errStr       string
		handleAction []string
	}{
		{name: "1.1 同一路径,action被分配完", path: "/path", testHandler: []interface{}{&testHandler3{}}, wantErr: true, errStr: "服务/path无法注册，所有action已分配"},
		{name: "1.2 同一路径,有多个需要分配的action的service", path: "/path", testHandler: []interface{}{&testHandler7{}, func(ctx context.IContext) (r interface{}) { return }}, wantErr: true, errStr: `重复注册的服务{"path":"/path","service":"/path","encoding":"utf-8","pages":["pages"]}`},

		{name: "2.1 获取指针所有routers", path: "/path", testHandler: []interface{}{&testHandler2{}}, handleAction: []string{"GET"}},
		{name: "2.2 获取函数所有routers", path: "/path", testHandler: []interface{}{func(ctx context.IContext) (r interface{}) { return }}, handleAction: []string{"GET", "POST"}},
		{name: "2.3 获取构建函数的所有routers", path: "/path", testHandler: []interface{}{newTestHandler9}, handleAction: []string{"GET", "POST"}},
		{name: "2.4 获取结构体的所有routers", path: "/path", testHandler: []interface{}{testHandler8{}}, handleAction: []string{"GET", "POST"}},
	}
	for _, tt := range tests {
		//添加routers
		rtests := []*testAdd{}
		for _, v := range tt.testHandler {
			getTests(tt.path, v, &rtests)
		}
		s := NewORouter("TEST")
		routers := []*router.Router{}
		for _, tt2 := range rtests {
			s.Add(tt2.path, tt2.service, tt2.action, tt2.opts...)
			temp := &router.Router{
				Path:     tt2.path,
				Action:   tt2.action,
				Service:  tt2.service,
				Encoding: "utf-8",
				Pages:    []string{"pages"},
			}
			if len(tt2.action) == 0 {
				temp.Action = tt.handleAction
			}
			routers = append(routers, temp)
		}

		//获取routers
		a, err := s.GetRouters()
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		//再次获取routers
		a, err = s.GetRouters()
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, len(routers), len(a.Routers), tt.name)
		for _, v1 := range routers {
			for _, v2 := range a.Routers {
				if v1.Service == v2.Service {
					assert.Equal(t, v1.Path, v2.Path, tt.name)

					sort.Strings(v1.Action)
					sort.Strings(v2.Action)
					assert.Equal(t, v1.Action, v2.Action, tt.name+"1")
					assert.Equal(t, v1.Encoding, v2.Encoding, tt.name)

					sort.Strings(v1.Pages)
					sort.Strings(v2.Pages)
					assert.Equal(t, v1.Pages, v2.Pages, tt.name)
				}
			}
		}
	}
}
