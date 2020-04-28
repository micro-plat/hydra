package middleware

import (
	"strings"

	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Options 请求处理
func Options() swap.Handler {
	return func(r swap.IContext) {

		//options请求则自动不再进行后续处理
		if strings.ToUpper(r.GetMethod()) == "OPTIONS" {
			return
		}
		r.Next()

	}
}
