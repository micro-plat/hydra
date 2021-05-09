# 安装部署

通过`hydra.NewApp`创建的应用，自动包含了对应用管理的若干功能，可通过`--help`查看，以`flowserver`为例:

```sh 
$ ./flowserver --help
NAME:
   flowserver - flowserver(A new hydra application)

USAGE:
   flowserver [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   conf     配置管理, 查看、安装配置信息
   install  安装服务，以服务方式安装到本地系统
   remove   删除服务，从本地服务器移除服务
   run      运行服务,以前台方式运行服务。通过终端输出日志，终端关闭后服务自动退出。
   update   更新应用，将服务发布到远程服务器
   start    启动服务，以后台方式运行服务
   status   查询状态，查询服务器运行、停止状态
   stop     停止服务，停止服务器运行
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     查看帮助信息
   --version, -v  查看版本信息

```

应用的`NAME`,`USAGE`,`VERSION`等都可以通过，`hydra.NewApp(hydra.With....)`指定

