package services

import "github.com/micro-plat/hydra/context"

type unit struct {
	path     string
	service  string
	Handling context.IHandler
	Handled  context.IHandler
	Handle   context.IHandler
	Fallback context.IHandler
	actions  []string
	group    *unitGroup
}

func (u *unit) GetHandlings() []context.IHandler {
	h := make([]context.IHandler, 0, 1)
	if u.group.Handling != nil {
		h = append(h, u.group.Handling)
	}
	if u.Handling != nil {
		h = append(h, u.Handling)
	}
	return h
}
func (u *unit) GetHandleds() []context.IHandler {
	h := make([]context.IHandler, 0, 1)
	if u.Handled != nil {
		h = append(h, u.Handled)
	}
	if u.group.Handled != nil {
		h = append(h, u.group.Handled)
	}

	return h
}
