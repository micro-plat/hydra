# 构建API服务器

本示例构建一个api服务器，为商户提供会员充值，结果查询两个接口。接口通过md5签名验证。

| 接口地址       | 功能说明     |
| -------------- | ------------ |
| /order/request | 充值订单提交 |
| /order/query   | 充值结果查询 |

知识点:

*  数据库配置与操作
*  消息队列配置与发送
*  自定义服务名
*  md5签名

#### 1. 搭建基础代码


`main.go`
```go
package main

import "github.com/micro-plat/hydra/hydra"

type apiserver struct {
	*hydra.MicroApp
}

func main() {
	app := &apiserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("apiserver"),
			hydra.WithServerTypes("api")),
	}
	app.Start()
}

```
`config.dev.go`
```go
// +build !prod

package main

func (api *apiserver) config() {
	api.IsDebug = true
	api.Conf.API.SetMainConf(`{"address":":8090","trace":true}`)
	api.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"mysql",
			"connString":"mrss:123456@tcp(192.168.0.36)/mrss?charset=utf8",
			"maxOpen":20,
			"maxIdle":10,
			"lifeTime":600		
	}`)
	api.Conf.Plat.SetVarConf("queue", "queue", `
		{
			"proto":"redis",
			"addrs":[
					"192.168.0.111:6379",
					"192.168.0.112:6379",
					"192.168.0.113:6379"
			],
			"db":1,
			"dial_timeout":10,
			"read_timeout":10,
			"write_timeout":10,
			"pool_size":10
	}
	`)
}
```


`config.prod.go`

```go

// +build prod

package main

func (api *apiserver) config() {
	api.Conf.API.SetMainConf(`{"address":":8090","trace":true}`)
	api.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"mysql",
			"connString":"#connString",
			"maxOpen":20,
			"maxIdle":10,
			"lifeTime":600		
	}`)
	api.Conf.Plat.SetVarConf("queue", "queue", `
		{
			"proto":"redis",
			"addrs":[
					#redis_addr
			],
			"db":1,
			"dial_timeout":10,
			"read_timeout":10,
			"write_timeout":10,
			"pool_size":20
	}
	`)
}
```


#### 2. 初始化检查与服务注册



```go
package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/services/order"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (api *apiserver) init() {
	api.Initializing(func(c component.IContainer) error {
		//检查db配置是否正确
		if _, err := c.GetDB(); err != nil {
			return err
		}
		if _, err := c.GetQueue(); err != nil {
			return err
		}	

		return nil
	})

	//服务注册
	api.Micro("/order", order.NewOrderHandler)
}
```

#### 3. 请求预处理，验证签名

```go
package main

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/modules/merchant"
)

func (api *apiserver) handling() {
	api.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {
		if err := ctx.Request.Check("merchant_id"); err != nil {
			return err
		}
		key, err := merchant.GetKey(ctx,ctx.Request.GetInt(merchant_id))
		if err != nil {
			return err
		}
		if !ctx.Request.CheckSign(key) {
			return fmt.Errorf(908, "商户签名错误")
		}
		return nil
	})
}
```


#### 3. 构建服务


`servers/order.go`

```go
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
```


#### 4. 业务逻辑
```go
package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/modules/const/keys"
	"github.com/micro-plat/qtask"
)

type IOrder interface {
	Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error)
	Query(merchantID string, orderNO string) (map[string]interface{}, error)
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
func (d *OrderDB) Create(merchantID string, orderNO string, account string, face int, num int) (map[string]interface{}, error) {
	order, err := d.db.Create(merchantID, orderNO, account, face, num)
	if err != nil {
		return nil, err
	}
	qtask.Create(d.c, "订单支付", order, 60, keys.ORDER_PAY)
	return order, err
}



```



```go
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



```