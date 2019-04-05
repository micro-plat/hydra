## 服务配置与安装

hydra服务所有配置都采用注册中心管理, 本地零配置. 

  *  目前已实现的注册中心有`zookeeper`,`file system`.
一般 `zookeeper`用于生产环境,`file system`用于本地测试 
*   如需使用其它注册中心(`etcd`,`consul`),需自行实现`registry.IRegistry`[接口](https://github.com/micro-plat/hydra/tree/master/registry)


参数配置方式有两种:
 

  * 手工编写[不推荐]
     > 使用三方工具连接到注册中心,创建节点, 添加配置数据. 此方式需了解配置节点结构, 操作较繁琐,不推荐


  *  代码 + 命令行[推荐]
  
     > 1. 固定不变的参数,写在代码中. 如: `平台名称`,`系统名称`,`服务器类型`,`开发环境的数据库连接`,`安全认证方式`,`跨域配置`,`静态服务配置`. 通过`install`命令安装到注册中心

     > 2. 启动时才能确定的参数,命令行指定,如:`服务集群名称`,`服务端口号`,`生产环境数据库配置`

     > 3. 启动必须的5个参数未在代码中指定时会自动设置为命令行参数,启动时必须输入

### 一. 注册中心参数配置与安装

#### 1. 必须参数配置

启动必须的参数只有5个,分别是: `平台名称`,`系统名称`,`服务器类型`,`集群名称`,`注册中心地址`. 代码中未指定时,自动设置为命令行参数, 执行`install`, `run`时引导用户输入

示例代码如下:


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
		hydra.WithSystemName("helloserver"), //系统名称
		hydra.WithServerTypes("api"), //服务器类型为http api文件系统作为注册中心
		hydra.WithDebug())

	app.API("/hello",hello)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```
  当前示例未在代码中指定`集群名称`,`注册中心地址`,在执行`install`,`run`命令时要求用户必须输入:

```sh
$ ./helloserver install 
  服务器运行缺少参数，请查看以下帮助信息

NAME:
   helloserver install - 安装注册中心配置和本地服务

USAGE:
   helloserver install [command options] [arguments...]

OPTIONS:
     ...
```
根据提示输入注册中心地址(`-r`)与集群地址(`-c`), 即可启动,进入安装向导

```sh
$ ./helloserver install -r fs://../ -c test
-> 创建注册中心配置数据? 如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):
```
输入`2`完成安装
```sh
$ sudo ./helloserver install -r fs://../ -c test
Password:
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		创建配置: /myplat_debug/helloserver/api/test/conf
Install helloserver:					[  OK  ]
```
   > 注册中心平台名称变为了`myplat_debug`. 因为我们通过`hydra.WithDebug()`指定为`debug`模式,该模式下仅对系统产生两点影响 1. 平台名称中增加`_dubug` 2. 系统发生错误时会将错误内容返回到请求调用方,便于查找原因. 正式环境中可修改为`false` `app.IsDebug=false`,或不指定`hydra.WithDebug()`

安装分为两步, 1. 创建注册中心配置    2. 安装本地服务

#### 2. 服务器配置参数
服务器配置参数可通过代码指定,不确定的参数可通过变量符`#`指定,安装时会引导用户输入.也可以直接通过三方工具连接到注册中心添加或修改

 #### 2.1 服务器启动端口

在`app.Start()`前添加代码
```go
app.Conf.API.SetMainConf(`{"address":":9999"}`)
```
完整代码如下: 
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
		hydra.WithSystemName("helloserver"), //系统名称
		hydra.WithServerTypes("api"), //服务器类型为http api文件系统作为注册中心
		hydra.WithDebug())
        app.Conf.API.SetMainConf(`{"address":":9999"}`)

	app.API("/hello",hello)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
```
以上代码在`install`时将`{"address":":9999"}`保存到` /myplat_debug/helloserver/api/test/conf`,启动时将自动拉取该配置, 并将服务器启动端口设置为`9999`

#### 2.2 数据库连接串配置

在`app.Start()`前添加代码
```go
app.Conf.Plat.SetVarConf("db", "db", `{			
				"provider":"mysql",
				"connString":"wechat:12345678@tcp(192.168.0.36)/wechat",
				"maxOpen":10,
				"maxIdle":1,
				"lifeTime":300		
		}`)
```
以上代码在执行`install`后, 自动向注册中心创建节点`/myplat_debug/var/db/db`并保存内容.

#### 2.3 生产环境配置

如果是生产环境,则连接字符串无法写在代码中,则用`#`开头的变量名作为占位符,在执行`install`命令时会由安装人员输入,代码如下: 
```go
	app.Conf.Plat.SetVarConf("db", "db", `{			
		"provider":"mysql",
		"connString":"#connstring",
		"maxOpen":10,
		"maxIdle":1,
		"lifeTime":300		
}`)
```
编译并执行`install`命令

```sh
$ go build
$ sudo ./helloserver install -r fs://../ -c test
Password:
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		修改配置: /myplat/helloserver/api/test/conf
		* 请输入connstring(/myplat/var/db/db等配置中使用):
```


#### 2.4 手动编写配置
除系统启动的5个参数外,其它参数不想写在代码里,可通过三方工具连接到注册中心进行配置.zookeeper可使用`ZooInspector`


### 二. 本地服务安装
执行`install`命令执行的第二步为安装本地服务, 安装成功后可直接使用`start`启动,`stop`停止,`status`查看服务是否运行,`remove`卸载服务.


下一节[服务注册与启动]((https://github.com/micro-plat/hydra/tree/master/docs/service.types.register.md))