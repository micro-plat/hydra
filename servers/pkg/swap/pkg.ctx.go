package swap

import "github.com/micro-plat/hydra/servers/pkg/dispatcher"

//PkgCtx   dispatcher.Context
type PkgCtx struct {
	*dispatcher.Context
}

//GetBody 获取body
func (c *PkgCtx) GetBody() (string, bool) {
	if body, ok := c.Request.GetForm()["__body_"]; ok {
		return body.(string), ok
	}
	return "", false
}

// //GetHeader 获取头信息
// func (c *PkgCtx) GetHeader(key string) string {
// 	return c.Context.GetHeader(key)
// }
