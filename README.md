# hydra

hydra 是基于 go 语言和众多开源项目实现的分布式微服务框架

通常后端系统包括接口服务和自动流程。接口服务有`http接口`，`rpc接口`，自动流程有`定时任务`，`MQ消息消费`。另外系统间实时推送消息还会用到`websocket`，前端静态网页还需要`web服务器`来运行。每种服务的架构模式不同,编码风格各不相同,接口服务化就需服务治理，分布式部署就需要对等，主备，分片模式支持。服务器数量达到一定规模还需要集群监控,日志归集。`hydra`致力于解决这些痛点，搭建统一框架，统一开发模式，持续完善基础设施; 提供快速开发、使用简单、功能强大、轻量级、易扩展的基础框架; 打造敏捷开发，快速交付的企业级后端服务开发方案

hydra 的优势

1.  使用简单

    > 只需几行代码便可集成到自己的应用，并具备 hydra 丰富的功能；扔到服务器通过命令启动起来即可，后台运行，自动重启。

2.  易学习，适合快速开发

    > 标准化输入输出`request`,`response`，统一公共组件`db`,`cache`，`mq`,`logger`,`metric`,`registry`等。所有服务的代码都相同， 只需注册为`api接口`，`web服务`,`rpc服务`,`定时任务`,`MQ消费流程`,`websocket`等运行即可。

3.  分布式服务协调

    > 自带`zookeeper`,`文件系统`作为`配置管理中心`和`注册中心`。集中管理配置信息，变更后自动下发到各服务器， 服务自动`注册`与`发现`，`监控`, 管理服务以不同模式运行(`主`，`备`)和处理流程任务分片等。 服务器初始化完成后`配置中心`，`注册中心`宕机不影响服务运行。恢复后自动连接到`注册中心`，`配置中心`

4.  服务器初始化与安装

    > 提供安装程序 可`配置中心数据`，`数据库`，`本地配置`等

5.  本地服务化运行

    > 以本地服务方式运行 `后台运行`，`自动重启`

6.  服务监控

    > 服务状态`cpu,内存，QPS,服务执行时长，执行结果统计`等自动上报到 influxdb,通过 grafana 实时查看状态

7.  统一日志服务

    > 为执行流程分配统一的`全局编号`，通过全局编号串联所有日志，可根据`全局编号`查询位于不同服务器， 不同流程服务的所有日志。并可以配置远程日志服务，自动上传日志。

8.  丰富功能

    > 支持服务`熔断`，`降级`,`RESTful`,`jwt`,`智能路由`，`静态文件服务`，支持数据库`oracle`,`mysql`,`sqlite`,`sqlserver`,`postgreSQL`,`influxdb` 缓存 `redis`,`memcache`,`本地缓存`,消息队列 `activeMQ`,`mqtt`,`redis`，websocket `消息推送`，定时任务和 MQ 消息消费以`主备模式，分片模式，对等模式运行`等

9.  易扩展
    > `数据库`，`缓存`，`消息队列`，`注册中心`等组件可自已扩展

## 安装 hydra

hydra 使用 godep 管理依赖包，执行`go get`即可下载完整源码

```sh
go get github.com/micro-plat/hydra
```

## 示例项目

示例代码请参考[examples](https://github.com/micro-plat/hydra/tree/master/examples)

### 1. hello world 示例代码

- 1.  编写代码

新建项目`hello`,并添加`main.go`文件输入以下代码，实现一个简单的`http api服务`

```go
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心
		hydra.WithDebug())

	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```

- 2.  编译安装

```sh
go install hello
```

- 3.  运行服务

```sh
./hello run
```

- 4.  测试服务

```sh
curl http://localhost:8090/hello

{"data":"hello world"}
```
