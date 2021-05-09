# 服务实现


服务的实现须采用统一的参数签名，任何`Service`可注册到任意的`ServerType`，有两种方式实现服务：


### 函数方式

函数参数为: `func(hydra.IContext)interface{}` 或 `func(hydra.IContext)`

示例:

```go
func hello(ctx hydra.IContext)interface{}{
    return "success"
}
```

或

```go
func hello(ctx hydra.IContext){
    ctx.Response().JSON(200,"success")
}
```

### 对象方式

定义`struct`并提供一个或多个以`Handle`名称结尾，函数签名为 `func(hydra.IContext)interface{}` 或 `func(hydra.IContext)`的服务处理函数。


```go
type OrderHandler struct{

}
func NewOrderHandler()*OrderHandler{
    return &OrderHandler{}
}

func(o *OrderHandler)QueryHandle(ctx hydra.IContext)interface{}{
    return "success"
}

func(o *OrderHandler)RequestHandle(ctx hydra.IContext)interface{}{
    return "success"
}

func(o *OrderHandler)GetHandle(ctx hydra.IContext)interface{}{
    return "success"
}
func(o *OrderHandler)PostHandle(ctx hydra.IContext)interface{}{
    return "success"
}
func(o *OrderHandler)PutHandle(ctx hydra.IContext)interface{}{
    return "success"
}
func(o *OrderHandler)DeleteHandle(ctx hydra.IContext)interface{}{
    return "success"
}
```
以上实现了6个服务，服务实际名称在注册时指定。

其中，`QueryHandle`,`RequestHandle`为普通服务，支持`GET`，`POST`请求

`GetHandle`，`PostHandle`，`PutHandle`，`DeleteHandle`为`resutful`服务，对应请求的`GET`,`POST`,`PUT`,`DELETE` Method。


> 使用`对象方式`实现`Service`实际是将一组服务进行`Group`, 但于对`统一状态`,`资源释放`，`勾子函数`等进行统一处理。
