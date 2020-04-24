package middleware

import (
	"strings"

	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

//Header 头设置
func Header(h header.IHeader) swap.Handler {
	return func(r swap.IRequest) {

		r.Next()

		headers, ok := h.GetConf()
		if !ok {
			return
		}
		if ok {
			origin := r.GetHeader("Origin")
			for k, v := range headers {
				if k != "Access-Control-Allow-Origin" { //非跨域设置
					r.Header(k, v)
					continue
				}
				if origin != "" && (v == "*" || strings.Contains(v, origin)) {
					r.Header(k, origin)
				}
			}
		}
	}
}
