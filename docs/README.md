

## 快速开始

  

### 1. 创建项目

创建`flowserver`,对外提供api服务，内部每隔5秒自动执行某个服务。


新建文件`main.go`,结构如下:

|----flowserver

|-------- main.go


编写如下代码：
```go
package main

import (
    "github.com/micro-plat/hydra"   
     "github.com/micro-plat/hydra/hydra/servers/http"
     "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {

	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithServerTypes(http.API,cron.CRON),
    )

    app.API("/hello",hello)
    app.CRON("/auto",hello,"@every 5s") 
    app.Start()
}
func hello(ctx hydra.IContext) interface{} {
        return "hello world"
}
```

### 2. 编译启动

初始化go mod 文件:

```sh
    go mod init
```

编译代码:
```sh
    go build

```

启动服务:

```sh
$ ./flowserver run
```

### 3. 服务运行

终端日志如下：

```sh
[2020/07/08 09:36:31.140432][i][29f63e41d]初始化: /test/flowserver/api-cron/1.0.0/conf
[2020/07/08 09:36:31.143027][i][29f63e41d]启动[api]服务...
[2020/07/08 09:36:31.643524][i][b65655312]启动成功(api,http://192.168.4.121:8080,1)
[2020/07/08 09:36:31.643885][i][29f63e41d]启动[cron]服务...
[2020/07/08 09:36:31.844844][i][3908a5ccc]启动成功(cron,cron://192.168.4.121,1)
[2020/07/08 09:36:32.346047][d][3908a5ccc]the cron server is started as master
[2020/07/08 09:36:36.648149][i][01751ece6]cron.request: GET /hello from 192.168.4.121
[2020/07/08 09:36:36.648244][i][01751ece6]cron.response: GET /hello 200  193.356µs
[2020/07/08 09:36:41.651858][i][00f45e17b]cron.request: GET /auto from 192.168.4.121
[2020/07/08 09:36:41.651911][i][00f45e17b]cron.response: GET /auto 200  159.694µs
```

服务已启动

> api服务： 默认8080端口，对外提供 `/hello`服务

> cron服务: 每隔5秒自动执行服务`/auto`


测试api服务：

```sh
$ curl http://192.168.4.121:8080/hello
```
返回内容：

```sh
hello world
```




## 服务构建


所有服务器都采用相同的方式构建，即`指定服务类型`,`注册服务`,`配置运行参数`

### 1. 服务类型

构建api server
```go    
	app := hydra.NewApp(hydra.WithServerTypes(http.API))
```
构建web server
```go   
	app := hydra.NewApp(hydra.WithServerTypes(http.Web))
```
构建RPC server
```go   
	app := hydra.NewApp( hydra.WithServerTypes(http.RPC))
```
构建CRON server
```go
	app := hydra.NewApp( hydra.WithServerTypes(cron.CRON))
```
构建MQC server
```go    
	app := hydra.NewApp(hydra.WithServerTypes(mqc.MQC))
```
### 2. 注册服务
服务器提供的服务，都采用相同的方式进行注册，同一个服务`Handler`可注册到任何服务器,可将函数、struct等通过注册接口进行注册

Handler示例：
```go
//Query 订单查询
func Query(ctx hydra.IContext) interface{} {
	ctx.Log().Debug("-------------处理订单查询----------------------")
	if err := ctx.Request().Check(fields.FieldMerNo, fields.FieldMerOrderNo); err != nil {
		return err
	}

	ctx.Log().Debug("1. 查询订单信息")
	order, err := orders.QueryDetail(ctx.Request().GetString(fields.FieldMerNo),
		ctx.Request().GetString(fields.FieldMerOrderNo))
	if err == nil && order.Len() > 0 {
		return order
	}

	ctx.Log().Debug("2. 订单不存在")
	return errs.NewError(int(enums.CodeOrderNotExists), "订单不存在")
}
```

注册为API服务：
```go
hydra.S.API("/order/request", Query)
```

注册为Web服务：
```go
hydra.S.Web("/order/request", Query)
```

注册为RPC服务：
```go
hydra.S.RPC("/order/request", Query)
```

注册为CRON服务：
```go
hydra.S.CRON("/order/request", Query)
```

注册为MQC服务：
```go
hydra.S.MQC("/order/request", Query)
```
