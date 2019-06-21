## API 服务器示例

#### 一、返回类型

```go

package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("helloserver"),
		hydra.WithServerTypes("api"), //服务器类型为http api
		hydra.WithDebug(),
	)

	app.Micro("/order", NewOrderHandler)
	app.Start()
}

type order struct{
    OrderNo string `json:"order_no" xml:"order_no"`
    ProductID int `json:"pid" xml:"pid"`
    Amount int `json:"amount" xml:"amount"`
}

type OrderHandler struct {
	container component.IContainer
}

func NewOrderHandler(container component.IContainer) (u *OrderHandler) {
	return &OrderHandler{container: container}
}
func (u *OrderHandler) RequestHandle(ctx *context.Context) (r interface{}) {
    tp:=ctx.Request.GetString("result_type")
    switch(tp){
        case "xml":
        ctx.Response.SetXML()　//指定content-type为text/xml
        case "json":
        ctx.Response.SetJSON() //指定content-type为text/json
        default: //默认content-type为text/json
    }
    return &order{
        OrderNo:"89769037987666",
        ProductID:5989234,
        Amount:10,
    }
}
```

安装服务

```sh

~/work/bin$ sudo ./tserver install -r zk://192.168.0.109
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		创建配置: /myplat_debug/helloserver/api/t/conf
```

运行服务

```sh
:~/work/bin$ sudo ./tserver run -r zk://192.168.0.109
[2019/06/21 14:24:40.387765][i][4a41971bb]Connected to 192.168.0.109:2181
[2019/06/21 14:24:40.391548][i][4a41971bb]Authenticated: id=246395503264333900, timeout=4000
[2019/06/21 14:24:40.391555][i][4a41971bb]Re-submitting `0` credentials after reconnect
[2019/06/21 14:24:40.436415][i][4a41971bb]初始化 /myplat_debug/helloserver/api/t
[2019/06/21 14:24:40.446795][i][ca27bca03]开始启动[API]服务...
[2019/06/21 14:24:40.447940][i][ca27bca03][启用 静态文件]
[2019/06/21 14:24:40.447958][d][ca27bca03][未启用 header设置]
[2019/06/21 14:24:40.447967][d][ca27bca03][未启用 熔断设置]
[2019/06/21 14:24:40.447975][d][ca27bca03][未启用 jwt设置]
[2019/06/21 14:24:40.447980][d][ca27bca03][未启用 ajax请求限制设置]
[2019/06/21 14:24:40.447989][d][ca27bca03][未启用 metric设置]
[2019/06/21 14:24:40.447996][d][ca27bca03][未启用 host设置]
[2019/06/21 14:24:40.965265][i][ca27bca03]服务启动成功(API,http://192.168.4.121:8090,1)

```

测试

```sh
~/work/bin$ curl http://localhost:8090/order/request?result_type=xml
<order><order_no>89769037987666</order_no><pid>5989234</pid><amount>10</amount></order>


~/work/bin$ curl http://localhost:8090/order/request?result_type=json
{"order_no":"89769037987666","pid":5989234,"amount":10}

~/work/bin$ curl http://localhost:8090/order/request
{"order_no":"89769037987666","pid":5989234,"amount":10}
```

#### 二、参数配置

```go
func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("helloserver"),
		hydra.WithServerTypes("api"), //服务器类型为http api
		hydra.WithDebug(),
    )
    app.Conf.API.SetMainConf(`{"address":":9098","trace":true,"rTimeout":"30","wTimeout":30}`)

	app.Micro("/order", NewOrderHandler)
	app.Start()
}
```

> 参数说明

- `address`: 服务器启动地址，合法格式为:`[ip或localhost]:port`

- `trace`:显示 response 响应信息

> 未设置 trace 参数或 trace 参数值为 false

```sh
[2019/06/21 14:44:33.115273][i][158ebae99]api.response GET /order/request 200  276.616µs
```

> trace 参数值设置为 true

```sh
[2019/06/21 14:44:33.115273][i][158ebae99]api.response GET /order/request 200  276.616µs {"order_no":"89769037987666","pid":5989234,"amount":10}

```

\*http 请求超时时长

> rTimeout，wTimeout：http 请求读写超时时长,默认 10 秒
