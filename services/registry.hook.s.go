package services

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
)

type serverHook struct {
	setups    []func(app.IAPPConf) error
	startings []func(app.IAPPConf) error
	starteds  []func(app.IAPPConf) error
	closings  []func(app.IAPPConf) error
	handlings []context.IHandler
	handleds  []context.IHandler
}

func newServerHook() *serverHook {
	return &serverHook{
		setups:    make([]func(app.IAPPConf) error, 0, 1),
		startings: make([]func(app.IAPPConf) error, 0, 1),
		starteds:  make([]func(app.IAPPConf) error, 0, 1),
		closings:  make([]func(app.IAPPConf) error, 0, 1),
		handlings: make([]context.IHandler, 0, 1),
		handleds:  make([]context.IHandler, 0, 1),
	}
}

//AddSetup 添加服务器启动前配置参数前执行勾子
func (s *serverHook) AddSetup(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("服务不能为空")
	}
	s.setups = append(s.setups, h)
	return nil
}

//AddStarted 添加服务器启动完成勾子
func (s *serverHook) AddStarted(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("服务不能为空")
	}
	s.starteds = append(s.starteds, h)
	return nil
}

//AddStarting 添加服务器启动勾子
func (s *serverHook) AddStarting(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("启动服务不能为空")
	}
	s.startings = append(s.startings, h)
	return nil
}

//AddClosing 设置关闭器关闭勾子
func (s *serverHook) AddClosing(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("关闭服务不能为空")
	}
	s.closings = append(s.closings, h)
	return nil
}

//AddHandleExecuting 添加handle预处理勾子
func (s *serverHook) AddHandleExecuting(h ...context.IHandler) error {
	if len(h) == 0 {
		return nil
	}
	for _, v := range h {
		s.Remove(v)
	}
	s.handlings = append(s.handlings, h...)
	return nil
}

//AddHandleExecuted 添加handle后处理勾子
func (s *serverHook) AddHandleExecuted(h ...context.IHandler) error {
	if len(h) == 0 {
		return nil
	}
	s.handleds = append(s.handleds, h...)
	return nil
}

//GetHandleExecutings 获取handle预处理勾子
func (s *serverHook) GetHandleExecutings() []context.IHandler {
	return s.handlings
}

//GetHandleExecuteds 获取handle后处理勾子
func (s *serverHook) GetHandleExecuteds() []context.IHandler {
	return s.handleds
}

//DoSetup 获取服务器配置勾子
func (s *serverHook) DoSetup(c app.IAPPConf) error {
	for _, setup := range s.setups {
		if err := setup(c); err != nil {
			return err
		}
	}
	return nil
}

//DoStarted 获取服务器启动完成处理函数
func (s *serverHook) DoStarted(c app.IAPPConf) error {
	for _, started := range s.starteds {
		if err := started(c); err != nil {
			return err
		}
	}
	return nil
}

//DoStarting 获取服务器启动预处理函数
func (s *serverHook) DoStarting(c app.IAPPConf) error {
	for _, starting := range s.startings {
		if err := starting(c); err != nil {
			return err
		}
	}
	return nil
}

//DoClosing 获取服务器关闭处理函数
func (s *serverHook) DoClosing(c app.IAPPConf) error {
	for _, closing := range s.closings {
		if err := closing(c); err != nil {
			return err
		}
	}
	return nil
}

//Remove 移除某个钩子函数
func (s *serverHook) Remove(hs ...interface{}) {
	for _, h := range hs {
		for i, p := range s.closings {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.closings = append(s.closings[:i], s.closings[i+1:]...)
			}
		}
		for i, p := range s.handleds {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.handleds = append(s.handleds[:i], s.handleds[i+1:]...)
			}
		}
		for i, p := range s.handlings {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.handlings = append(s.handlings[:i], s.handlings[i+1:]...)
			}
		}
		for i, p := range s.setups {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.setups = append(s.setups[:i], s.setups[i+1:]...)
			}
		}
		for i, p := range s.starteds {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.starteds = append(s.starteds[:i], s.starteds[i+1:]...)
			}
		}
		for i, p := range s.startings {
			if strings.EqualFold(fmt.Sprintf("%p", p), fmt.Sprintf("%p", h)) {
				s.startings = append(s.startings[:i], s.startings[i+1:]...)
			}
		}
	}
}
