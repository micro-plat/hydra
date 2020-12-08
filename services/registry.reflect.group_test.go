package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
)

func TestUnitGroup_getPaths(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		mName       string
		wantRpath   string
		wantService string
		wantAction  []string
	}{
		{name: "1.1 函数前缀为空,注册路径没有*", path: "/path", mName: "", wantRpath: "/path", wantService: "/path", wantAction: []string{}},
		{name: "1.2 函数前缀为空,注册路径有一个*", path: "/path/*", mName: "", wantRpath: "/path/handle", wantService: "/path/handle", wantAction: []string{}},
		{name: "1.3 函数前缀为空,注册路径有多个*", path: "/*/path*", mName: "", wantRpath: "/*/pathhandle", wantService: "/*/pathhandle", wantAction: []string{}},

		{name: "2.1 函数前缀为[GET],注册路径没有*", path: "/path", mName: "GET", wantRpath: "/path", wantService: "/path/$GET", wantAction: []string{"GET"}},
		{name: "2.2 函数前缀为[GET],注册路径有一个*", path: "/path/*", mName: "GET", wantRpath: "/path/GET", wantService: "/path/GET/$GET", wantAction: []string{"GET"}},
		{name: "2.3 函数前缀为[GET],注册路径有多个*", path: "/*/path*", mName: "GET", wantRpath: "/*/pathGET", wantService: "/*/pathGET/$GET", wantAction: []string{"GET"}},

		{name: "3.1 函数前缀为[POST],注册路径没有*", path: "/path", mName: "POST", wantRpath: "/path", wantService: "/path/$POST", wantAction: []string{"POST"}},
		{name: "3.2 函数前缀为[POST],注册路径没有*", path: "/path/*", mName: "POST", wantRpath: "/path/POST", wantService: "/path/POST/$POST", wantAction: []string{"POST"}},
		{name: "3.3 函数前缀为[POST],注册路径没有*", path: "/*/path*", mName: "POST", wantRpath: "/*/pathPOST", wantService: "/*/pathPOST/$POST", wantAction: []string{"POST"}},

		{name: "1.4 函数前缀为[ORDER],注册路径没有*", path: "/path", mName: "ORDER", wantRpath: "/path/ORDER", wantService: "/path/ORDER", wantAction: []string{"GET", "POST"}},
		{name: "1.4 函数前缀为[ORDER],注册路径有一个*", path: "/path/*", mName: "FILE", wantRpath: "/path/FILE", wantService: "/path/FILE", wantAction: []string{"GET", "POST"}},
		{name: "1.4 函数前缀为[ORDER],注册路径有多个*", path: "/*/path*", mName: "FILE", wantRpath: "/*/pathFILE", wantService: "/*/pathFILE", wantAction: []string{"GET", "POST"}},
	}
	g := &UnitGroup{}
	for _, tt := range tests {
		gotRpath, gotService, gotAction := g.getPaths(tt.path, tt.mName)
		assert.Equal(t, tt.wantRpath, gotRpath, tt.name)
		assert.Equal(t, tt.wantService, gotService, tt.name)
		assert.Equal(t, tt.wantAction, gotAction, tt.name)
	}
}

func TestUnitGroup_AddHandling(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		h           context.IHandler
		wantService string
	}{
		{name: "1.1 预处理函数前缀为空", h: hander1{}},
		{name: "1.2 预处理函数前缀为GET", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "1.3 预处理函数前缀为GET,替换存在对应的handler", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "1.4 预处理函数前缀为FILE", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandling(tt.mName, tt.h)
		if tt.mName == "" {
			assert.Equal(t, tt.h, g.Handling, tt.name)
			continue
		}
		u := g.Services[tt.wantService]
		assert.Equal(t, g, u.Group, tt.name)
		assert.Equal(t, tt.h, u.Handling, tt.name)
		assert.Equal(t, tt.wantService, u.Service, tt.name)

	}
}

func TestUnitGroup_AddHandled(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		h           context.IHandler
		wantService string
	}{
		{name: "1.1 后处理函数前缀为空", h: hander1{}},
		{name: "1.2 后处理函数前缀为GET", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "1.3 后处理函数前缀为GET,替换存在对应的handler", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "1.4 后处理函数前缀为FILE", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandled(tt.mName, tt.h)
		if tt.mName == "" {
			assert.Equal(t, tt.h, g.Handled, tt.name)
			continue
		}
		u := g.Services[tt.wantService]
		assert.Equal(t, g, u.Group, tt.name)
		assert.Equal(t, tt.h, u.Handled, tt.name)
		assert.Equal(t, tt.wantService, u.Service, tt.name)
	}
}

func TestUnitGroup_AddHandl(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		h           context.IHandler
		wantPath    string
		wantService string
		wantAction  []string
	}{
		{name: "1.1 处理函数前缀为空", h: hander1{}, wantPath: "path", wantService: "path", wantAction: []string{}},
		{name: "1.2 处理函数前缀为GET", mName: "GET", wantPath: "path", wantService: "/path/$GET", h: hander1{}, wantAction: []string{"GET"}},
		{name: "1.3 处理函数前缀为GET,替换存在对应的handler", mName: "GET", wantPath: "path", wantService: "/path/$GET", h: hander2{}, wantAction: []string{"GET"}},
		{name: "1.4 处理函数前缀为FILE", mName: "FILE", wantPath: "/path/FILE", wantService: "/path/FILE", h: hander2{}, wantAction: []string{"GET", "POST"}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandle(tt.mName, tt.h)
		u := g.Services[tt.wantService]

		assert.Equal(t, g, u.Group, tt.name)
		assert.Equal(t, tt.h, u.Handle, tt.name)
		assert.Equal(t, tt.wantPath, u.Path, tt.name)
		assert.Equal(t, tt.wantService, u.Service, tt.name)
		assert.Equal(t, tt.wantAction, u.Actions, tt.name)
	}
}

func TestUnitGroup_AddFallback(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		h           context.IHandler
		wantService string
	}{
		{name: "1.1 降级函数前缀为空", h: hander1{}, wantService: "path"},
		{name: "1.2 降级函数前缀为GET", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "1.3 降级函数前缀为GET,替换存在对应的handler", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "1.4 降级函数前缀为FILE", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddFallback(tt.mName, tt.h)
		u := g.Services[tt.wantService]

		assert.Equal(t, g, u.Group, tt.name)
		assert.Equal(t, tt.h, u.Fallback, tt.name)
		assert.Equal(t, tt.wantService, u.Service, tt.name)
	}
}
