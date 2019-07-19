package middleware

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//Body 处理请求的body参数
func Body() gin.HandlerFunc {
	return func(ctx *gin.Context) {	
		if body, err := ioutil.ReadAll(ctx.Request.Body); err == nil {
			ctx.Set("__body_", body)
		}
		ctx.Next()
	}
}
