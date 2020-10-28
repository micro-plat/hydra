package services

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_serverHook_AddStarting(t *testing.T) {
	h := func(server.IServerConf) error {
		return nil
	}
	tests := []struct {
		name    string
		h       func(server.IServerConf) error
		wantErr bool
		ErrStr  string
	}{
		{name: "添加空服务", h: nil, wantErr: true, ErrStr: "启动服务不能为空"},
		{name: "添加服务", h: h},
		{name: "再次添加服务", h: h, wantErr: true, ErrStr: "启动服务不能重复注册"},
	}
	s := &serverHook{}
	for _, tt := range tests {
		err := s.AddStarting(tt.h)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.ErrStr, err.Error(), tt.name)
		}
	}
}

func Test_serverHook_AddClosing(t *testing.T) {
	h := func(server.IServerConf) error {
		return nil
	}
	tests := []struct {
		name    string
		h       func(server.IServerConf) error
		wantErr bool
		ErrStr  string
	}{
		{name: "添加空服务", h: nil, wantErr: true, ErrStr: "关闭服务不能为空"},
		{name: "添加服务", h: h},
		{name: "再次添加服务", h: h, wantErr: true, ErrStr: "关闭服务不能重复注册"},
	}
	s := &serverHook{}
	for _, tt := range tests {
		err := s.AddClosing(tt.h)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.Equal(t, tt.ErrStr, err.Error(), tt.name)
		}
	}
}

func Test_serverHook_AddHandleExecuting(t *testing.T) {
	tests := []struct {
		name       string
		h          []context.IHandler
		wantErr    bool
		wantHandle []context.IHandler
	}{
		{name: "添加空接口", h: nil, wantHandle: nil},
		{name: "添加单个接口", h: []context.IHandler{hander1{}}, wantHandle: []context.IHandler{hander1{}}},
		{name: "添加单个接口", h: []context.IHandler{hander1{}, hander2{}}, wantHandle: []context.IHandler{hander1{}, hander1{}, hander2{}}},
	}
	s := &serverHook{}
	for _, tt := range tests {
		err := s.AddHandleExecuting(tt.h...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		gotHandle := s.GetHandleExecutings()
		assert.Equal(t, tt.wantHandle, gotHandle, tt.name)
	}
}

func Test_serverHook_AddHandleExecuted(t *testing.T) {
	tests := []struct {
		name        string
		h           []context.IHandler
		wantErr     bool
		wantHandled []context.IHandler
	}{
		{name: "添加空接口", h: nil, wantHandled: nil},
		{name: "添加单个接口", h: []context.IHandler{hander1{}}, wantHandled: []context.IHandler{hander1{}}},
		{name: "添加单个接口", h: []context.IHandler{hander1{}, hander2{}}, wantHandled: []context.IHandler{hander1{}, hander1{}, hander2{}}},
	}
	s := &serverHook{}
	for _, tt := range tests {
		err := s.AddHandleExecuted(tt.h...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		gotHandled := s.GetHandleExecuteds()
		assert.Equal(t, tt.wantHandled, gotHandled, tt.name)
	}
}

func Test_serverHook_DoStarting(t *testing.T) {

	s := &serverHook{}
	err := s.DoStarting(nil)
	assert.Equal(t, nil, err, "启动空服务")

	//添加关闭服务
	err = s.AddStarting(func(server.IServerConf) error {
		return fmt.Errorf("test")
	})
	assert.Equal(t, false, err != nil, "添加服务")
	err = s.DoStarting(nil)
	assert.Equal(t, "test", err.Error(), "启动服务")
}

func Test_serverHook_DoClosing(t *testing.T) {
	s := &serverHook{}
	err := s.DoClosing(nil)
	assert.Equal(t, nil, err, "关闭空服务")

	//添加关闭服务
	err = s.AddClosing(func(server.IServerConf) error {
		return fmt.Errorf("test")
	})
	assert.Equal(t, false, err != nil, "关闭服务")
	err = s.DoClosing(nil)
	assert.Equal(t, "test", err.Error(), "关闭服务")
}
