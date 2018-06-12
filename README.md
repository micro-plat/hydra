# hydra 
hydra 是基于go语言和众多开源项目，实现的分布式服务框架。其核心设计目标是快速开发、使用简单、功能强大、轻量级、易扩展。敏捷开发，快速交付，致力于打造企业级后端服务的快速开发方案。使用hydra可开发`http接口`,`rpc服务器`,`web站点`,`定时任务`,`MQ消费流程`,`websocket` 等服务

  hydra特点
* 简单轻量: 几行代码即可集成到自己的应用
* 多服务集成: 支持 `http接口`,`rpc服务器`,`web站点`,`定时任务`,`MQ消费流程`,`websocket` 等服务
* 优秀性能: 底层使用`gin`,`时间轮`等框架和算法，借助`golang`自身优势，提供优秀的服务性能 
* 部署简单: 编译后只有一个可执行程序，复制到目标服务器，通过自带的命令行参数启动即可 
* 本地零配置: 本地无需任何配置，启动时从注册中心拉取
* 配置热更新: 配置发生变化后可自动通知到服务器，配置变更后自动更新到服务器，必要时自动重启服务器
* 服务注册发现: 基于 zookeeper 和本地文件系统(用于单机测试)的服务注册与发现
* 丰富功能: `智能路由`,`静态文件`,`安全认证`,`熔断降级`等
* 性能调优： 集成`pprof`等提供可实时监控的性能调优工具
* 实时监控: 服务器运行状况如:QPS,服务执行时长，执行结果，CPU，内存等信息自动上报到influxdb,通过grafana配置后即可实时查看服务状况
* 统一日志: 配置远程日志接收接口，则本地日志每隔一定时间压缩后发送到远程日志服务([远程日志](https://github.com/micro-plat/logsaver))
* 自动更新: 通过服务器远程控制接口，远程发送更新信息，服务器自动检查版本，并下载，更新服务


##  安装hydra
```sh
go get github.com/micro-plat/hydra
```

## 示例项目

示例代码请参考[examples](https://github.com/micro-plat/hydra/tree/master/examples)


### 1. hello world 示例代码
* 1. 编写代码

新建项目`hello`,并添加`main.go`文件输入以下代码

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
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}

```

* 2. 编译项目
```sh
go install hello
```
* 3. 启动服务
```sh
hello start
```
* 4. 测试服务
```sh
  curl http://localhost:8090/hello
 hello world
```


### 2. 通过命令行指定运行参数

```go
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {

//以下hydra.With开头的参数都可以使用命令行传入
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		//hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}

```
```
* 启动服务时指定注册中心
```sh
hello start --registry "fs://../"
```

### 3. 使用对象注册服务
```sh
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", newHelloService) //使用对象的构造函数注册服务
	app.Start()
}

type helloService struct {
	container component.IContainer
}

func newHelloService(container component.IContainer) (u *helloService) {
	return &helloService{container: container}
}
func (u *helloService) Handle(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```



### 4. RESTful服务
```sh
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", newHelloService) //使用对象的构造函数注册服务
	app.Start()
}

type helloService struct {
	container component.IContainer
}

func newHelloService(container component.IContainer) (u *helloService) {
	return &helloService{container: container}
}
func (u *helloService) GetHandle(ctx *context.Context) (r interface{}) {
	return "get ->  hello world"
}
func (u *helloService) PostHandle(ctx *context.Context) (r interface{}) {
	return "post -> hello world"
}
func (u *helloService) PutHandle(ctx *context.Context) (r interface{}) {
	return "put -> hello world"
}
func (u *helloService) DeleteHandle(ctx *context.Context) (r interface{}) {
	return "delete -> hello world"
}
```




### 5. 输入参数绑定到对象
`ctx.Request.Bind`绑定输入参数到指定的对象
`valid:"required"`参数用于设置参数验证方式，参考[govalidator](https://github.com/asaskevich/govalidator)
```sh
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

//NewsInfo 新的新闻消息
type NewsInfo struct {
	ID       string `json:"id"` //程序中指定
	Title    string `form:"title" json:"title" valid:"required"`
	Brief    string `form:"brief" json:"brief" valid:"required"`
	Image    string `form:"image" json:"image" valid:"required"`
	Keywords string `form:"keywords" json:"keywords" valid:"required"`
	Content  string `form:"content" json:"content" valid:"required"`
}

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", newHelloService) //使用对象的构造函数注册服务
	app.Start()
}

type helloService struct {
	container component.IContainer
}

func newHelloService(container component.IContainer) (u *helloService) {
	return &helloService{container: container}
}
func (u *helloService) PostHandle(ctx *context.Context) (r interface{}) {
	var input NewsInfo
	if err := ctx.Request.Bind(&input); err != nil { //输入参数绑定到对象
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}
	return input
}
```



### 5. 完整示例


`ctx.Request.Check`检查输入参数

`ctx.Request.Get...`获取输入参数

`ctx.Log...`打印日志

`返回值`可以是`error`,`context.IError`,`map`,`struct`
```sh
package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"), //平台名称
		hydra.WithSystemName("demo"), //系统名称
		hydra.WithClusterName("test"), //集群名称
		hydra.WithServerTypes("api"), //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心	
		hydra.WithDebug())

	app.Micro("/hello", newHelloService) //使用对象的构造函数注册服务
	app.Start()
}

type helloService struct {
	container component.IContainer
}

func newHelloService(container component.IContainer) (u *helloService) {
	return &helloService{container: container}
}
func (u *helloService) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------处理用户请求----------------")
	ctx.Log.Info("1、获取参数")
	err := ctx.Request.Check("user_id", "media_id")
	if err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err) //返回错误
	}
	uid := ctx.Request.GetInt("user_id", 0)
	mid := ctx.Request.GetInt("media_id", 0)
	return map[string]interface{}{
		"user_id":uid,
		"media_id":mid,
	} //返回map,默然以json方式输出
}
```