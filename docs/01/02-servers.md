构建六大服务器
-------------

本示例通过最简单的代码，介绍构建6大服务器及混合服务器的方法

[TOC]

### 一、构建api server

对外提供http api 服务

```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
  
	app := hydra.NewApp(
            hydra.WithServerTypes(http.API),
    )

    //注册服务
    app.API("/hello", hello)
    
    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
对外提供 http://localhost:8080/hello 服务，通过`GET、POST`等请求时返回结果为`success`

### 二、构建mqc  server

监听消息队列，有新的消息时执行对应的服务

```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/mqc"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(mqc.MQC),
    )

    //注册服务
    app.MQC("/hello",hello,"queue-name")

    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
监听消息队列`queue-name`，有新消息到达时执行服务`/hello`即`hello`函数



### 三、构建cron  server

提供定时任务服务器

```go
package main

import (
    "github.com/micro-plat/hydra"
   
     "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(cron.CRON),
    )
   
   //注册服务
    app.CRON("/hello",hello,"@every 5s") 

  
    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
每隔5秒执行一次"/hello"服务



### 四、构建rpc server

提供基于json的rpc协议服务器

```go
package main

import (
    "github.com/micro-plat/hydra"
   
     "github.com/micro-plat/hydra/hydra/servers/rpc"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(rpc.RPC),
    )

    //注册服务
    app.RPC("/hello",hello)

    app.Start()
}
func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
对外提供基于RPC协议的服务，可通过hydra.C.RPC组件进行调用

### 五、构建ws server

提供web socket 服务器

```go
package main

import (
    "github.com/micro-plat/hydra"
   
     "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(http.WS),
    )

    //注册服务
    app.WS("/hello",hello)
     
    app.Start()
}


func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
对外提供`/ws`服务，外部可通过json方式请求服务`/hello`


### 六、构建web server

提供静态文件及API服务

```go
package main

import (
    "github.com/micro-plat/hydra"
   
     "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(http.WEB),
    )

    //配置静态文件
    hydra.Conf.WEB.Static(static.WithArchive("./static.zip")) //系统自动解压static.zip自动路由到包中对应的文件

    //注册服务
    app.WEB("/hello",hello)

    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```

外部可访问`static.zip`中包含的所有静态文件和`/hello`服务



以上6种类型服务器都以服务注册的方式将服务实现对象注册到服务器，并以相同的方式启动即可使用。


### 七、混合服务器

同一个程序中包含多个服务器

```go
package main

import (
    "github.com/micro-plat/hydra"   
     "github.com/micro-plat/hydra/hydra/servers/http"
     "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(http.API,cron.CRON),
    )

    app.API("/hello",hello)

    app.CRON("/hello",hello,"@every 5s") 

    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```


以编译后的二进制文件`flowserver`为例：

```sh
$ ./flowserver run --plat test
```

日志如下：
```sh
[2020/07/08 09:36:31.140432][i][29f63e41d]初始化: /test/flowserver/api-cron/1.0.0/conf
[2020/07/08 09:36:31.143027][i][29f63e41d]启动[api]服务...
[2020/07/08 09:36:31.643524][i][b65655312]启动成功(api,http://192.168.4.121:8080,1)
[2020/07/08 09:36:31.643885][i][29f63e41d]启动[cron]服务...
[2020/07/08 09:36:31.844844][i][3908a5ccc]启动成功(cron,cron://192.168.4.121,1)
[2020/07/08 09:36:32.346047][d][3908a5ccc]this cron server is started as master
[2020/07/08 09:36:36.648149][i][01751ece6]cron.request: GET /hello from 192.168.4.121
[2020/07/08 09:36:36.648244][i][01751ece6]cron.response: GET /hello 200  193.356µs
[2020/07/08 09:36:41.651858][i][00f45e17b]cron.request: GET /hello from 192.168.4.121
[2020/07/08 09:36:41.651911][i][00f45e17b]cron.response: GET /hello 200  159.694µs
```

服务已经启动，cron server也开始按周期(每隔5秒执行一次服务`/hello`)执行任务