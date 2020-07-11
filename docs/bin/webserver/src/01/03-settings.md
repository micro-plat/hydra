服务器配置
--------------------------
[TOC]

### 一、启动参数

hydra实现了本地零配置。所有服务器配置都通过`注册中心`集中管理，但是在启动前需要指定`注册中心地址`，`平台名`、`系统名`、`服务器类型`和`集群名`。启动时连接到`注册中心`，并从注册中心拉取`/平台名/系统名/服务器类型/集群名称/conf`服务器参数，启动服务器。


|名称|必须|说明|
|---|---|----|
|platform|是|平台名称或项目名称，指一类系统的总称|
|system name|否|系统名称或应用名,默认为当前应用程序名称，如apiserver,flowserver,rpcserver|
|server type|是|服务器类型,指本程序提供的服务类型，包括api,web,rpc,ws,cron,mqc。每个应用可以包含多个服务器|
|cluster name|否|集群名称,默认为当前应用的版本号或hydra版本号|
|registry address|否|注册中心地址,默认为`lm://.`即本地内存作为注册中心。相当于以单机版运行|

* 启动参数可通过代码指定:
```go
    app := hydra.NewApp(
        hydra.WithPlatName("taobao"),
        hydra.WithSystemName("rechargeServer"),
        hydra.WithClusterName("prod"),
        hydra.WithRegistry("zk://192.168.0.109"),
        hydra.WithServerTypes(rpc.RPC, http.API))
```


* 也可以通过命令行指定:

```go
package main

import "github.com/micro-plat/hydra"

func main() {
	app := hydra.NewApp()
	app.Start()
}
```

可查看应用的`run`命令的请求参数
```sh
~/work/bin$ apiserver run -h
NAME:
   apiserver run - 运行服务

USAGE:
   apiserver run [command options] [arguments...]

OPTIONS:
   --registry value, -r value  *注册中心地址。目前支持zookeeper(zk)和本地文件系统(fs)。注册中心用于保存服务启动和运行参数，
                               服务注册与发现等数据，格式:proto://host。proto的取值有zk,fs; host的取值根据不同的注册中心各不同,
                               如zookeeper则为ip地址(加端口号),多个ip用逗号分隔,如:zk://192.168.0.2,192.168.0.107:12181。本地文
                               件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../;  此参数可以通过命令行参数指定，程序指
                               定，也可从环境变量中获取，环境变量名为: (default: "lm://.") [$registry]
   --name value, -n value      *服务全名，指服务在注册中心的完整名称，该名称是以/分隔的多级目录结构，完整的表示该服务所在平台，系统，服务
                               类型，集群名称，格式：/平台名称/系统名称/服务器类型/集群名称; 平台名称，系统名称，集群名称可以是任意字母
                               下划线或数字，服务器类型则为目前支持的几种服务器类型有:api,web,rpc,mqc,cron,ws。该参数可从环境变量中获取，
                               环境变量名为: [$name]
   --trace value, -t value     -性能跟踪，可选项。用于生成golang的pprof的性能分析数据,支持的模式有:cpu,mem,block,mutex,web。其中web是以http
                                服务的方式提供pprof数据。

```

* 也可以混合指定，代码中指定不能变化的如：平台名称、系统名称、服务器类型。命令行中指定集群名称、注册中心地址。


* 应用启动后会持续监听`/平台名/系统名/服务器类型/集群名称/conf`节点，未配置或配置发生变化后将拉取配置信息，并判断首次启动或关闭服务器重启。


### 二、服务器配置
指服务器启动时和运行时的配置参数，除MQC服务器需指定消息队列服务器外，其它服务器都有默认参数，不用配置即可启动。其它配置如白名单、黑名单、安全认证等运行时配置可通过代码、或三方工具连接到注册中心设置。

通过代码指定的配置或hydra默认配置，必须调用`conf install`命令，安装到注册中心才能使用(使用lm://.本地内存作为注册中心除外)。

服务器配置根路径：`/平台名/系统名/服务器类型/集群名称/conf`
如:`/hydra/rechargeSystem/api/prod/conf`

服务器其它配置都在根节点以下，如：
路由：`/hydra/rechargeSystem/api/prod/conf/router/`
限流配置:`/hydra/rechargeSystem/api/prod/conf/acl/limit`
JWT:`/hydra/rechargeSystem/api/prod/conf/auth/jwt`

* 代码设置：
```go
hydra.OnReady(func() {
      hydra.Conf.API(":8080", api.WithTrace()).
      WhiteList(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
      BlackList(blacklist.WithIP("192.168.4.120")).	
      Header()	
})
```

* 三方工具设置:
修改zookeeper配置为例
![zookeeper](../src/settings01.png)



### 三、应用配置
指在应用逻辑中使用的配置，如数据库连接、消息队列服务器地址等作为平台的公共参数，放到`/平台名/var/类型/名称`目录下，如:`/hydra/var/db/db`,`/hydra/var/mq/redis`。同样可以通过代码或三方工具进行设置。


* 代码设置：
```go
hydra.OnReady(func() {
     hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
})
```

* 三方工具设置:
修改zookeeper配置为例
![zookeeper](../src/settings02.png)