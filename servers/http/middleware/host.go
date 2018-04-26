package middleware

import (
	x "net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

//Host 处理服务器的主机头
func Host(cnf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hosts, ok := cnf.GetMetadata("hosts").(conf.Hosts)
		if !ok || len(hosts) == 0 {
			ctx.Next()
			return
		}
		if !checkHost(hosts, ctx) {
			getLogger(ctx).Errorf("必须使用host:%v访问", strings.Join(hosts, "或"))
			ctx.AbortWithStatus(x.StatusNotAcceptable)
			return
		}
		ctx.Next()
	}
}

func checkHost(hosts []string, ctx *gin.Context) bool {
	chost := ctx.Request.Host
	if i := strings.Index(chost, ":"); i != -1 {
		chost = chost[:i]
	}
	for _, host := range hosts {
		if host == chost {
			return true
		}
	}
	return false

}
