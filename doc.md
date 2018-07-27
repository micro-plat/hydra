# hydra 开发手册

## 一、 hydra 是什么

hydra(读音: ['haɪdrə])是一套构建后端服务，流程的快速开发框架。

通常后端系统包括接口服务和自动流程。接口服务有`http接口`，`rpc接口`，自动流程有`定时任务`，`MQ消息消费`。另外系统间实时推送消息还会用到`websocket`，前端静态网页还需要`web站点`来运行。这些服务器`hydra`都提供。

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

## 三、 hydra 服务构建

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

| 参数名称     | 必须 | 程序传入              | 命令行参数           | 环境变量           | 说明                                   |
| ------------ | ---- | --------------------- | -------------------- | ------------------ | -------------------------------------- |
| 注册中心地址 | √    | hydra.WithRegistry    | --registry 或-r      | hydra_registry     | proto://host                           |
| 系统全名     | √    | hydra.WithName        | --name 或 -n         | hydra_name         | /平台名称/系统名称/服务器类型/集群名称 |
| 平台名称     | √    | hydra.WithPlatName    | --plat 或-p          | hydra_plat         | 字母下划线或数字                       |
| 系统名称     | √    | hydra.WithSystemName  | --system 或-s        | hydra_system       | 字母下划线或数字                       |
| 服务类型     | √    | hydra.WithServerTypes | --server-types 或 -S | hydra_server_types | api,web,rpc,mqc,cron,ws                |
| 集群名称     | √    | hydra.WithClusterName | --cluster 或 -c      | hydra_cluster      | 字母下划线或数字                       |
| 调试模式     | ×    | hydra.WithDebug       | ---                  | ---                | ---                                    |
| 性能跟踪     | ×    | ---                   | --trace 或-t         | hydra_trace        | cpu,mem,block,mutex,web                |
| 远程日志     | ×    | ---                   | --rpclog 或-r        | hydra_rpclog       | ---                                    |
| 远程服务     | ×    | ---                   | --rs,-R              | hydra-rs           | ---                                    |

了解了启动参数，我们尝试启动程序:

```sh
~/work/bin$ sudo app run --registry zk://192.168.0.107 --name /myapp/sys/api/t
[2018/07/26 18:31:06.497813][i][c46cf206e]Connected to 192.168.0.107:2181
[2018/07/26 18:31:06.512927][i][c46cf206e]Re-submitting `0` credentials after reconnect
[2018/07/26 18:31:06.512919][i][c46cf206e]Authenticated: id=100437066149265444, timeout=4000
[2018/07/26 18:31:07.497153][i][c46cf206e]初始化 /myapp/sys/api/t
[2018/07/26 18:31:17.498113][w][c46cf206e]/myapp/sys/api/t 未配置
```

我们看到服务已成功连接到了 zookeeper 服务器，但提示`/myapp/sys/api/t`未配置。这是因为服务器运行的端口号等信息并未指定，服务器无法启动。我们可使用`sudo app install`命令来初始化`启动参数`

```sh
:~/work/bin$ sudo app install --registry zk://192.168.0.107 --name /myapp/sys/api/t
	-> 创建注册中心配置数据?,如果存在则不修改(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):1
	   创建配置: /myapp/sys/api/t/conf
Install app:					[  OK  ]
yanglei@yanglei-H97-HD3:~/work/bin$
```

注意：

1.  `sudo app install`我们使用了和`sudo app run`相同的参数启动。因为`install`同样需要知道`注册中心地址`，`需要初始化的服务地址`

2.  创建注册中心的数据我们未指定任何参数（如服务端口号等）,是采用默认参数进行配置。

提示 `Install app: [OK]`则安装成功，我们再次运行`sudo app run`

```sh
~/work/bin$ sudo app run --registry zk://192.168.0.107 --name /app/sys/api/t
[2018/07/26 18:31:41.960876][i][9ba7a63ae]Connected to 192.168.0.107:2181
[2018/07/26 18:31:41.974092][i][9ba7a63ae]Re-submitting `0` credentials after reconnect
[2018/07/26 18:31:41.974087][i][9ba7a63ae]Authenticated: id=100437066149265445, timeout=4000
[2018/07/26 18:31:42.960822][i][9ba7a63ae]初始化 /app/sys/api/t
[2018/07/26 18:31:42.963070][i][4c4f2ec06]开始启动...
[2018/07/26 18:31:42.963349][d][4c4f2ec06][未启用 header设置]
[2018/07/26 18:31:42.963344][i][4c4f2ec06][启用 静态文件]
[2018/07/26 18:31:42.963365][d][4c4f2ec06][未启用 jwt设置]
[2018/07/26 18:31:42.963370][d][4c4f2ec06][未启用 metric设置]
[2018/07/26 18:31:42.963367][d][4c4f2ec06][未启用 ajax请求限制设置]
[2018/07/26 18:31:42.963362][d][4c4f2ec06][未启用 熔断设置]
[2018/07/26 18:31:42.963372][d][4c4f2ec06][未启用 host设置]
[2018/07/26 18:31:43.499808][i][4c4f2ec06]启动成功(http://192.168.5.71:8090,0)
```

看到`启动成功...`,恭喜你服务成功启动了。但是后面有一个`0`，是因为我们还没注册服务。

## 四、 服务注册
