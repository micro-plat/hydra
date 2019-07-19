package middleware

import (
	x "net/http"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Host 处理服务器的主机头
func Host(cnf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		hosts, ok := cnf.GetMetadata("hosts").(conf.Hosts)
		if !ok {
			ctx.Next()
			return
		}
		correct := checkHost(hosts, ctx)
		if !correct {
			getLogger(ctx).Errorf("访问被拒绝,必须使用:%v访问", hosts)
			ctx.AbortWithStatus(x.StatusNotAcceptable)
			return
		}
		ctx.Next()
	}
}

func checkHost(hosts []string, ctx *dispatcher.Context) bool {
	chost, ok := ctx.Request.GetHeader()["host"]
	if !ok {
		return true
	}
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
