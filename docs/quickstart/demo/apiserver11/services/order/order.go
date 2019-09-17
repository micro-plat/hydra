package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver10/modules/order"
)

type OrderHandler struct {
	container component.IContainer
	o         order.IOrder
}

func NewOrderHandler(container component.IContainer) (u *OrderHandler) {
	return &OrderHandler{
		container: container,
		o:         order.NewOrder(container),
	}
}

//RequestHandle 会员充值订单请求
func (u *OrderHandler) RequestHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------会员充值订单请求---------------")
	ctx.Log.Info("1.检查请求参数")
	if err := ctx.Request.Check("merchant_id","order_no", "account", "face", "num"); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}

	ctx.Log.Info("2. 创建充值订单")
	result, err := u.o.Create(
		ctx.Request.GetString("merchant_id")
		ctx.Request.GetString("order_no"),
		ctx.Request.GetString("account"),
		ctx.Request.GetInt("face"),
		ctx.Request.GetInt("num"))
	if err != nil {
		return err
	}
	return result
}

//QueryHandle 充值结果查询
func (u *OrderHandler) QueryHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------充值结果查询---------------")
	ctx.Log.Info("1.检查请求参数")
	if err := ctx.Request.Check("merchant_id","order_no"); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}

	ctx.Log.Info("2. 查询充值结果")
	result, err := u.o.Query(
		ctx.Request.GetString("merchant_id",
		ctx.Request.GetString("order_no"))
	if err != nil {
		return err
	}
	return result
}
