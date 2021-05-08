

## 背景介绍

通常面向业务的应用包括三大系统：外部服务系统、流程处理系统、运营管理系统。

外部服务系统： 为APP, H5,下游商户等提供交互能力，通常是通过http api或rpc提供服务

流程处理系统： 通常进行业务的内部处理，如：库存扣还、交易记账等或与上游供货商系统交互完成出库、发货等。通常包括：异步消息处理、定时任务处理等

运营管理系统： 目前主流采用前后端分离方式，前端使用vue react angular等框架，后端提供http或websocket接口





## 服务构建

### 服务类型
### 参数配置
### 服务注册


## 简单示例

### 创建项目

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




## 安装部署



### 编译启动

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

### 服务运行

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

```sh
$ curl http://192.168.4.121:8080/hello
```
返回内容：

```sh
hello world
``


* 1. 将服务安装到本地，以后台方式运行，服务器重启后自动启动:

```sh
./flowserver install
```
> `./flowserver remove` 卸载服务

2. 启动服务

```sh
./flowserver start
```
> `./flowserver stop` 停止服务
> `./flowserver status` 查看服务状态