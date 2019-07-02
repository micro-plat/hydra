# 构建远程调用服务(RPC)

本示例介绍如何构建RPC服务器，以及RPC调用方法。


RPC服务器的构建方法与API服务器的构建相同，只需修改服务器类型为`rpc`即可，其它完全相同。


#### 1. 创建服务器


`main.go`
```go
package main

import "github.com/micro-plat/hydra/hydra"

type rpcserver struct {
	*hydra.MicroApp
}

func main() {
	app := &rpcserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("rpcserver"),
			hydra.WithServerTypes("rpc")),
	}
	app.init()
	app.Start()
}

```
> app.init 用于挂载服务配置，注册等处理


`config.dev.go`
```go
// +build !prod

package main

func (rpc *rpcserver) config() {
	rpc.IsDebug = true
	rpc.Conf.RPC.SetMainConf(`{"address":":9090","trace":true}`)
	rpc.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"mysql",
			"connString":"mrss:123456@tcp(192.168.0.36)/mrss?charset=utf8",
			"maxOpen":20,
			"maxIdle":10,
			"lifeTime":600		
	}`)
	
}
```

#### 2. 初始化检查与服务注册



```go
package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/rpcserver01/services/order"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (rpc *rpcserver) init() {
	rpc.config()
	rpc.handling()

	rpc.Initializing(func(c component.IContainer) error {
		//检查db配置是否正确
		if _, err := c.GetDB(); err != nil {
			return err
		}

		return nil
	})

	//服务注册
	rpc.Micro("/order", order.NewOrderHandler)
}
```
> `rpc.Micro`注册为`api`,`rpc`类型，也可使用`rpc.RPC("/order", order.NewOrderHandler)`只注册为`rpc`服务


#### 3. 请求预处理，验证签名

```go
package main

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/quickstart/demo/rpcserver11/modules/merchant"
)

func (rpc *rpcserver) handling() {
	rpc.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {
		if err := ctx.Request.Check("merchant_id"); err != nil {
			return err
		}
		key, err := merchant.GetKey(ctx,ctx.Request.GetInt(merchant_id))
		if err != nil {
			return err
		}
		if !ctx.Request.CheckSign(key) {
			return fmt.Errorf(908, "商户签名错误")
		}
		return nil
	})
}
```


#### 3. 构建服务


`servers/order.go`

```go
package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type OrderHandler struct {
	container component.IContainer
}

func NewOrderHandler(container component.IContainer) (u *OrderHandler) {
	return &OrderHandler{
		container: container,
	}
}
//QueryHandle 充值结果查询
func (u *OrderHandler) QueryHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------充值结果查询---------------")	
	return "SUCCESS"
}
```
说明:
* 设置Response输出类型为`json`,`xml`,`plain`等，返回给调用方的数据不会转换为相应的格式，只会将此信息保存到返回参数的`HeaderMap`中。可直接获取`HeaderMap`的值设置到`api`服务的`Response`中输出到`http`响应流


#### 4. 安装并启动RPC服务器

安装配置信息:


```sh
~/work/bin$ rpcserver01 registry -r zk://192.168.0.109 -c yl
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		创建配置: /mall_debug/rpcserver/rpc/yl/conf
		创建配置: /mall_debug/var/db/db
```

运行服务：

```sh
~/work/bin$ rpcserver01 run -r zk://192.168.0.109 -c yl
[2019/07/02 11:44:12.274275][i][49f56b398]Connected to 192.168.0.109:2181
[2019/07/02 11:44:12.289725][i][49f56b398]Authenticated: id=246395503264334090, timeout=4000
[2019/07/02 11:44:12.289766][i][49f56b398]Re-submitting `0` credentials after reconnect
[2019/07/02 11:44:12.321996][i][49f56b398]初始化 /mall_debug/rpcserver/rpc/yl
[2019/07/02 11:44:12.336878][i][4f24e906a]开始启动[RPC]服务...
[2019/07/02 11:44:12.337520][d][4f24e906a][未启用 jwt设置]
[2019/07/02 11:44:12.337530][d][4f24e906a][未启用 header设置]
[2019/07/02 11:44:12.337535][d][4f24e906a][未启用 metric设置]
[2019/07/02 11:44:12.337538][d][4f24e906a][未启用 host设置]
[2019/07/02 11:44:12.962662][i][4f24e906a]服务启动成功(RPC,tcp://192.168.4.121:9090,1)
```


#### 5. 查看服务注册

使用`ZooInspector`连接到zookeeper

查找当前服务节点`/[platName]/services/rpc/[systeName]/[serviceName]/providers/`


当前服务的节点如下:

mall_debug
-------services
-----------rpc
--------------order
-----------------query
--------------------providers
-----------------------192.168.4.121

> 启动多台服务器则所有服务器都会注册到`/[platName]/services/rpc/[systeName]/[serviceName]/providers/`目录