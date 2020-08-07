API服务器
----------------------

[TOC]

### 一、简单的服务器
```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
  
    app := hydra.NewApp(hydra.WithServerTypes(http.API))

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



### 二、服务器配置

API服务器可选的配置参数有：

|函数名|可选|说明|
|-----|:----:|:----|
|api.WithTrace|是|显示请求、响应的详细参数，默认不显示|
|api.WithTimeout|是|设置服务器读取请求与写入响应的超时时间长,默认30秒|
|api.WithHeaderReadTimeout|是|设置读取请求头的超时时长,默认30秒|
|api.WithHost|是|设置当前服务器的主机名,设置后只能通过主机名访问|
|api.WithDisable|是|停止当前服务器，默认false|
|api.WithEnable|是|启动当前服务器|
|api.WithDNS|是|设置是否发布到DNS服务节点，使用[DDNS](https://github.com/micro-plat/ddns)可直接通过域名访问到本机服务|

代码设置:
```go
hydra.Conf.API(":8081", api.WithTrace(),api.WithTimeout(5,5))
```


* 任何配置都可以操作注册中心添加、修改、删除。
* 固定不变的配置可以通过代码设置。
* 开发环境的配置可通过代码配置，但建议通过预编译命令对不同环境进行编译隔离。

### 三、中件间

|名称|说明|
|-----|:-----|
|Limiter|服务器限流配置|
|WhiteList|IP白名单，根据请求地址进行限制|
|BlackList|IP黑名单，指定后所有服务都不能访问|
|Delay|请求延迟处理组件|
|API KEY|指定固定密钥或通过RPC获取密钥，对请求的参数进行MD5,SHA1,SHA256等验证|
|Basic Auth|Basic Auth登录认证|
|JWT|JSON Web Token登录认证|
|RAS|远程认证系统验证|
|Header|可对响应header进行设置|
|Render|通过模板配置响应响应内容|


代码设置:
```go
hydra.OnReady(func() {
	hydra.Conf.API(":8080", api.WithTrace()).		
	WhiteListGroup(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
	BlackList("192.168.4.120").	
	Header()
})
```



### 四、示例


```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
    
    app := hydra.NewApp(hydra.WithServerTypes(http.API))

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

将以上代码编译为`apiserver`,并启动:

```sh
~/work/bin$ apiserver run -p plat
[2020/07/08 14:42:38.380283][i][882f188bc]初始化: /plat/apiserver/api/1.0.0/conf
[2020/07/08 14:42:38.381712][i][882f188bc]启动[api]服务...
[2020/07/08 14:42:38.881994][i][970cf0e3b]启动成功(api,http://192.168.4.121:8080,1)
```

请求服务:

```sh
curl http://localhost:8080/api
{"id":100,"name":"colin"}
```

服务器日志：


```sh
[2020/07/08 14:42:40.886998][i][883e3529d]api.request: GET /api?id=100 from 192.168.4.121
[2020/07/08 14:42:40.887060][i][883e3529d]--------api--------------
[2020/07/08 14:42:40.887203][i][883e3529d]api.response: GET /api?id=100 200  344.076µs
```