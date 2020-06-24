package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/micro-plat/hydra/components/ws"
	"github.com/micro-plat/hydra/conf/server/router"
)

//upgrader 处理ws请求
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//WSExecuteHandler 业务处理Handler
func WSExecuteHandler(service string, routers ...*router.Router) Handler {
	return func(ctx IMiddleContext) {
		n, ok := ctx.Meta().Get("__context_")
		if !ok {
			panic("ws获取context错误，未获取到__context_对象")
		}
		c := n.(*gin.Context)

		//构建ws处理对象
		headers := ctx.ServerConf().GetHeaderConf()
		hh := headers.GetHTTPHeaderByOrigin(ctx.Request().Path().GetHeader(originName))
		conn, err := upgrader.Upgrade(c.Writer, c.Request, hh)
		if err != nil {
			ctx.Response().Abort(http.StatusNotAcceptable, fmt.Errorf("无法初始化ws.upgrader %w", err))
			return
		}

		//构建处理函数
		h := newWSHandler(conn, ctx.User().GetRequestID(), routers...)
		ws.WSExchange.Subscribe(ctx.User().GetRequestID(), h.recvNotify(c))
		defer ws.WSExchange.Unsubscribe(ctx.User().GetRequestID())

		//异步读取与写入
		go h.readPump()
		h.writePump()
	}
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}
