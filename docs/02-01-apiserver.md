API服务器
----------------------

[TOC]

### 一、最简单的服务器
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
    app.API("/api", api)
    
    app.Start()
}

func api(ctx hydra.IContext) interface{} {
    ctx.Log().Info("--------api--------------")
    return map[string]interface{}{
        "name":"colin",
        "id":ctx.Request().GetInt("id"),
    }
}
```
打印日志信息,返回对象为map则api服务器自动转换为json返回

[>>  服务请求--了解更多关于输入参数的获取、验证处理](./03-04-req.md)
[>>  服务响应--了解更多关于输出参数、响应码处理](./03-05-resp.md)
[>>  模板输出--了解更多关于模板方式转换输出响应](./03-04-render.md)



### 二、服务器高级配置

API服务器可选的配置参数有：

|函数名|可选|说明|
|-----|:----:|:----|
|api.WithTrace|是|是否显示输入、转出参数信息，默认false|
|api.WithTimeout|是|设置服务器读取请求与写入响应的超时时间长,默认30秒|
|api.WithHeaderReadTimeout|是|设置读取http的超时时长,默认30秒|
|api.WithHost|是|设置当前服务器的主机名,设置后只能通过主机名访问|
|api.WithDisable|是|停止当前服务器，默认false|
|api.WithEnable|是|启动当前服务器|
|api.WithDNS|是|设置是否发布到DNS服务节点，使用[DDNS](https://github.com/micro-plat/ddns)可直接通过域名访问到本机服务|

列如设置请求与响应时分别输出详细的参数信息:
```go
hydra.Conf.API(":8081", api.WithTrace())
```
设置后的日志如下:


```sh
curl http://localhost:8080/api
```

```sh
[2020/07/07 19:38:43.824752][i][2456d0644]api.request: GET /api?id=100 from 192.168.4.121
[2020/07/07 19:38:43.824798][d][2456d0644]> trace.request: map[id:100]
[2020/07/07 19:38:43.824829][i][2456d0644]------request------
[2020/07/07 19:38:43.824922][d][2456d0644]> trace.response: 200 {"name":"colin","id":100}
[2020/07/07 19:38:43.824941][i][2456d0644]api.response: GET /api?id=100 200 trace 192.799µs
```
