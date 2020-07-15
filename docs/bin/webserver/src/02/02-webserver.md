WEB服务器
----------------------
提供静态文件访问、API接口调用服务


[TOC]

### 一、简单的服务器
```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {  
    app := hydra.NewApp(hydra.WithServerTypes(http.WEB))
    app.Start()
}
```
编译后将`html静态文件`放到`./src`目录，通过`run`命令启动即可访问


WEB服务器同样支持api服务:

```go
app.WEB("/api",api)

func api(ctx hydra.IContext) interface{} {
    ctx.Log().Info("--------api--------------")
    return map[string]interface{}{
        "name":"colin",
        "id":ctx.Request().GetInt("id"),
    }
}
```

* 启动后可通过`/api`请求到服务。
* 对于前后端分离的WEB项目，“WEB服务器” 可以作为静态文件和后端接口的统一站点。



### 二、服务器配置

WEB服务器的配置参数与API服务器相同有：

|函数名|可选|说明|
|-----|:----:|:----|
|api.WithTrace|是|显示请求、响应的详细参数，默认不显示|
|api.WithTimeout|是|设置服务器读取请求与写入响应的超时时间长,默认30秒|
|api.WithHeaderReadTimeout|是|设置读取请求头的超时时长,默认30秒|
|api.WithHost|是|设置当前服务器的主机名,设置后只能通过主机名访问|
|api.WithDisable|是|停止当前服务器，默认false|
|api.WithEnable|是|启动当前服务器|
|api.WithDNS|是|设置是否发布到DNS服务节点，使用[DDNS](https://github.com/micro-plat/ddns)可直接通过域名访问到本机服务|


WEB服务器新增`Static`配置,可通过此配置，设置静态文件路径、扩展名、首页重写，甚至支持指定压缩包作为静态文件源。未指定时系统默认配置信息如下:

```json
{"archive":"webserver.zip","dir":"./src","exclude":["/views/",".exe",".so"],
"exts":[".txt",".html",".htm",".js",".css",".map",".ttf",".woff",".woff2",".woff2",".jpg",".jpeg",".png",".gif",".ico",".tif",".pcx",".tga",".exif",".fpx",".svg",".psd",".cdr",".pcd",".dxf",".ufo",".eps",".ai",".raw",".WMF",".webp"],
"first-page":"index.html",
"rewriters":["/","index.htm","default.html"]}
```
* archive：静态文件压缩包，默认查找与应用程序同名的zip压缩包，有则自动解压作为静态文件源。
* dir：静态文件存放的根目录
* exclude: 排除文件名或扩展名，请求路径中包含这些名称则不能访问
* exts: 支持的文件扩展名
* first-page:当访问`rewriters`指定的页面时，自动重写到当前页面
* rewriters:需要重写到首页的页面地址。即默认访问`index.htm`时，自动访问到`index.html`文件

以上配置实际上是通过以下代码生成的:

```go
   hydra.Conf.WEB.Static()
```
修改配置：
```go
   hydra.Conf.WEB.Static(static.WithArchive("./static.zip"))
```

将vuejs生成的静态文件直接放到`./src`目录或打包生成`"webserver.zip`放到与进程同级目录，即可访问。


### 三、中件间

静态文件访问控制支持以下组件
|名称|说明|
|-----|:-----|
|WhiteList|IP白名单，根据请求地址进行限制|
|BlackList|IP黑名单，指定后所有服务都不能访问|
|Delay|请求延迟处理组件|


API接口与API服务支持的组件相同：

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


中件间可通过代码设置:
```go
hydra.OnReady(func() {
	hydra.Conf.WEB(":8090", api.WithTrace()).		
	WhiteList(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
	BlackList(blacklist.WithIP("192.168.4.120")).	
	Header()
})
```
* 任何配置都可以操作注册中心添加、修改、删除。
* 固定不变的配置可以通过代码设置。
* 开发环境配置可通过代码设置，但建议通过预编译命令对不同环境进行编译隔离。

### 四、示例


```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
    app := hydra.NewApp(hydra.WithServerTypes(http.WEB))
    app.Start()
}
```

将以上代码编译为`webserver`,并将vuejs生成的静态文件打包生成[webserver.zip](./src/webserver.zip) 放到与webserver同级目录，并启动服务器：


```sh
~/work/bin$ webserver run -p plat
[2020/07/08 14:17:34.34813][i][7fafb0738]初始化: /plat/webserver/web/1.0.0/conf
[2020/07/08 14:17:34.36474][i][7fafb0738]启动[web]服务...
[2020/07/08 14:17:34.536734][i][8786eedf6]启动成功(web,http://192.168.4.121:8089)
```


通过浏览器访问 `http://localhost:8089`

![图片](../img/webserver--01.png)



服务器日志:
```sh
[2020/07/08 14:17:37.861767][i][e4060d3cc]web.request: GET / from 192.168.4.121
[2020/07/08 14:17:37.873062][i][e4060d3cc]web.response: GET / 200 static 11.537422ms
[2020/07/08 14:17:37.879381][i][6e1e2248a]web.request: GET /css/app.3c31093f.css from 192.168.4.121
[2020/07/08 14:17:37.879441][i][6e1e2248a]web.response: GET /css/app.3c31093f.css 200 static 185.822µs
[2020/07/08 14:17:37.879486][i][74db2d212]web.request: GET /js/chunk-vendors.6697755a.js from 192.168.4.121
[2020/07/08 14:17:37.879489][i][298d8efb3]web.request: GET /js/app.c6ab6812.js from 192.168.4.121
[2020/07/08 14:17:37.879585][i][298d8efb3]web.response: GET /js/app.c6ab6812.js 200 static 236.593µs
[2020/07/08 14:17:37.879622][i][74db2d212]web.response: GET /js/chunk-vendors.6697755a.js 200 static 277.615µs
[2020/07/08 14:17:37.880874][i][97a437623]web.request: GET /js/about.d211a758.js from 192.168.4.121
[2020/07/08 14:17:37.880954][i][97a437623]web.response: GET /js/about.d211a758.js 200 static 197.021µs
[2020/07/08 14:17:37.966417][i][ed30d2e25]web.request: GET /img/logo.82b9c7a5.png from 192.168.4.121
[2020/07/08 14:17:37.966532][i][ed30d2e25]web.response: GET /img/logo.82b9c7a5.png 200 static 227.871µs
[2020/07/08 14:17:37.981811][i][6f33a87b0]web.request: GET /favicon.ico from 192.168.4.121
[2020/07/08 14:17:37.981964][i][6f33a87b0]web.response: GET /favicon.ico 200 static 240.265µs
```
