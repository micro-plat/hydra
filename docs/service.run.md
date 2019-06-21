### 服务启动

一个应用程序实例可启动 6 种服务器的任意组合,只需使用`-`连接,可通过代码或命令行指定:

#### 1. 代码中指定

- 启动`api`,`rpc`服务器实例

```go
   hydra.WithServerTypes("api-rpc"),
```

- 启动`api`,`cron`,`mqc`服务器实例

```go
   hydra.WithServerTypes("api-rpc-mqc"),
```

#### 2. 命令行中指定

启动`api`和`rpc`实例

```sh
$ sudo ./helloserver run -r fs://../ -c test -S "api-rpc"
```

#### 3. 服务启动

可使用命令`run`和`start`启动服务,区别是:

> `run` 直接运行服务. 所有日志输出到控制台, 并根据级别显示不同颜色,便于调试,一般开发时使用此命令

> `start` 服务安装后可使用`start`命令启动, 服务将在在后台运行, 异常关闭或服务器重启会自动启动应用. 日志存入日志文件或远程日志归集系统, 控制台不显示日志. 可使用`stop`停止服务,`status`查看服务是否运行,`remove`卸载服务.

示例代码:

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
		hydra.WithSystemName("helloserver"), //系统名称
		hydra.WithDebug())

	app.Micro("/hello",hello)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```

使用`Micro`注册服务

重新编译`go install`并安装服务

- 修改任何配置,请重新执行`install`命令
  > 执行`install`时返回`Service has already been installed`错误,则需执行`remove`命令

```sh
$ sudo ./helloserver install -r fs://../ -c test -S "api-rpc"
Service has already been installed

$ sudo ./helloserver remove
Removing helloserver:					[  OK  ]
```

再次执行`install`命令

```sh
$ sudo ./helloserver install -r fs://../ -c test -S "api-rpc"
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		修改配置: /myplat_debug/helloserver/api/test/conf
		创建配置: /myplat_debug/helloserver/rpc/test/conf
Install helloserver:					[  OK  ]
```

- `run`启动服务

  一般`run`命令参数与`install`一致(`start`时不需要任何参数)

```sh
$ sudo ./helloserver run -r fs://../ -c test -S "api-rpc"

[2019/06/21 14:09:59.973885][i][fe0ccfe5b]初始化 /myplat_debug/helloserver/rpc-api/test
[2019/06/21 14:09:59.975250][i][629ded880]开始启动[RPC]服务...
[2019/06/21 14:09:59.975536][d][629ded880][未启用 jwt设置]
[2019/06/21 14:09:59.975541][d][629ded880][未启用 header设置]
[2019/06/21 14:09:59.975554][d][629ded880][未启用 metric设置]
[2019/06/21 14:09:59.975556][d][629ded880][未启用 host设置]
[2019/06/21 14:10:00.476655][i][629ded880]服务启动成功(RPC,tcp://192.168.4.121:8081,1)
[2019/06/21 14:10:00.476889][i][4eaf41f15]开始启动[API]服务...
[2019/06/21 14:10:00.477953][i][4eaf41f15][启用 静态文件]
[2019/06/21 14:10:00.477977][d][4eaf41f15][未启用 header设置]
[2019/06/21 14:10:00.477987][d][4eaf41f15][未启用 熔断设置]
[2019/06/21 14:10:00.477997][d][4eaf41f15][未启用 jwt设置]
[2019/06/21 14:10:00.478004][d][4eaf41f15][未启用 ajax请求限制设置]
[2019/06/21 14:10:00.478016][d][4eaf41f15][未启用 metric设置]
[2019/06/21 14:10:00.478023][d][4eaf41f15][未启用 host设置]
[2019/06/21 14:10:01.13949][i][4eaf41f15]服务启动成功(API,http://192.168.4.121:8090,1)


```

控制台打印出了两次`启动成功`,分别是`api`服务器(http 协议),`rpc`服务器(甚于 grpc,tcp 协议),包含服务提供地址和启动的服务个数

同一个服务器的日志可根据`session_id`(当前启动实例为:`8edb58733`,`7650a8ecf`)查看上下文日志

- `start`启动服务

```sh
$ sudo ./helloserver start
Starting helloserver:					[  OK  ]
```

控制台只会输出启动成功,不会显示运行时日志

#### 4. 服务发布信息查询

服务器启动时会自动将服务器及服务路径等添加到注册中心, 便于监控服务和服务发现者查找服务.

- 1. 服务器监控节点

  监控服务查询`/plat/system/[server-type]/cluster/servers`目录可获得正在运行的服务器

        `api`服务: `/myplat_debug/helloserver/api/test/servers/192.168.1.8:8090`

        `rpc`服务:`/myplat_debug/helloserver/rpc/test/servers/192.168.1.8:8081`

* 2. 服务提供者节点

     服务发布到`/plat/services/[server-type]/system/service-name/providers`目录,便于服务调用方获取服务:

     `api`服务: `/myplat_debug/services/api/helloserver/hello/providers/192.168.1.8:8090`

     `rpc`服务: `/myplat_debug/services/rpc/helloserver/hello/providers/192.168.1.8:8081`

     如使用`fs://../`指定的注册中心,则运行以下命令查看:

     ```sh
     $ cd ../myplat_debug/helloserver/api/test/servers/

     $ ls
     192.168.1.8:8090

     $ cat 192.168.1.8\:8090
     {"service":"http://192.168.1.8:8090"}
     ```
