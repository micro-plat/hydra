package middleware

import (
	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/servers/pkg/swap"
)

var originName = "Origin"

//Header 头设置
func Header(h header.IHeader) swap.Handler {
	return func(r swap.IRequest) {

		//1. 业务处理
		r.Next()

		//2. 获取header配置
		headers, ok := h.GetConf()
		if !ok {
			return
		}

		//3. 处理响应header参数
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
