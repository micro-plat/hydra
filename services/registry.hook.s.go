package services

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
)

type serverHook struct {
	setup     func(app.IAPPConf) error
	starting  func(app.IAPPConf) error
	started   func(app.IAPPConf) error
	closing   func(app.IAPPConf) error
	handlings []context.IHandler
	handleds  []context.IHandler
}

//AddSetup 添加服务器启动前配置参数前执行勾子
func (s *serverHook) AddSetup(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("服务不能为空")
	}
	if s.setup != nil {
		return fmt.Errorf("服务不能重复注册")
	}
	s.setup = h
	return nil
}

//AddStarted 添加服务器启动完成勾子
func (s *serverHook) AddStarted(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("服务不能为空")
	}
	if s.started != nil {
		return fmt.Errorf("服务不能重复注册")
	}
	s.started = h
	return nil
}

//AddStarting 添加服务器启动勾子
func (s *serverHook) AddStarting(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("启动服务不能为空")
	}
	if s.starting != nil {
		return fmt.Errorf("启动服务不能重复注册")
	}
	s.starting = h
	return nil
}

//AddClosing 设置关闭器关闭勾子
func (s *serverHook) AddClosing(h func(app.IAPPConf) error) error {
	if h == nil {
		return fmt.Errorf("关闭服务不能为空")
	}
	if s.closing != nil {
		return fmt.Errorf("关闭服务不能重复注册")
	}
	s.closing = h
	return nil
}

//AddHandleExecuting 添加handle预处理勾子
func (s *serverHook) AddHandleExecuting(h ...context.IHandler) error {
	if len(h) == 0 {
		return nil
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
	if s.setup == nil {
		return nil
	}
	return s.setup(c)
}

//DoStarted 获取服务器启动完成处理函数
func (s *serverHook) DoStarted(c app.IAPPConf) error {
	if s.started == nil {
		return nil
	}
	return s.started(c)
}

//DoStarting 获取服务器启动预处理函数
func (s *serverHook) DoStarting(c app.IAPPConf) error {
	if s.starting == nil {
		return nil
	}
	return s.starting(c)
}

//DoClosing 获取服务器关闭处理函数
func (s *serverHook) DoClosing(c app.IAPPConf) error {
	if s.closing == nil {
		return nil
	}
	return s.closing(c)
}
