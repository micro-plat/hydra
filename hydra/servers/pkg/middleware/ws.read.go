package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//readPump 循环从读取客户端传入数据
func (c *wsHandler) readPump() {
	defer func() {
		c.close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-c.closeChan:
			return
		default:
			c.wsAction()
		}
	}
}

//wsAction 调用内部服务处理类型逻辑
func (c *wsHandler) wsAction() {

	//读取传入消息
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
		return
	}

	//构建请求
	req, err := NewRequest(http.MethodGet, msg, c.uuid)
	if err != nil {
		c.sendNow("init", http.StatusNotAcceptable, err)
		return
	}

	//处理请求
	writer, err := c.engine.HandleRequest(req)
	if err != nil {
		c.sendNow(req.GetService(), http.StatusInternalServerError, err)
		return
	}
	if writer.Status() == http.StatusUnauthorized || writer.Status() == http.StatusNotAcceptable {
		websocket.IsUnexpectedCloseError(err, writer.Status())
		return
	}

	//向客户端写入消息
	c.sendNow(req.GetService(), writer.Status(), writer.Data)
	return
}

//close 关闭当前连接
func (c *wsHandler) close() {
	c.once.Do(func() {
		close(c.closeChan)
		close(c.send)
	})
}
