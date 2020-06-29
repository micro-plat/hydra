package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/micro-plat/hydra/components/ws"
)

//upgrader 处理ws请求

//WSExecuteHandler 业务处理Handler
func WSExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {
		n, ok := ctx.Meta().Get("__context_")
		if !ok {
			panic("ws获取context错误，未获取到__context_对象")
		}
		c := n.(*gin.Context)

		conn, err := getUpgrader(c.Writer, c.Request, c.Request.Header)
		if err != nil {
			ctx.Response().Write(http.StatusNotAcceptable, fmt.Errorf("无法初始化ws.upgrader %w", err))
			return
		}

		//构建处理函数
		h := newWSHandler(conn, ctx.User().GetRequestID(), ctx.User().GetClientIP())
		ws.WSExchange.Subscribe(ctx.User().GetRequestID(), h.recvNotify(c))
		defer ws.WSExchange.Unsubscribe(ctx.User().GetRequestID())

		//异步读取与写入
		go h.readPump()
		h.writePump()
		ctx.Response().NoNeedWrite(c.Writer.Status())
	}
}
func getUpgrader(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{h.Get("Sec-WebSocket-Extensions")},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return upgrader.Upgrade(w, r, nil)
}
