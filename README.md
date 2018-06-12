# hydra
hydra是一个能够快速开发http接口, web应用，RPC服务，流程服务，任务调度，消息消费(MQ Consumer)的服务框架，你只需要关注你所提供的服务自身，简单的配置,即可对外提供服务。


  hydra特点
* 多服务集成: 支持 `http接口`,`rpc服务器`,`web站点`,`定时任务`,`MQ消费流程`,`websocket` 等服务
* 开发简单: 10行代码即可实现一个服务，并可以以6种服务(`api`,`rpc`,`web`,`cron`,`mqc`,`ws`)运行，对外提供服务
* 部署简单: 编译后只有一个可执行程序，复制到目标服务器，通过自带的命令行参数启动即可 
* 本地零配置: 本地无需任何配置，启动时从注册中心拉取
* 配置自动更新: 配置发生变化后可自动通知到服务器，配置变更后自动更新到服务器，必要时自动重启服务器
* 服务注册与发现: 基于 zookeeper 和本地文件系统(用于单机测试)的服务注册与发现
* 实时监控: 服务器运行状况如:QPS,服务执行时长，执行结果，CPU，内存等信息自动上报到influxdb,通过grafana配置后即可实时查看服务状况
* 统一日志: 配置远程日志接收接口，则本地日志每隔一定时间压缩后发送到远程日志服务([远程日志](https://github.com/micro-plat/logsaver))
* 自动更新: 通过服务器远程控制接口，远程发送更新信息，服务器自动检查版本，并下载，更新服务






## hydra架构图

![架构图](https://github.com/micro-plat/hydra/blob/master/quickstart/hydra.png?raw=true)


## hydra启动过程


![架构图](https://github.com/micro-plat/hydra/blob/master/quickstart/flow.png?raw=true)

## 文档目录
1. [快速入门](README.md#hydra)
      * [hydra安装](https://github.com/micro-plat/hydra/blob/master/quickstart/2_install.md)
      * [gaea工具简介](https://github.com/micro-plat/hydra/blob/master/quickstart/3.install_gaea.md)
       * [创建第一个项目](https://github.com/micro-plat/hydra/blob/master/quickstart/6.first_project.md)
      
2. [服务器管理](https://github.com/micro-plat/hydra/blob/master/quickstart/7.server.intro.md)
      * [接口服务器](https://github.com/micro-plat/hydra/blob/master/quickstart/api/1.api_intro.md)
          + [路由配置](https://github.com/micro-plat/hydra/blob/master/quickstart/api/2.api_router.md)         
          + [静态文件](https://github.com/micro-plat/hydra/blob/master/quickstart/api/3.api_static.md)
          + [metric](https://github.com/micro-plat/hydra/blob/master/quickstart/api/4.api_metric.md)
          + [jwt选项](https://github.com/micro-plat/hydra/blob/master/quickstart/api/5.api_auth.md)
          + [附加header](https://github.com/micro-plat/hydra/blob/master/quickstart/api/6.api_header.md)
          + [熔断降级](https://github.com/micro-plat/hydra/blob/master/quickstart/api/7.api_circuit.md)
      * web服务器
         + [路由配置](https://github.com/micro-plat/hydra/blob/master/quickstart/api/2.api_router.md)  
         + view配置       
          + [静态文件](https://github.com/micro-plat/hydra/blob/master/quickstart/api/3.api_static.md)
          + [metric](https://github.com/micro-plat/hydra/blob/master/quickstart/api/4.api_metric.md)
          + [jwt选项](https://github.com/micro-plat/hydra/blob/master/quickstart/api/5.api_auth.md)
          + [附加header](https://github.com/micro-plat/hydra/blob/master/quickstart/api/6.api_header.md)
      * rpc服务器
          + 路由配置
          + [jwt选项](https://github.com/micro-plat/hydra/blob/master/quickstart/api/5.api_auth.md)
          + 限流配置
          + [metric](https://github.com/micro-plat/hydra/blob/master/quickstart/api/4.api_metric.md)
      * mq consumer
          + 队列配置
          + [metric](https://github.com/micro-plat/hydra/blob/master/quickstart/api/4.api_metric.md)
      * 定时服务
          + 任务配置
          + [metric](https://github.com/micro-plat/hydra/blob/master/quickstart/api/4.api_metric.md)
3. 日志组件
4. Context
      * 输入参数
      * 缓存操作
      * 数据库处理
      * RPC请求
5. Response
6. [监控与报警](https://github.com/micro-plat/hydra/blob/master/quickstart/alarm/1.alarm.md)