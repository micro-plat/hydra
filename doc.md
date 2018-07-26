# hydra 开发手册

## 一、 hydra 是什么

hydra(读音: ['haɪdrə])是一套构建后端服务，流程的快速开发框架。与其它框架不同的是`hydra`集成 api 接口， web 服务, rpc 服务, 定时任务, MQ 消费流程, websocket 等多种服务; 致力于快速开发、使用简单、功能丰富， 一次开发， 多服务运行等。最大限度解决重复造轮子，多种框架开发不统一等问题。

## 二、 起步

### 1. 安装 hydra

hydra 使用 godep 管理依赖包，执行`go get`即可下载完整源码

```sh
go get github.com/micro-plat/hydra
```

### 2. 示例项目

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

## 三、 hydra 服务创建

```go
app := hydra.NewApp()
app.Start()
```

以上代码通过`hydra.NewApp()`创建了 hydra app 实例， 并通过`app.Start()`运行该实例。但实例需要一些必须参数才能正常启动。

1.  注册中心地址，服务器通过`注册中心`管理`配置`，服务`注册`，`发现`等
2.  完整名称, 当前应用在`注册中心`的路径，hydra 通过该路径从注册中心拉取`配置数据`。该参数可以通过完整的路径方式传入`hydra.WithName`
    如:/平台名称/系统名称/服务器类型/集群名称，也分为 4 个字段分别传入。如 `hydra.WithPlatName`,`hydra.WithSystemName`,`hydra.WithServerTypes`,`hydra.WithClusterName`,`hydra.WithRegistry`, 其中`ServerTypes`可传`api`,`web`,`ws`,`cron`,`mqc`,也可以传多个服务类型，用`-`连接，如:`api-rpc-cron`。
    通过`hydra.With...`指定运行参数中是方法之一， 其它方式指定参数，请继续往下阅读。

我们可以通过以下三种方式指定运行参数：

1.  `hydra.NewApp()`初始化时通过`hydra.With...`参数指定
2.  通过命令行指定参数，如: app run --registry zk://192.168.0.168 -name /plat/sys/api/t
3.  通过环境变量指定参数,查看应用程序帮助信息 如 `app run -h` 获取各参数的环境变量名，配置完成后执行`app run`即可启动服务

考虑到生产环境环境实际情况,部分参数如：`PlatName`,`SystemName`,`ServerTypes`,可以通过`hydra.With...`指定; `注册中心地址`,`集群名称`可以在生产环境安装时指定; 便于开发环境通过简单的参数就可以启动系统， 可以将`注册中心地址`,`集群名称`配置到环境变量中。
