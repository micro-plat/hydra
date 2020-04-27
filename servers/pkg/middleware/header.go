package middleware

import (
	"strings"

	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

var originName = "Origin"

//Header 头设置
func Header(h header.IHeader) swap.Handler {
	return func(r swap.IRequest) {

		if strings.ToUpper(r.GetMethod()) != "OPTIONS" {
			r.Next()
		}

		headers, ok := h.GetConf()
		if !ok {
			return
		}
		if ok {
			origin := r.GetHeader(originName)
			for k, v := range headers {
				if !headers.IsAccessControlAllowOrigin(k) { //非跨域设置
					r.Header(k, v)
					continue
				}
				if headers.AllowOrigin(k, v, origin) {
					r.Header(k, origin)
				}
			}
		}
	}
}
