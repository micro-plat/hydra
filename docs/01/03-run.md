启动原理
--------------------------
[TOC]

### 一、启动参数



|名称|必须|说明|
|---|---|----|
|platform|是|平台名称，或项目名称，指一类系统的总称|
|system name|否|系统名称,默认为当前应用程序名称，如apiserver,flowserver,rpcserver|
|server type|是|服务器类型,指本程序提供的服务类型，包括api,web,rpc,ws,cron,mqc|
|cluster name|否|集群名称,默认为当前应用的版本号或hydra版本号|
|registry address|是|注册中心地址,默认为`lm://.`即本地内存作为注册中心。相当于以单机版运行|

以上参数可通过代码中指定:
```go
    app := hydra.NewApp(
        hydra.WithPlatName("taobao"),
        hydra.WithSystemName("rechargeServer"),
        hydra.WithClusterName("prod"),
        hydra.WithRegistry("zk://192.168.0.109"),
        hydra.WithServerTypes(rpc.RPC, http.API))
```






### 二、配置参数
### 三、启动顺序
### 四、热更新
