# 构建API服务器二
本示例构建一个供前端网站使用的api服务器。提供4个API接口：产品查询，添加，修改，删除。


| 接口地址       | 功能     |说明|
| -------------- | ------------ | ------------ |
| /member/login | 用户登录 |登录成功后，使用jwt返回登录状态|
| /product   | 添加，修改，查询，删除 |使用RESTful风格实现|

知识点:
* 跨域配置
* jwt配置
* 编写RESTful服务
* 登录状态获取


#### 1. 服务配置

```go

package main

func (api *apiserver) config() {
	api.IsDebug = true
	api.Conf.API.SetSubConf("header", `
				{
					"Access-Control-Allow-Origin": "*",
					"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,PATCH,OPTIONS",
					"Access-Control-Allow-Headers": "X-Requested-With,Content-Type",
					"Access-Control-Allow-Credentials": "true"
				}
            `)
    api.Conf.API.SetSubConf("auth", `
		{
			"jwt": {
				"exclude": ["/member/login"],
				"expireAt": 36000,
				"mode": "HS512",
				"name": "__jwt__",
				"secret": "45d25cb71f3bee254c2bc6fc0dc0caf1"
			}
		}
		`)
}
```

#### 2. 检查用户登录状态
```go
package main

import (
	"fmt"

	"github.com/micro-plat/hydra/component"

	"github.com/micro-plat/hydra/context"
	mem "github.com/micro-plat/sso/flowserver/modules/member"
	xmenu "github.com/micro-plat/sso/flowserver/modules/menu"
)

//bind 检查应用程序配置文件，并根据配置初始化服务
func (api *apiserver) handling() {
	//每个请求执行前执行
	api.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {

		//获取jwt
		jwt, err := ctx.Request.GetJWTConfig() //获取jwt配置
		if err != nil {
			return err
		}
		for _, u := range jwt.Exclude { //排除请求
			if u == ctx.Service {
				return nil
			}
		}

		//缓存用户信息
		var m mem.LoginState
		if err = ctx.Request.GetJWT(&m); err != nil {
			return context.NewError(context.ERR_FORBIDDEN, err)
		}
		if err = mem.Save(ctx, &m); err != nil {
			return err
		}


		//检查用户权限
		tags := r.GetTags(ctx.Service)
		menu := xmenu.Get(ctx.GetContainer().(component.IContainer))
		for _, tag := range tags {
			if tag == "*" {
				return nil
			}
			if err = menu.Verify(m.UserID, m.SystemID, tag, ctx.Request.GetMethod()); err == nil {
				return nil
			}
		}
		return context.NewError(context.ERR_NOT_ACCEPTABLE, fmt.Sprintf("没有权限:%v", tags))
	})
}

```



#### 2. 编写登录接口
```go
package member

import (
	"fmt"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver11/modules/member")

//LoginHandler 用户登录对象
type LoginHandler struct {
	c      component.IContainer
	m      member.IMember
}

//NewLoginHandler 创建登录对象
func NewLoginHandler(container component.IContainer) (u *LoginHandler) {
	return &LoginHandler{
		c:      container,
		m:      member.NewMember(container),
	}
}


//SysHandle 子系统远程登录
func (u *LoginHandler) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("-------用户登录---------")
	//检查输入参数
	if err := ctx.Request.Check("username", "password"); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}
	ctx.Log.Info("2.执行操作")
	//处理用户登录
	member, err := u.m.Login(ctx.Request.GetString("username"),ctx.Request.GetString("password"))
	if err != nil {
		return err
	}
	
	ctx.Log.Info("3.返回数据")
	//设置jwt数据
	ctx.Response.SetJWT(member)
	//记录登录行为
	return member

}


```


#### 3. RESTfull服务
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

//GetHandle 查询
func (u *OrderHandler) GetHandle(ctx *context.Context) (r interface{}) {
 }

 //PostHandle 新增
func (u *OrderHandler) PostHandle(ctx *context.Context) (r interface{}) {
 }

//PutHandle 修改
 func (u *OrderHandler) PutHandle(ctx *context.Context) (r interface{}) {
 }

//DeleteHandle 删除
 func (u *OrderHandler) DeleteHandle(ctx *context.Context) (r interface{}) {
 }
```