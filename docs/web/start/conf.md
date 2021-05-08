
# 参数配置

指服务启动、运行所必须的配置。

这些配置使用注册中心集中管理，本地零配置。

配置发生变化后自动通知到集群中的各应用并立即生效，必要时应用会自动重启。

有三种方式构建配置：


### 代码构建
可通过`hydra.Conf.[ServerType]`提供的函数对服务参数与服务组件进行构建，如：
```go
hydra.Conf.API("6689")
```
#### 1. 支持链式调用与交互参数输入(hydra.ByInstall):

```go
//go:embed loginweb/dist/static
var staticFs embed.FS
var archive = "loginweb/dist/static"

hydra.Conf.Web("6687", api.WithTimeout(300, 300), api.WithDNS(hydra.ByInstall)).
			Static(static.WithAutoRewrite(), static.WithEmbed(archive, staticFs)).
			Processor(processor.WithServicePrefix("/web")).
			Header(header.WithCrossDomain()).
			Jwt(jwt.WithMode("HS512"),
				jwt.WithSecret("f0abd74b09bcc61449d66ae5d8128c18"),
				jwt.WithExpireAt(36000),
				jwt.WithAuthURL("/"),
				jwt.WithHeader(),
				jwt.WithExcludes(
					"/*/system/config/get",
					"/*/member/login",
					"/*/member/bind/*",
					"/*/member/sendcode",
					"/*/logout",
				),
			)

```

> hydra.ByInstall 会在配置安装时通过向导方式引导用户输入

> 其它`ServerType`可通过`hydra.Conf.API`,`hydra.Conf.CRON`,`hydra.Conf.MQC`等进行参数配置

####  2. 公共参数配置
指同一个平台下所有`ServerType`可共同使用的配置, 通过```hydra.Conf.Vars()....```指定


DB配置：

```go
    hydra.Conf.Vars().DB().MySQLByConnStr("db", "hydra:123456@tcp(192.168.0.36:10036)/hydra?charset=utf8")
```
缓存配置：

```go
    hydra.Conf.Vars().Cache().GoCache("cache")
```

用户自定义配置:

```go
    hydra.Conf.Vars().Custom("app", "conf", &model.LoginConf{
        LoginHost: hydra.ByInstall,
        APIHost:   hydra.ByInstall,
        UserLoginFailLimit: 5,
        UserLockTime:       24 * 60 * 60,
    })
```



### 文件构建



### 手动配置


