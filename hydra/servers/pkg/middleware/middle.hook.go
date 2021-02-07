package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

//Handlers 中间件处理函数
type Handlers []Handler

//ICustomMiddleware 用户自定义组件
type ICustomMiddleware interface {
	Add(handler ...Handler)
	Count() int
}

//DispFunc 返回DispFunc
func (c Handlers) DispFunc() []dispatcher.HandlerFunc {
	list := make([]dispatcher.HandlerFunc, 0, len(c))
	for _, item := range c {
		list = append(list, item.DispFunc())
	}
	return list
}

//GinFunc 返回GinFunc
func (c Handlers) GinFunc() []gin.HandlerFunc {
	list := make([]gin.HandlerFunc, 0, len(c))
	for _, item := range c {
		list = append(list, item.GinFunc())
	}
	return list
}

//Add 添加组件
func (c Handlers) Add(handler ...Handler) {
	c = append(c, handler...)
}

//Count 统计组件数量
func (c Handlers) Count() int {
	return len(c)
}
