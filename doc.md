# hydra

hydra 是基于 go 语言和众多开源项目实现的分布式微服务框架

hydra['haɪdrə]致力于提供统一，丰富的后端开发框架，降低后端开发的复杂性，提高开发效率。目前已支持的服务类型包括：`http api`服务，`rpc`服务，`websocket`,`mqc`消息消费服务，`cron`定时任务,`web`服务，静态文件服务。


特性


* 后端一体化框架, 支持6+服务器类型
* 微服务的基础设施, 服务注册发现，熔断降级，监控与配置管理
* 多集群模式支持，对等，主备等
* 20+线上项目实践经验
* 全golang原生实现



###  示例

- 1.  编写代码

新建文件夹`hello`,并添加`main.go`输入以下代码

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
		hydra.WithServerTypes("api"), //服务器类型为http api
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心
		hydra.WithDebug())

	app.API("/hello",hello)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```

- 2.  编译安装

```sh
go install hello

```
3. 安装服务
```sh
hello install
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