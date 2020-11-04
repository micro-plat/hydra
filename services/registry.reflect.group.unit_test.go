package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/test/assert"
)

func TestUnit_GetHandlings(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		wantService string
		h           context.IHandler
		want        []context.IHandler
	}{
		{name: "name为支持的http请求类型,group添加handing", mName: "GET", wantService: "/path/$GET", h: hander1{}, want: []context.IHandler{hander1{}}},
		{name: "name为空,group添加handing", mName: "", wantService: "/path/$GET", h: hander2{}, want: []context.IHandler{hander2{}, hander1{}}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandling(tt.mName, tt.h)
		got := g.Services[tt.wantService].GetHandlings()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestUnit_GetHandleds(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		wantService string
		h           context.IHandler
		want        []context.IHandler
	}{
		{name: "name为支持的http请求类型,group添加handing", mName: "GET", wantService: "/path/$GET", h: hander1{}, want: []context.IHandler{hander1{}}},
		{name: "name为空,group添加handing", mName: "", wantService: "/path/$GET", h: hander2{}, want: []context.IHandler{hander1{}, hander2{}}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandled(tt.mName, tt.h)
		got := g.Services[tt.wantService].GetHandleds()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
