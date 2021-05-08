# 服务类型

目前支持六种服务类型, 为便于理解，统称为`ServerType`,多个`ServerType`可集成到同一个`app`应用中提供不同的服务。

 - 支持的`ServerType`:


| ServerType | 引用方式  | 说明                                         |
| :--------: | --------- | -------------------------------------------- |
|    API     | http.API  | 提供http服务                                 |
|    Web     | http.Web  | 提供Http服务                                 |
|     WS     | http.WS   | 提供websocket服务                            |
|    RPC     | rpc.RPC   | 提供基于grpc协议的通用远程调用方式           |
|    CRON    | cron.CRON | 提供定时任务服务，指定cron表达式定时执行任务 |
|    MQC     | mqc.MQC   | 提供消息消费服务，即message queue consumer   |


- 构建`app`:
```go
app := hydra.NewApp()
```

- 指定`ServerType`:

```go
app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithServerTypes(http.API,http.Web,cron,CRON,mqc.MQC,rpc.RPC,http.WS),
 )

```
同一个`app`支持多种ServerType，提供多种Service。


- 完整示例

```go
package main

import (
    "github.com/micro-plat/hydra"   
     "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithServerTypes(http.API),
    )
    app.Start()
}
```

此`app`未提供任何服务，服务实现与注册请参考后续章节。