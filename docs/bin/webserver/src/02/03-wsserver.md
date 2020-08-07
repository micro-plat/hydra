Web Socket服务器
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
  
    app := hydra.NewApp(hydra.WithServerTypes(http.WS))

    //注册服务
    app.WS("/api", api)
    
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

### 二、服务器配置

WS服务器可选的配置参数有：

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
hydra.Conf.WS(":8081", api.WithTrace(),api.WithTimeout(5,5))
```


* 任何配置都可以操作注册中心添加、修改、删除。
* 固定不变的配置可以通过代码设置。
* 开发环境的配置可通过代码配置，但建议通过预编译命令对不同环境进行编译隔离。

### 三、中件间

|名称|说明|
|-----|:-----|
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
	hydra.Conf.WS(":8080", api.WithTrace()).		
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
    
    app := hydra.NewApp(hydra.WithServerTypes(http.WS))

    //注册服务
    app.WS("/api", api)
    
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

将以上代码编译为`wsserver`,并启动:

```sh
~/work/bin$ wsserver run -p plat
[2020/07/08 19:12:29.863859][i][5a00a15db]初始化: /plat/wsserver/ws/1.0.0/conf
[2020/07/08 19:12:29.865100][i][5a00a15db]启动[ws]服务...
[2020/07/08 19:12:30.365351][i][56cf22bfa]启动成功(ws,ws://192.168.4.121:8070,1)

```

请求服务:
可通过[在线web socket测试工具](http://www.websocket-test.com/)连接 ws://localhost:8070 与服务器建立连接
通过json发送服务请求:
```json
{"service":"/api","id":100}
```

![在线发送](../img/ws01.png)

222.209.84.37:8121为localhost:8070的公网访问入口

* service 参数为必须参数，指定后端处理的服务名称，其它参数为业务参数


服务器日志：


```sh
[2020/07/09 14:27:31.59993][i][a7c3fa907]ws.request: GET /ws from 125.69.28.75
[2020/07/09 14:27:45.185950][i][a7c3fa907]ws.request: GET /api from 125.69.28.75
[2020/07/09 14:27:45.185977][i][a7c3fa907]--------api--------------
[2020/07/09 14:27:45.186039][i][a7c3fa907]ws.response: GET /api 200 ws 94.211µs
```