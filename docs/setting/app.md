
### 1. 服务器参数

指平台、系统、集群等参数的设置，有两种设置方式

* 代码设置，app初始化时通过`hydra.With...`设置
  
```go
	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithSystemName("apiserver"),
            hydra.WithServerTypes(http.API),
            hydra.WithClusterName("prod"),			
    )
    app.Start()
```

* 通过cli命令行设置:
```bash
$ ./flowserver run --plat test --system apiserver
```

* 可设置参数

|  参数名  | 代码设置              | cli设置 | 说明                                               |
| :------: | --------------------- | :-----: | -------------------------------------------------- |
|  平台名  | hydra.WithPlatName    |  支持   | 平台或大系统名                                     |
|  系统名  | hydra.WithSystemName  |  支持   | 系统或server名                                     |
| 集群名称 | hydra.WithClusterName |  支持   | 区分普通集群与灰度集群                             |
| 注册中心 | hydra.WithRegistry    |  支持   | 注册中心地址，用于获取配置与服务注册，单机版不指定 |
| 服务类型 | hydra.WithServerTypes | 不支持  | 当前应用支持的服务器类型                           |
| 调试模式 | hydra.WithDebug       |  支持   | 打开详细日志输出                                   |
|  版本号  | hydra.WithVersion     | 不支持  | 当前应用的版本号                                   |
|   用途   | hydra.WithUsage       | 不支持  | 当前应用的说明信息                                 |
|  IP掩码  | hydra.WithIPMask      | 不支持  | 主机有多个IP时使用哪个IP进行服务注册               |
| 性能分析 | hydra.WithTrace       |  支持   | 启动pprof性能跟踪                                  |