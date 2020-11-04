package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/test/assert"
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
		{name: "name为空,handleName为空", path: "path", mName: "", wantRpath: "path", wantService: "path", wantAction: []string{}},
		{name: "name为支持的http请求类型", path: "path", mName: "GET", wantRpath: "path", wantService: "/path/$GET", wantAction: []string{"GET"}},
		{name: "name为支持的http请求类型", path: "path", mName: "GET", wantRpath: "path", wantService: "/path/$GET", wantAction: []string{"GET"}},
		{name: "name为不支持的http请求类型,且HandleName为空", path: "path", mName: "FILE", wantRpath: "/path/FILE", wantService: "/path/FILE", wantAction: []string{"GET", "POST"}},
		{name: "name为不支持的http请求类型,且HandleName不为空", path: "*/**path/*/*", mName: "FILE", wantRpath: "*/**path/*/FILE", wantService: "*/**path/*/FILE", wantAction: []string{"GET", "POST"}},
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
		{name: "name为空", h: hander1{}},
		{name: "name为支持的http请求类型,添加不存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "name为支持的http请求类型,添加存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "name为不支持的http请求类型", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
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
		{name: "name为空", h: hander1{}},
		{name: "name为支持的http请求类型,添加不存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "name为支持的http请求类型,添加存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "name为不支持的http请求类型", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
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
		{name: "name为空", h: hander1{}, wantPath: "path", wantService: "path", wantAction: []string{}},
		{name: "name为支持的http请求类型,添加不存在serviecs", mName: "GET", wantPath: "path", wantService: "/path/$GET", h: hander1{}, wantAction: []string{"GET"}},
		{name: "name为支持的http请求类型,添加存在serviecs", mName: "GET", wantPath: "path", wantService: "/path/$GET", h: hander2{}, wantAction: []string{"GET"}},
		{name: "name为不支持的http请求类型", mName: "FILE", wantPath: "/path/FILE", wantService: "/path/FILE", h: hander2{}, wantAction: []string{"GET", "POST"}},
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
		{name: "name为空", h: hander1{}, wantService: "path"},
		{name: "name为支持的http请求类型,添加不存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander1{}},
		{name: "name为支持的http请求类型,添加存在serviecs", mName: "GET", wantService: "/path/$GET", h: hander2{}},
		{name: "name为不支持的http请求类型", mName: "FILE", wantService: "/path/FILE", h: hander2{}},
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
