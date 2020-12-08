package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
)

func Test_handleHook_AddHandling(t *testing.T) {
	tests := []struct {
		name        string
		service     string
		h           []context.IHandler
		wantErr     bool
		wantHanding []context.IHandler
	}{
		{name: "1.1 添加空接口", service: "service", wantHanding: nil},
		{name: "1.2 添加非空接口", service: "service", h: []context.IHandler{hander1{}}, wantHanding: []context.IHandler{hander1{}}},
		{name: "1.3 添加不存在接口", service: "service", h: []context.IHandler{hander2{}}, wantHanding: []context.IHandler{hander1{}, hander2{}}},
		{name: "1.4 添加已存在理接口", service: "service", h: []context.IHandler{hander2{}}, wantHanding: []context.IHandler{hander1{}, hander2{}, hander2{}}},
	}
	s := newHandleHook()
	for _, tt := range tests {
		err := s.AddHandling(tt.service, tt.h...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		got := s.GetHandlings(tt.service)
		assert.Equal(t, tt.wantHanding, got, tt.name)
	}
}

func Test_handleHook_AddHandled(t *testing.T) {
	tests := []struct {
		name        string
		service     string
		h           []context.IHandler
		wantErr     bool
		wantHandled []context.IHandler
	}{
		{name: "1.1 添加空接口", service: "service", wantHandled: nil},
		{name: "1.2 添加非空接口", service: "service", h: []context.IHandler{hander1{}}, wantHandled: []context.IHandler{hander1{}}},
		{name: "1.3 添加不存在接口", service: "service", h: []context.IHandler{hander2{}}, wantHandled: []context.IHandler{hander1{}, hander2{}}},
		{name: "1.4 添加已存在理接口", service: "service", h: []context.IHandler{hander2{}}, wantHandled: []context.IHandler{hander1{}, hander2{}, hander2{}}},
	}
	s := newHandleHook()
	for _, tt := range tests {
		err := s.AddHandled(tt.service, tt.h...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		got := s.GetHandleds(tt.service)
		assert.Equal(t, tt.wantHandled, got, tt.name)
	}
}

func Test_handleHook_AddClosingHanle(t *testing.T) {
	f1 := func() error { return nil }
	f2 := func() {}
	tests := []struct {
		name           string
		c              []interface{}
		wantErr        bool
		wantClosingLen int
	}{
		{name: "1.1 添加空钩子"},
		{name: "1.2 添加不支持类型的钩子", c: []interface{}{"123456"}, wantErr: true},
		{name: "1.3 添加钩子", c: []interface{}{f1}, wantClosingLen: 1},
		{name: "1.4 添加多个钩子", c: []interface{}{f1, nil, f2}, wantClosingLen: 3},
	}
	s := newHandleHook()
	for _, tt := range tests {
		err := s.AddClosingHanle(tt.c...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		got := s.GetClosingHandlers()
		assert.Equal(t, tt.wantClosingLen, len(got), tt.name)
	}
}
