构建六大服务器
-------------

本示例介绍通过最简单的代码，构建服务器的方法

[TOC]

### 一、构建 api server

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



### 七、启动服务器

以编译后的二进制文件`apiserver`为例：

```sh
$ ./apiserver run --plat test
```

日志如下：
```sh
[2020/07/07 19:06:00.349344][i][3a9367088]初始化: /test/apiserver/api/1.0.0/conf
[2020/07/07 19:06:00.351695][i][3a9367088]启动[api]服务...
[2020/07/07 19:06:00.852165][i][dc0045470]启动成功(api,http://192.168.4.121:8081,1)
```