## 服务器生命周期

hydra支持6种服务器, 每个hydra应用(hydra.MicroApp)实例可包含多个服务器（同一个进程中）, 不同服务器相互独立.

![流程](https://github.com/micro-plat/hydra/blob/master/docs/imgs/lifetime.png)

#### 1. 自定义命令行参数
   某些参数需要用户在执行安装，启动等命令时输入：
   
 1. 参数绑定：
```go
app.Cli.Append(hydra.ModeRun, cli.StringFlag{
		Name:  "ip,i",
		Usage: "IP地址",
	})
```
2. 执行`run`命令时引导输入：

```sh
argserver run -r zk://192.168.0.177 -n /test/test/api/yl -ip 192.168.5.71
```

3. 参数验证:
```go
app.Cli.Validate(hydra.ModeRun, func(c *cli.Context) error {
		if !c.IsSet("ip") {
			return fmt.Errorf("未设置ip地址")
		}
		return nil
	})
```
4. 服务器初始化时获取参数值:

```go
	app.Initializing(func(component.IContainer) error {
		fmt.Println("ip:", app.Cli.Context().String("ip"))
		return nil
	})
```


#### 2. 配置参数设置

1. 可通过`conf`中各服务器的配置函数，设置每种服务器的配置参数，[参考](https://github.com/micro-plat/hydra/blob/master/docs/service.conf.install.md)
```go
app.Conf.API.SetMainConf
app.Conf.API.SetSubConf

app.Conf.RPC.SetMainConf
app.Conf.RPC.SetSubConf

app.Conf.WEB.SetMainConf
app.Conf.WEB.SetSubConf

app.Conf.WS.SetMainConf
app.Conf.WS.SetSubConf

app.Conf.CRON.SetMainConf
app.Conf.CRON.SetSubConf

app.Conf.MQC.SetMainConf
app.Conf.MQC.SetSubConf
```
2. 平台内公共参数:
```go
app.Conf.Plat.SetVarConf
```

3. 自定义安装，可使用installer函数:
```go
app.Conf.CRON.Installer(func(c component.IContainer) error {

})
```




