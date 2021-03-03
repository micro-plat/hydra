package services

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/assert"
)

// func Test_metaServices_AddHanler(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		service string
// 		h       context.IHandler
// 		wantErr bool
// 	}{
// 		{name: "1.1 service和handler为空", service: "", h: nil},
// 		{name: "1.2 service和handler不为空", service: "service", h: hander1{}},
// 		{name: "1.3 更新已存在service", service: "service", wantErr: true},
// 	}
// 	s := newService()
// 	services := []string{}
// 	for _, tt := range tests {
// 		err := s.AddHanler(tt.service, "", tt.h, nil)
// 		assert.Equal(t, tt.wantErr, err != nil, tt.name)
// 		if tt.wantErr {
// 			continue
// 		}
// 		services = append(services, tt.service)
// 		gotH, ok := s.GetHandlers(tt.service)
// 		assert.Equal(t, true, ok, tt.name)
// 		assert.Equal(t, tt.h, gotH, tt.name)

// 		gotS := s.GetServices()
// 		assert.Equal(t, services, gotS, tt.name)
// 	}
// }

func Test_metaServices_AddFallback(t *testing.T) {
	tests := []struct {
		name    string
		service string
		h       context.IHandler
		wantErr bool
	}{
		{name: "1.1 service和handler为空", service: "", h: nil},
		{name: "1.2 添加service和handler", service: "service", h: hander1{}},
		{name: "1.3 更新已存在service", service: "service", h: hander2{}},
	}
	s := newService()
	for _, tt := range tests {
		err := s.AddFallback(tt.service, tt.h)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.h == nil {
			continue
		}
		got, ok := s.GetFallback(tt.service)
		assert.Equal(t, true, ok, tt.name)
		assert.Equal(t, tt.h, got, tt.name)
	}
}
