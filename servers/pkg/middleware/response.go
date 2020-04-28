package middleware

import (
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Response 处理服务器响应
func Response() swap.Handler {
	return func(r swap.IContext) {
		r.Next()
	}
}
