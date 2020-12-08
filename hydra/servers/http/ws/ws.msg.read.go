package ws

import (
	"fmt"
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
			if err := c.wsAction(); err != nil {
				c.log.Errorf("连接已断开:%v", err)
				c.close()
				return
			}
		}
	}
}

//wsAction 调用内部服务处理类型逻辑
func (c *wsHandler) wsAction() error {

	//读取消息
	tp, msg, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return fmt.Errorf("连接已关闭:%w", err)
		}
		return err
	}

	//检查消息类型
	if tp != websocket.TextMessage && tp != websocket.BinaryMessage {
		return nil
	}

	//构建请求
	req, err := NewRequest(http.MethodGet, msg, c.uuid, c.clientip)
	if err != nil {
		c.log.Errorf("消息有误:%v", err)
		c.sendNow("/ws.init", http.StatusNotAcceptable, err)
		return nil
	}

	//使用引擎执行请求
	writer, err := c.engine.HandleRequest(req)
	if err != nil {
		c.sendNow(req.GetService(), http.StatusInternalServerError, err)
		return nil
	}
	if writer.Status() == http.StatusUnauthorized || writer.Status() == http.StatusNotAcceptable {
		return fmt.Errorf("验证失败，断开网络连接:%d", writer.Status())
	}
	//向客户端写入消息
	// data, _ := base64.Decode(string(writer.Data()))
	c.sendNow(req.GetService(), writer.Status(), string(writer.Data()))
	return nil
}

//close 关闭当前连接
func (c *wsHandler) close() {
	c.once.Do(func() {
		close(c.closeChan)
		close(c.send)
	})
}
