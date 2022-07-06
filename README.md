[hydra](https://github.com/micro-plat/hydra)微服务容器
======================
基于golang实现。

hydra 提供简单的、统一的、易扩展的服务容器框架。通过少量的代码集成，即可实现的丰富功能，如：集群管理、配置管理、服务注册与发现、系统监控、日志归集、链路追踪、安全访问、常用组件等。

目前已应用于20+生产项目，主要功能：

- #### ✓  六类服务
   支持*API*、*WEB*、*WEBSOCKET*、*RPC*、*定时任务*、*消息消费* 等服务器，可在单个应用中组合使用。

- #### ✓ 跨平台
    支持windows, mac, linux以服务方式安装、运行、停止、卸载等。

- #### ✓ 多种部署
  支持分布式集群部署、单机伪集群部署、单机部署。

- #### ✓ 多种集群

    提供对等、分片、主从等集群模式。


- #### ✓ 配置管理
    采用配置中心，中心化管理配置，本地零配置。支持zookeeper, redis, 本地文件，进程内管理配置。

 - #### ✓ 热更新
    配置变更后自动生效，无须手动重启服务

- #### ✓ 注册与发现
    支持zookeeper, etcd, redis等作为注册中心，为远程调用提供服务管理。

- #### ✓ 业务监控

    支持将metric信息(如：QPS、处理时长、响应等)定时上报到influxdb，用于系统运行状况监控大屏显示。

- #### ✓ 链路跟踪

    支持Skywalking，Cat等APM工具进行分布式追踪、性能指标分析、应用和服务依赖分析。
- #### ✓ 日志归集

    支持将本地日志提交到日志归集服务器( [rlog](https://github.com/micro-plat/rlog))，用于日志集中查询分析。

   

- #### ✓ 访问控制
  
    支持白名单、黑名单访问控制，Basic Auth、API KEY、Cookie, JWT等安全验证。
    支持远程认证服务, 提供加解密、验证签等服务。

- #### ✓ 服务器限流
    支持服务器限流、降级、熔断等控制。

- #### ✓ 灰度发布

    支持根据业务规则编写灰度脚本，将用户请求转发到不同集群。

- #### ✓ 提供常用组件库

    redis,memcached,数据库，mqtt,activeMQ,rpc,uuid,分布式锁，http client,rpc client等。





[hydra](https://github.com/micro-plat/hydra)

[hello world示例](/01guide/01helloworld.md)


[构建六类服务器](/01guide/02servers.md)


[服务运行](/02component/01service.md)


[配置中心](/02component/02conf.md)

## 一、示例

- ### 1. 构建API服务
```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {

    //创建app
	app := hydra.NewApp(
            hydra.WithPlatName("test"), //平台名
            hydra.WithSystemName("apiserver"), //系统或应用名
            hydra.WithServerTypes(http.API),
    )

    //注册服务
    app.API("/hello", func(ctx hydra.IContext) interface{} {
        return "hello world"
    })

    //启动app
    app.Start()
}
```
- ### 2. 构建RPC服务

```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/rpc"
)

func main() {

    //创建app
	app := hydra.NewApp(
            hydra.WithServerTypes(rpc.RPC),
    )

    //注册服务
    app.RPC("/hello", func(ctx hydra.IContext) interface{} {
        return "hello world"
    })

    //启动app
    app.Start()
}
```
- ### 3. 构建定时任务服务

```go
package main

import (
    "github.com/micro-plat/hydra"   
    "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(cron.CRON),
    )
   
   //注册服务
    app.CRON("/hello",hello,"@every 5s")   
    app.Start()
}
func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```
- ### 4. 构建消息消费服务

```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/mqc"
    "github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(mqc.MQC),
    )

    //注册服务
    app.MQC("/hello",hello,"queue-name")

    //设置消息队列服务器(本地内存MQ，支持redis,mqtt等)  
    hydra.Conf.MQC(lmq.MQ)

    app.Start()
}

func hello(ctx hydra.IContext) interface{} {
        return "success"
}
```


## 二、组合服务

```go
package main

import (
    "github.com/micro-plat/hydra"   
     "github.com/micro-plat/hydra/hydra/servers/http"
     "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {

	app := hydra.NewApp(
            hydra.WithServerTypes(http.API,cron.CRON),
    )

    app.API("/hello",hello)
    app.CRON("/hello",hello,"@every 5s") 
    app.Start()
}
func hello(ctx hydra.IContext) interface{} {
        return "hello world"
}
```


```sh
$ ./flowserver run --plat test
```

日志如下：
```sh
[2020/07/08 09:36:31.140432][i][29f63e41d]初始化: /test/flowserver/api-cron/1.0.0/conf
[2020/07/08 09:36:31.143027][i][29f63e41d]启动[api]服务...
[2020/07/08 09:36:31.643524][i][b65655312]启动成功(api,http://192.168.4.121:8080,1)
[2020/07/08 09:36:31.643885][i][29f63e41d]启动[cron]服务...
[2020/07/08 09:36:31.844844][i][3908a5ccc]启动成功(cron,cron://192.168.4.121,1)
[2020/07/08 09:36:32.346047][d][3908a5ccc]当前server启动为: master
[2020/07/08 09:36:36.648149][i][01751ece6]cron.request: GET /hello from 192.168.4.121
[2020/07/08 09:36:36.648244][i][01751ece6]cron.response: GET /hello 200  193.356µs
[2020/07/08 09:36:41.651858][i][00f45e17b]cron.request: GET /hello from 192.168.4.121
[2020/07/08 09:36:41.651911][i][00f45e17b]cron.response: GET /hello 200  159.694µs
```

服务已启动，cron server开始周期(每隔5秒)执行任务


```sh
$ curl http://192.168.4.121:8080/hello
```
返回内容：
```sh
hello world
```


## 三、 服务器配置


#### 1. 设置平台、系统、集群名称

```go
	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithSystemName("apiserver"),
            hydra.WithServerTypes(http.API),
            hydra.WithClusterName("prod"),			
    )
    app.Start()

```

#### 2. 运行时 通过cli参数指定平台、系统、集群名称

```go
	app := hydra.NewApp(
            hydra.WithServerTypes(http.API),	
    )
    app.Start()

```

#### 3. *API*监听端口：
```go
	hydra.Conf.API("8090")	
```

#### 4. *API*超时时长
```go
	hydra.Conf.API("8090", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)) 
```

#### 5. 白名单、黑名单
```go
    hydra.Conf.API("8080").			
	WhiteList(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
	BlackList(blacklist.WithIP("192.168.4.120"))
```
#### 6. JWT认证
- 1. 配置*/member*开头人路径不验证，将jwt串串存到header中：
```go
    hydra.Conf.API("8080").Jwt(jwt.WithExcludes("/member/**"), jwt.WithHeader())
```
- 2. 登录成功设置jwt信息:

```go
type UserInfo struct{
    Name string `json:"name"`
    UID string `json:"id"`
}


func login(ctx hydra.IContext) interface{} {
    userInfo := UserInfo{UID:"209867923",Name:"colin"}
	ctx.User().Auth().Response(&userInfo)
	return "success"
}

```

- 3. 从jwt中获取用户信息:
```go
func login(ctx hydra.IContext) interface{} {
     userInfo :=new(UserInfo)
	ctx.User().Auth().Bind(&userInfo)
	return "success"
}
```
#### 7. 灰度发布

以 *202.222.* 开头的IP，转发到名称中包含*gray*的集群（tengo脚本）:

```go
	hydra.Conf.API("8090").Proxy(`	
                request := import("request")
                app := import("app")
                text := import("text")
                types :=import("types")
                fmt := import("fmt")

                getUpCluster := func(){
                    ip := request.getClientIP()
                    current:= app.getCurrentClusterName()
                    if text.has_prefix(ip,"202.222."){
                        return app.getClusterNameBy("gray")
                    }
                    return current
                }
                upcluster := getUpCluster()
		`)
```
#### 8. 自定义输出格式

- 1. 处理输出（tengo脚本）：
```go
    hydra.Conf.API("8070").Render(`
            request := import("request")
            response := import("response")
            text := import("text")
            types :=import("types")

            rc:="<response><code>{@status}</code><msg>{@content}</msg></response>"

            getContent := func(){  
                input:={status:response.getStatus(),content:response.getRaw()["id"]}

                if text.has_prefix(request.getPath(),"/tx/request"){
                    return [200,types.translate(rc,input)]
                }
                if text.has_prefix(request.getPath(),"/tx/query"){
                    return [200,types.translate(rc,input),"text/plain"]
                }
                return [200,response.getContent()]
            }

            render := getContent()
            `)
```

- 2. 服务实现
```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)

	app.API("/tx/request", request)
	app.API("/tx/query", request)
	app.API("/request", request)
	app.Start()
}
func request(ctx hydra.IContext) interface{} {
	return map[string]interface{}{
		"id": 101010,
	}
}
```

- 3. 请求服务
```sh
~/work/bin$ curl http://localhost:8070/tx/request
<response><code>200</code><msg>101010</msg></response>

~/work/bin$ curl http://localhost:8070/request
{"id":101010}
```

#### 9. 数据库配置
```go
 hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra"))
```

使用方法:

```go
 func GetUserInfo(id int64)(map[string]interface{},error)
    db:=hydra.C.GetDB()
    userInfo,err:=db.Query("select t.* from user_info t where t.id=@id",map[string]interface{}{"id",id})
    return userInfo,err
```

#### 10. 缓存配置
```go
 hydra.Conf.Vars().Cache("cache", redis.New("192.168.0.109"))
```

使用方法:

```go
 func GetUserInfo(id int64)(string,error)
    cache:=hydra.C.GetCache()
    userInfo,err:=cache.Get(fmt.Sprintf("user:info:%d",id))
    return userInfo,err
```

## 四、 服务注册

- 1. 服务函数
```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
	app.API("/request", request)
	app.Start()
}
func request(ctx hydra.IContext) interface{} {
	return "success"
}
```

- 2. 包含Handle服务的struct构造函数或对象

```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
    app.API("/order", NewOrderService)
    //或
    // app.API("/order", &OrderService{})
	app.Start()
}

type OrderService struct {
}
func NewOrderService()*OrderService{
    return &OrderService{}
}

func (o *OrderService) Handle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------订单处理----------"))
	return "success"
}
```

- 3. 多个服务注册

```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
	app.API("/order/*", NewOrderService)
	app.Start()
}

type OrderService struct {
}
func NewOrderService()*OrderService{
    return &OrderService{}
}

func (o *OrderService) RequestHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------收单处理----------"))
	return "success"
}
func (o *OrderService) QueryHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------查单处理----------"))
	return "success"
}
```
> 实际注册了两个服务 */order/request*,*/order/query*


- 3. RESTful服务注册

```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
	app.API("/product", &ProductService{})
	app.Start()
}

type ProductService struct {
}



func (o *ProductService) GetHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------查询产品----------"))
	return "success"
}
func (o *ProductService) PostHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------新增产品----------"))
	return "success"
}
func (o *OrderService) PutHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------修改产品----------"))
	return "success"
}
func (o *OrderService) DeleteHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------禁用产品----------"))
	return "success"
}
```

> 实际注册了一个服务 */product*



- 4. 钩子函数

```go
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
    app.API("/order", NewOrderService)
	app.Start()
}

type OrderService struct {
}
func NewOrderService()*OrderService{
    return &OrderService{}
}
func (o *OrderService) Handling(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------预处理----------"))
	return "success"
}
func (o *OrderService) Handle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------订单处理----------"))
	return "success"
}
func (o *OrderService) Handled(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------后处理----------"))
	return "success"
}

```
## 五、 服务响应


- 1. 普通字符串
```go
func request(ctx hydra.IContext) interface{} {
	return "success"
}
```

> 响应内容: *success*  
> Content-Type: text/plain; charset=utf-8



- 2. JSON字符串

```go
func request(ctx hydra.IContext) interface{} {
	return `{"name":"colin"}`
}
```

> 响应内容: *{"name":"colin"}*  
> Content-Type: application/json; charset=utf-8


- 3. XML字符串

```go
func request(ctx hydra.IContext) interface{} {
	return `<name>colin</name>`
}
```

> 响应内容: &#60;name&#62; colin&#60;/name&#62; 
> Content-Type: application/xml; charset=utf-8


- 4. Map
```go
func request(ctx hydra.IContext) interface{} {
	return map[string]interface{}{"name":"colin"}
}
```
> 响应内容: *{"name":"colin"}*  
> Content-Type: application/json; charset=utf-8

- 5. Struct

```go
func request(ctx hydra.IContext) interface{} {
	return &UserInfo{Name:"colin"}
}
```

> 响应内容: *{"name":"colin"}*  
> Content-Type: application/json; charset=utf-8


- 6. Map,指定ContentType

```go
func request(ctx hydra.IContext) interface{} {
    ctx.Response().ContentType(context.UTF8XML)
	return map[string]interface{}{"name":"colin"}
}
```

> 响应内容:  *&#60;name/&#62;colin &#60;name/&#62;*
> Content-Type: application/xml; charset=utf-8


- 7. Struct,指定ContentType

```go
func request(ctx hydra.IContext) interface{} {
    ctx.Response().ContentType(context.UTF8XML)
	return &UserInfo{Name:"colin"}
}
```

> 响应内容: *&#60;name/&#62;colin &#60;name/&#62;*
> Content-Type: application/xml; charset=utf-8

- 8. 普通错误

```go
func request(ctx hydra.IContext) interface{} {
   return fmt.Errorf("err")
}
```

> 响应内容: err
> Content-Type: text/plain; charset=utf-8

- 9. 带状态码的错误

```go
func request(ctx hydra.IContext) interface{} {
   return errs.New(500,"err")
}
```

> 状态码: 500
> 响应内容: err
> Content-Type: text/plain; charset=utf-8


