package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver10/modules/const/keys"
)

type IOrderCache interface {
}

type OrderCache struct {
	c component.IContainer
}

func NewOrderCache(c component.IContainer) *OrderCache {
	return &OrderCache{
		c: c,
	}
}

func (o *OrderCache) Query(orderNO string) (string, error) {
	cache := o.c.GetRegularCache()
	rows, err := cache.Get(keys.ORDER_QUERY)
}
