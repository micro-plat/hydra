package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
)

func TestUnit_GetHandlings(t *testing.T) {
	tests := []struct {
		name        string
		mName       string
		wantService string
		h           context.IHandler
		want        []context.IHandler
	}{
		{name: "1.1 添加GET的预处理函数", mName: "GET", wantService: "/path/$GET", h: hander1{}, want: []context.IHandler{hander1{}}},
		{name: "1.2 添加没有前缀的预处理函数", mName: "", wantService: "/path/$GET", h: hander2{}, want: []context.IHandler{hander2{}, hander1{}}},
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
		{name: "1.1 添加GET的后处理函数", mName: "GET", wantService: "/path/$GET", h: hander1{}, want: []context.IHandler{hander1{}}},
		{name: "1.2 添加没有前缀的后处理函数", mName: "", wantService: "/path/$GET", h: hander2{}, want: []context.IHandler{hander1{}, hander2{}}},
	}
	g := newUnitGroup("path")
	for _, tt := range tests {
		g.AddHandled(tt.mName, tt.h)
		got := g.Services[tt.wantService].GetHandleds()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
