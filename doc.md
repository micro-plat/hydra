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

以上代码通过`hydra.NewApp()`创建了 hydra app 实例， 并通过`app.Start()`启动服务。但要成功运行服务还需要指定以下参数：

1.  注册中心地址，支持 zookeeper(zk)和本地文件系统(fs)。用于保存服务启动和运行参数，服务注册与发现等数据，格式:proto://host。proto 的取值有 zk,fs; host 的取值各不相同,如 zookeeper 则为 ip 地址(加端口号),多个 ip 用逗号分隔:zk://192.168.0.2，192.168.0.107:12181。本地文件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../; 此参数可以通过命令行参数(--registry 或-r)指定,或通过环境变量($hydra_registry)指定,或在程序中通过 hydra.WithRegistry 指定

2.  完整名称, 当前应用在`注册中心`的路径，hydra 通过该路径从注册中心拉取配置数据。以`/`分隔的多级目录结构，完整的表示该服务所在平台，系统，服务
    类型，集群名称，格式：/平台名称/系统名称/服务器类型/集群名称; 平台名称，系统名称，集群名称可以是任意字母
    下划线或数字，服务器类型则为目前支持的几种服务器:api,web,rpc,mqc,cron,ws。该参数可以通过完整的路径方式传入`hydra.WithName`
    如: /平台名称/系统名称/服务器类型/集群名称，也分为 4 个字段分别传入： `hydra.WithPlatName`,`hydra.WithSystemName`,`hydra.WithServerTypes`,`hydra.WithClusterName`,`hydra.WithRegistry`， 这些参数都可以通过程序指定，命令行指定，或环境变量指定
