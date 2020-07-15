参数配置
--------------------------
[TOC]

### 一、启动参数

hydra实现了本地零配置。所有服务器配置都通过`注册中心`集中管理，但需要指定`注册中心地址`，`平台名`、`系统名`、`服务器类型`和`集群名`。启动时连接到`注册中心`，拉取`/平台名/系统名/服务器类型/集群名称/conf`服务器参数，启动服务器。


|名称|必须|说明|
|---|---|----|
|platform|是|平台名称或项目名称|
|system name|否|系统名称或应用名, 默认为当前应用程序名称，如apiserver,flowserver,rpcserver|
|server type|是|服务器类型,指 本程序提供的服务类型，包括api,web,rpc,ws,cron,mqc。每个程序可以包含多个服务类型|
|cluster name|否|集群名称, 默认为当前应用的版本号或hydra版本号|
|registry address|否|注册中心地址, 默认为`lm://.`即本地内存作为注册中心，相当于以单机版运行|

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

* 可以在代码中指定：平台名称、系统名称、服务器类型等确定的，不会变化的参数。命令行中指定集群名称、注册中心地址等随环境发生变化的参数。


* 应用启动后会持续监听`/平台名/系统名/服务器类型/集群名称/conf`节点，未配置或配置发生变化后将拉取配置信息，并判断首次启动还是配置变更，关键配置变更需要重启服务器。非关键参数变更立即生效不会重启服务器。


### 二、服务器参数
指服务器启动时和运行时的配置参数，除MQC服务器需指定消息队列服务器外，其它服务器都有默认参数，不用配置即可启动。配置可通过代码设置或三方工具连接到注册中心设置。

通过代码设置或hydra的默认配置，必须调用`[应用名] conf install`命令，安装到注册中心才能使用(lm://.本地内存注册中心除外)。

注册中心配置可通过`[应用名] conf show`命令查看:
```sh
~/work/bin$ mqcserver conf show -p plat
└─plat
  └─mqcserver
    └─mqc
      └─1.0.0
        └─conf
          └─main[1]
          └─queue[2]
    └─api
      └─1.0.0
        └─conf
          └─router[3]
          └─main[4]
请输入数字序号 > 
```
输入配置名后面的数字序号可查看配置内容




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
修改zookeeper配置
![zookeeper](../img/settings01.png)



### 三、应用参数
如数据库连接、消息队列、缓存等配置信息，放到`/平台名/var/类型/名称`目录下，如:`/hydra/var/db/db`,`/hydra/var/mq/redis`


* 代码设置：
```go
hydra.OnReady(func() {
     hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
})
```

* 三方工具设置:
修改zookeeper配置
![zookeeper](../img/settings02.png)