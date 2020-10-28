package services

import (
	"github.com/micro-plat/hydra/context"
)

type Unit struct {
	Path     string
	Service  string
	Handling context.IHandler
	Handled  context.IHandler
	Handle   context.IHandler
	Fallback context.IHandler
	Actions  []string
	Group    *UnitGroup
}

//GetHandlings 获取所有预处理函数
func (u *Unit) GetHandlings() []context.IHandler {
	h := make([]context.IHandler, 0, 1)
	if u.Group.Handling != nil {
		h = append(h, u.Group.Handling)
	}
	if u.Handling != nil {
		h = append(h, u.Handling)
	}
	return h
}

//GetHandleds 获取所有后处理函数
func (u *Unit) GetHandleds() []context.IHandler {
	h := make([]context.IHandler, 0, 1)
	if u.Handled != nil {
		h = append(h, u.Handled)
	}
	if u.Group.Handled != nil {
		h = append(h, u.Group.Handled)
	}

	return h
}
