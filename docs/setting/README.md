
## 通用配置

所有server都支持的配置

### 1. 平台、系统、集群名称

通过代码指定
```go
	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithSystemName("apiserver"),
            hydra.WithServerTypes(http.API),
            hydra.WithClusterName("prod"),			
    )
    app.Start()

```

运行时指定

```go
	app := hydra.NewApp(
            hydra.WithServerTypes(http.API),	
    )
    app.Start()

```


```sh

$ ./apiserver run --plat test

```

### 2. 缓存
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

### 3. 数据库

```go
 hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra"))

 hydra.Conf.Vars().DB().MySQL("db", "hydra", "123456", "121.196.168.242:23306", "hydra", db.WithConnect(20, 10, 600))
```
使用方法:

```go
 func GetUserInfo(id int64)(map[string]interface{},error)
    db:=hydra.C.GetDB()
    userInfo,err:=db.Query("select t.* from user_info t where t.id=@id",map[string]interface{}{"id",id})
    return userInfo,err
```

## 外部服务

指对外提供服务的API,WEB,WS,RPC支持的配置

### 监听端口
```go
	hydra.Conf.API("8090")	

	hydra.Conf.Web("8080")	

	hydra.Conf.RPC("8070")	
```

### 服务超时时长
```go
	hydra.Conf.API("8090", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)) 
```

### 白名单、黑名单
```go
    hydra.Conf.API("8080").WhiteList(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
	BlackList(blacklist.WithIP("192.168.4.120"))
```
### JWT认证

- 配置*/member*开头不验证，将jwt串存到header中：
```go
    hydra.Conf.API("8080").Jwt(jwt.WithExcludes("/member/**"), jwt.WithHeader())
```
- 登录成功设置jwt信息:

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

- 从jwt转换为struct
```go
func login(ctx hydra.IContext) interface{} {
     userInfo :=new(UserInfo)
	ctx.User().Auth().Bind(&userInfo)
	return "success"
}
```
### 灰度发布

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
### 自定义输出格式

1. 处理输出（tengo脚本）：
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

2. 服务实现
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

  3. 请求服务
  
```sh
~/work/bin$ curl http://localhost:8070/tx/request
<response><code>200</code><msg>101010</msg></response>
~/work/bin$ curl http://localhost:8070/request
{"id":101010}
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


