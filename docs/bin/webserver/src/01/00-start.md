hydra
======================

后端一站式服务框架，提供统一的编码方式开发API、WEB、WEBSOCKET，RPC、CRON、MQC服务器。

具有如下特点：

- 1. 服务化
相同的编码方式编写服务，注册到六大服务器即可运行

- 2. 统一配置
本地零配置，注册中心集中管理，提供配置热更新

- 3. 守护进程
集成systemd,systemv,upstart等daemon本地服务管理。通过install,start,stop,remove等命令管理本地进程


- 4. 日志归集
跟踪请求调用链、提供远程日志归集

- 5. 业务监控
对QPS、并发数、处理时长、响应统计等进行采集并存入influxdb。使用grafana配置图表即可监控系统状态

- 6. 安全认证
提供白名单、黑名单、Basic Auth、API KEY、JWT、RAS等访问控制


- 7. 灰度控制
提供根据规则判断请求参数，进行集群转发，用于灰度发布。

- 8. 集群模式
提供多集群发布及对等、分片、主从集群模式

- 9.  服务治理
提供服务注册、服务发现、负载均衡、流量控制、服务降级、熔断等

- 10. 基础组件
提供缓存、数据库、消息队列、远程调用(RPC)、分布式锁、全局编号等组件工具

 hydra已运行在20+线上项目中


#### 一、api server示例


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

#### 二、包含多个服务器

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
        return "hello world"
}
```


编译后的二进制文件`flowserver`为例：

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


```sh
$ curl http://192.168.4.121:8080/hello
```
返回内容：
```sh
hello world
```
