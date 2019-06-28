package order

import (
	"fmt"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/modules/const/sqls"
)

type IOrderDB interface {
	Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error)
	Query(merchantID string, orderNO string) (map[string]interface{}, error)
}

type OrderDB struct {
	c component.IContainer
}

func NewOrderDB(c component.IContainer) *OrderDB {
	return &OrderDB{
		c: c,
	}
}
func (d *OrderDB) Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error) {
	db := d.c.GetDB()
	input := map[string]interface{}{
		"merchant_id": merchantID,
		"order_no":    orderNO,
		"account":     account,
		"face":        face,
		"num":         num,
	}
	orderID, _, _, err := db.Scalar(sqls.Get_ORDER_ID, input)
	if err != nil {
		return nil, err
	}
	input["order_id"] = orderID
	row, _, _, err := db.Execute(sqls.ORDER_CREATE, input)
	if err != nil || row == 0 {
		return nil, fmt.Errorf("系统错误暂时无法创建订单%v", err)
	}
	return map[string]interface{}{
		"order_id": orderID,
	}, nil
}

func (d *OrderDB) Query(merchantID string, orderNO string) (map[string]interface{}, error) {
	db := d.c.GetDB()
	row, _, _, err := db.Execute(sqls.ORDER_QUERY, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return row, nil
}
