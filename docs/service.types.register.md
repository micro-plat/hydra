## 服务注册

hydra 已支持 6 种服务器类型:`http api`服务，`rpc`服务，`websocket`服务,`mqc`消息消费服务，`cron`定时任务,`web`服务. 分别对应的服务器类型名为:`api`,`rpc`,`ws`,`mqc`,`cron`,`web`

#### 1. 注册函数

`hydra`实例提供了 8 个服务注册函数, 可注册到不同的服务器,见下表:

| 注册函数 | api | rpc | web | ws  | mqc | cron |
| -------- | --- | --- | --- | --- | --- | ---- |
| Micro    | √   | √   | √   | ×   | ×   | ×    |
| Flow     | ×   | ×   | ×   | ×   | √   | √    |
| API      | √   | ×   | ×   | ×   | ×   | ×    |
| RPC      | ×   | √   | ×   | ×   | ×   | ×    |
| WEB      | ×   | ×   | √   | ×   | ×   | ×    |
| WS       | ×   | ×   | ×   | √   | ×   | ×    |
| MQC      | ×   | ×   | ×   | ×   | √   | ×    |
| CRON     | ×   | ×   | ×   | ×   | ×   | √    |

示例:

```go
    app.API("/hello",hello)
    app.MQC("/hello",hello)

    func hello(ctx *context.Context) (r interface{}) {
	    return "hello world"
    }
```

服务支持两种类型注册:

- 1.  函数注册: 服务实现代码放在函数中,函数签名格式为:`(*context.Context) (interface{})`,示例:

```go
        func hello(ctx *context.Context) (r interface{}) {
            return "hello world"
    }
```

- 2.  实例注册: 服务实现代码放到`struct`中,传入`struct`实例的引用的构造函数

      示例:

```go
           app.API("/hello",token.NewQueryHandler)
```

          添加服务实现文件`query.handler.go`

```go

            package token

            import (
                "github.com/micro-plat/hydra/component"
                "github.com/micro-plat/hydra/context"
            )

            type QueryHandler struct {
                container component.IContainer
            }


            //NewQueryHandler 创建服务
            func NewQueryHandler(container component.IContainer) (u *QueryHandler) {
                return &QueryHandler{
                    container: container,
                }
            }
            func (u *QueryHandler) Handle(ctx *context.Context) (r interface{}) {
                var result struct {
                    ErrCode int64  `json:"errcode"`
                    ErrMsg  string `json:"errmsg"`
                }
                result.ErrCode = 0
                result.ErrMsg = "success"
                return result
            }
```

该`struct`需具备两个条件:

1. 服务构造函数`NewQueryHandler`, 只能有两种格式:

   `(container component.IContainer) (*QueryHandler)`

   或

   `(container component.IContainer) (*QueryHandler,error)`

2. 对象中至少包含一个命名为`...Handle`的函数,且签名为:
   `(*context.Context) (interface{})`格式.

#### 2. 服务名称

```go
    app.API("/order",order.NewOrderHandler)
```

第一个参数`/order`为服务名, 一般都以`/`开头,支持`/`分隔的多段名称如:

```go
    app.API("/order/request",order.NewOrderRequestHandler)
```

第二个参数`order.NewOrderHandler`为服务实现函数

请求的服务名一般与注册的服务器名一致, 但服务注册函数返回的是引用`实例`,且内部实现的函数名为`xxxHandle`签名为`(*context.Context) (interface{})`的函数时,请求的服务器为`注册名`+`/`+`函数名`
如:

```go
    app.API("/order",order.NewOrderHandler)
```

```go
package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type OrderHandler struct {
	container component.IContainer
}

func NewOrderHandler(container component.IContainer) (u *QueryHandler) {
	return &OrderHandler{container: container}
}

func (u *OrderHandler) RequestHandle(ctx *context.Context) (r interface{}) {
	return "success"
}
func (u *OrderHandler) QueryHandle(ctx *context.Context) (r interface{}) {
	return "success"
}
```

以上示例实际注册了两个服务:
`/order/request`,`/order/query`,分别对应`RequestHandle`,
`QueryHandle`服务处理函数

#### 3. 生命周期

服务是在服务器初始化时挂载的，外部请求到达时直接执行服务名对应的服务处理函数(`Handle`,`...Handle`)，服务实例可实现`func Close()error`函数，用于释放服务相关资源。服务器关闭时会自动调用每个服务已实现的`Close`函数。

服务函数尽量不要依赖全局资源，必须依赖时应充分考虑多个服务器启，停对该资源的影响。

如：

```go
type Input struct {
	ID   string `form:"id" json:"id" valid:"int,required"` //绑定输入参数，并验证类型否是否必须输入
	Name string `form:"name" json:"name"`
}
type BindHandler struct {
	container component.IContainer
}

func NewBindHandler(container component.IContainer) (u *BindHandler) {
	return &BindHandler{container: container}
}
func (u *BindHandler) GetHandle(ctx *context.Context) (r interface{}) {
	var input Input
	if err := ctx.Request.Bind(&input); err != nil {
		return err
	}
	return input
}
func (u *BindHandler) Close()error{
    return nil
}

```

服务注册代码:

```go
app := hydra.NewApp(
		hydra.WithPlatName("hydra-test"),
		hydra.WithSystemName("micro"),
		hydra.WithServerTypes("api-rpc"),
		hydra.WithDebug())
	app.Micro("/order/bind",NewBindHandler)
	app.Start()
```

1. `api`服务器和`rpc`服务器启动时会分别执行`NewBindHandler`创建两个`BindHandler`实例

2. 某一个服务器关闭时(通过注册中心配置关闭),会调用每一个服务实例的`Close`函数(假如实现了`Close`函数),如当前示例的`BindHandler.Close`函数

> 全局数据保存与获取，可使用`component.IContainer`中提供的`SaveGlobalObject`和`GetGlobalObject`函数。
