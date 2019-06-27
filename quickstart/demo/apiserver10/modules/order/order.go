package order

import "github.com/micro-plat/hydra/component"

type IOrder interface {
}

type Order struct {
	c component.IContainer
}

func NewOrder(c component.IContainer) *Order {
	return &Order{
		c: c,
	}
}
