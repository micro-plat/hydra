package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/modules/const/keys"
	"github.com/micro-plat/lib4go/db"
	// "github.com/micro-plat/qtask"
)

type IOrder interface {
	Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error)
	Query(merchantID string, orderNO string) (db.QueryRows, error)
}

type Order struct {
	c  component.IContainer
	db IOrderDB
}

func NewOrder(c component.IContainer) *Order {
	return &Order{
		c:  c,
		db: NewOrderDB(c),
	}
}
func (d *Order) Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error) {
	order, err := d.db.Create(merchantID, orderNO, account, face, num)
	if err != nil {
		return nil, err
	}
	// qtask.Create(d.c, "订单支付", order, 60, keys.ORDER_PAY)
	return order, err
}
