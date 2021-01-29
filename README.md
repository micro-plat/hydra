[hydra](https://github.com/micro-plat/hydra)
======================

后端一站式服务框架,提供 *API*、*WEB*、*WEBSOCKET*、*RPC*、*定时任务*、*消息消费* 等服务，可任意组合以单体方式部署，并提供丰富的服务治理功能：




- #### 1. 配置中心、注册中心

使用zookeeper, etcd, redis, 本地文件系统，内存等作为配置与注册中心, 管理服务器配置、 服务注册与发现、 服务监控等。所有配置支持热更新。


- #### 2. 服务化(SOA)与注册发现

一切皆服务，使用统一的request, response编写业务代码，注册到任意服务器即可。
系统自动将服务发布到注册中心，调用时根据规则自动进行负载均衡。


- #### 3. 链路跟踪与业务监控

支持Skywalking，Cat等APM工具进行分布式追踪、性能指标分析、应用和服务依赖分析。
采集服务QPS、处理时长、响应等统计信息存入influxdb，使用grafana配置图表即可多维度监控系统状态。


- #### 4. 分布式日志归集

使用[rlog](https://github.com/micro-plat/rlog)将分布式服务器日志提交到日志归集服务器，并存储到Elastic search中，通过[rlog](https://github.com/micro-plat/rlog)查询界面可查询服务的全链路日志。


- #### 5. 访问控制与安全认证

提供白名单、黑名单访问控制，服务器限流、降级、熔断等控制。
Basic Auth、API KEY、Cookie, JWT等安全认证。
[SAS]远程服务认证，API KEY, RSA, 动态秘钥，DES，AES等签名与加解密服务。

- #### 6. 灰度发布与集群转发

可根据请求路径、用户信息、请求头、Cookie等编写tengo脚本，灵活编写灰度或集群转发规则。


- #### 7. 流程任务多集群模式

提供对等、分片、主从等集群模式。


- #### 8. 基础组件

提供缓存、数据库、消息队列、RPC、分布式锁等组件工具。

- #### 9. 跨平台支持
支持Windows、Linux和macOS跨平台运行，以系统服务方式安装、卸载、运行、停止、自动启动管理等。


----------


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
[2020/07/08 09:36:32.346047][d][3908a5ccc]the cron server is started as master
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


