package swap

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//GinCtx gin.context
type GinCtx struct {
	*gin.Context
}

//GetBody 获取body
func (c *GinCtx) GetBody() (string, bool) {
	if body, err := ioutil.ReadAll(c.Request.Body); err == nil {
		return string(body), true
	}
	return "", false
}
