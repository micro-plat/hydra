
# 简单示例


##  创建项目

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


初始化go mod 文件:``` go mod init```

编译代码:```  go build ```



## 调试运行


运行服务:```$ ./flowserver run```


直接运行服务，方便查看日志，调试代码，终端日志如下：

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

```$ curl http://192.168.4.121:8080/hello```

返回内容：
```hello world```