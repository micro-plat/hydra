package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/micro-plat/lib4go/jsons"
)

//writePump 向客户端写入响应消息
func (c *wsHandler) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.conn.Close()
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.close()
				break
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				c.close()
				break
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.close()
				break
			}
		case <-c.closeChan:
			break
		}
	}
}

func (c *wsHandler) sendNow(ctx *gin.Context, code int, i interface{}) {
	buff, err := getWSMessage(code, i)
	if err != nil {
		getLogger(ctx).Error(err)
		return
	}
	c.send <- buff
}
func getWSMessage(code interface{}, i interface{}) ([]byte, error) {
	var input interface{}
	switch v := i.(type) {
	case error:
		input = map[string]interface{}{
			"code": code,
			"err":  v.Error(),
		}
	default:
		input = map[string]interface{}{
			"code": code,
			"data": v,
		}
	}
	return jsons.Marshal(input)
}
func (c *wsHandler) recvNotify(ctx *gin.Context) func(...interface{}) error {
	return func(input ...interface{}) error {
		if len(input) == 0 {
			return nil
		}
		var code interface{}
		var i interface{}
		if len(input) == 1 {
			code = 200
			i = input[0]
		} else {
			code = input[0]
			i = input[1]
		}
		buff, err := getWSMessage(code, i)
		if err != nil {
			getLogger(ctx).Error(err)
			return err
		}
		getLogger(ctx).Info("ws.response", "PUSH", "/", code)
		c.send <- buff
		return nil
	}
}
