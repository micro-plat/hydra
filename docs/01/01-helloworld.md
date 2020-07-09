创建项目
-------
本示例介绍如何创建一个简单的http api服务


[toc]


### 一. 新建项目

在`$gopath/src`下添加目录`apiserver`,并在目录下添加文件`main.go`，如下结构:

|----apiserver

|-------- main.go

### 二. 添加代码

```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {

    //创建app
	app := hydra.NewApp(
            hydra.WithServerTypes(http.API),
    )

    //注册服务
    app.API("/hello", func(ctx hydra.IContext) interface{} {
        return "hello world"
    })

    //启动app
    app.Start()
}
```

### 三. 运行服务

```sh
#编译apiserver
$ go build apiserver


#运行服务
$ ./apiserver run --plat test

[2020/05/11 11:14:54.772320][i][a665e52a0]初始化: /test/apiserver/api/1.0.0/conf
[2020/05/11 11:14:54.772420][i][a665e52a0]开始启动[api]服务...
[2020/05/11 11:14:55.273198][i][a665e52a0]服务启动成功(api,http://192.168.4.121:8080,1)
```

这样应用就在`8080`端口启动起来了。

### 四. 测试服务

```sh
$ curl http://192.168.4.121:8080/hello
```
返回内容：
```sh
hello world
```
