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


#### 2. RESTfull服务
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