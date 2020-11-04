package services

import (
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/context"
)

type handleHook struct {
	handlings map[string][]context.IHandler
	handleds  map[string][]context.IHandler
	closings  []func() error
}

func newHandleHook() *handleHook {
	return &handleHook{
		handlings: make(map[string][]context.IHandler),
		handleds:  make(map[string][]context.IHandler),
		closings:  make([]func() error, 0, 1),
	}
}

//AddHandling 添加handle预处理钩子
func (s *handleHook) AddHandling(service string, h ...context.IHandler) error {
	if len(h) == 0 {
		return nil
	}
	s.handlings[service] = append(s.handlings[service], h...)
	return nil
}

//AddHandled 添加handle后处理钩子
func (s *handleHook) AddHandled(service string, h ...context.IHandler) error {
	if len(h) == 0 {
		return nil
	}
	s.handleds[service] = append(s.handleds[service], h...)
	return nil
}

//AddClosingHanle 添加服务关闭钩子
func (s *handleHook) AddClosingHanle(c ...interface{}) error {
	if len(c) == 0 {
		return nil
	}
	for _, h := range c {
		if h == nil {
			continue
		}
		if vv, ok := h.(func() error); ok {
			s.closings = append(s.closings, vv)
			//return nil  //此时不该return  @hujun
		} else if vv, ok := h.(func()); ok {
			s.closings = append(s.closings, func() error {
				vv()
				return nil
			})
		} else {
			return fmt.Errorf("钩子函数签名不支持 %v", reflect.TypeOf(h))
		}
	}
	return nil

}

//GetHandlings 获取服务对应的预处理函数
func (s *handleHook) GetHandlings(service string) (h []context.IHandler) {
	return s.handlings[service]
}

//GetHandleds 获取服务对应的后处理后处理函数
func (s *handleHook) GetHandleds(service string) (h []context.IHandler) {
	return s.handleds[service]
}

//GetHandleClosingHandlers 获取提供有资源释放的服务
func (s *handleHook) GetClosingHandlers() []func() error {
	return s.closings
}
