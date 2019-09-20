package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver10/modules/const/sqls"
	"github.com/micro-plat/lib4go/db"
)

type IOrderDB interface {
}

type OrderDB struct {
	c component.IContainer
}

func NewOrderDB(c component.IContainer) *OrderDB {
	return &OrderDB{
		c: c,
	}
}
func (o *OrderDB) Query(orderNO string) (db.QueryRows, error) {
	db := o.c.GetRegularDB()
	rows, _, _, err := db.Query(sqls.ORDER_QUERY, map[string]interface{}{
		"order_no": orderNO,
	})
	return rows, err
}
