启动原理
--------------------------
[TOC]

### 一、启动参数



|名称|必须|说明|
|---|---|----|
|platform|是|平台名称或项目名称，指一类系统的总称|
|system name|否|系统名称,默认为当前应用程序名称，如apiserver,flowserver,rpcserver|
|server type|是|服务器类型,指本程序提供的服务类型，包括api,web,rpc,ws,cron,mqc|
|cluster name|否|集群名称,默认为当前应用的版本号或hydra版本号|
|registry address|是|注册中心地址,默认为`lm://.`即本地内存作为注册中心。相当于以单机版运行|

* 以上参数可通过代码中指定:
```go
    app := hydra.NewApp(
        hydra.WithPlatName("taobao"),
        hydra.WithSystemName("rechargeServer"),
        hydra.WithClusterName("prod"),
        hydra.WithRegistry("zk://192.168.0.109"),
        hydra.WithServerTypes(rpc.RPC, http.API))
```
* 也可以在通过命令行指定:

可查看应用的`run`命令的请求参数
```sh
~/work/bin$ wsserver01 run -h
NAME:
   wsserver01 run - 运行服务

USAGE:
   wsserver01 run [command options] [arguments...]

OPTIONS:
   --registry value, -r value  *注册中心地址。目前支持zookeeper(zk)和本地文件系统(fs)。注册中心用于保存服务启动和运行参数，
                               服务注册与发现等数据，格式:proto://host。proto的取值有zk,fs; host的取值根据不同的注册中心各不同,
                               如zookeeper则为ip地址(加端口号),多个ip用逗号分隔,如:zk://192.168.0.2,192.168.0.107:12181。本地文
                               件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../;  此参数可以通过命令行参数指定，程序指
                               定，也可从环境变量中获取，环境变量名为: (default: "lm://.") [$registry]
   --plat value, -p value      *平台名称
   --system value, -s value    *系统名称 (default: "wsserver01")
   --cluster value, -c value   *集群名称 (default: "1.0.0")
   --trace value, -t value     -性能跟踪，可选项。用于生成golang的pprof的性能分析数据,支持的模式有:cpu,mem,block,mutex,web。其中web是以http
                                服务的方式提供pprof数据。
```
### 二、配置参数


### 三、配置更新

### 四、启动顺序

