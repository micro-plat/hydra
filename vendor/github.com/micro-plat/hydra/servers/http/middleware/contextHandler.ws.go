package middleware

import (
	x "net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
)

//WSContextHandler api请求处理程序
func WSContextHandler(exhandler interface{}, name string, engine string, service string, mSetting map[string]string) gin.HandlerFunc {
	handler, ok := exhandler.(servers.IExecuter)
	if !ok {
		panic("不是有效的servers.IExecuter接口")
	}
	ctn, _ := exhandler.(context.IContainer)
	return func(c *gin.Context) {
		cnf := getMetadataConf(c)
		header := getCrossHeader(cnf, c)
		upgrader.CheckOrigin = func(r *x.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, header)
		if err != nil {
			getLogger(c).Error(err)
			c.AbortWithStatus(x.StatusNotAcceptable)
			return
		}
		h := newWSHandler(conn)
		context.WSExchange.Subscribe(getUUID(c), h.recvNotify(c))
		defer context.WSExchange.Unsubscribe(getUUID(c))

		go h.readPump(exhandler, c, conn, handler, ctn, name, engine, service, mSetting)
		h.writePump()
	}
}
