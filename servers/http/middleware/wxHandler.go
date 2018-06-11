// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/jsons"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// wxHandler is a middleman between the websocket connection and the hub.
type wxHandler struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *wxHandler) readPump(cx *gin.Context, conn *websocket.Conn, handler servers.IExecuter, ctn context.IContainer, name string, engine string, service string, mSetting map[string]string) {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			getLogger(cx).Error(err)
			break
		}
		c.wsAction(cx, conn, handler, ctn, name, engine, service, msg, mSetting)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *wxHandler) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (cli *wxHandler) wsAction(c *gin.Context, conn *websocket.Conn, handler servers.IExecuter, ctn context.IContainer, name string, engine string, service string, msg []byte, mSetting map[string]string) {
	input, err := jsons.Unmarshal([]byte(msg))
	if err != nil {
		getLogger(c).Error(err)
		cli.sendNow(c, err)
		return
	}
	var ok bool
	if service, ok = input["service"].(string); !ok {
		err = errors.New("请求未包含服务名称字段")
		getLogger(c).Error(err)
		cli.sendNow(c, err)
		return
	}
	ctx := context.GetContext(name, engine, service, ctn, makeQueyStringData(c), makeMapData(input), makeParamsData(c), makeSettingData(c, mSetting), makeExtData(c), getLogger(c))

	defer setServiceName(c, ctx.Service)
	defer setCTX(c, ctx)
	//调用执行引擎进行逻辑处理

	result := handler.Execute(ctx)
	if result != nil {
		ctx.Response.ShouldContent(result)
	}
	//处理错误err,5xx
	if err := ctx.Response.GetError(); err != nil {
		err = fmt.Errorf("error:%v", err)
		if !servers.IsDebug {
			err = errors.New("error:Internal Server Error")
		}
		ctx.Response.ShouldContent(err)
	}
	//处理跳转3xx
	if url, ok := ctx.Response.IsRedirect(); ok {
		ctx.Response.MustContent(ctx.Response.GetStatus(), map[string]interface{}{
			"status": ctx.Response.GetStatus(),
			"url":    url,
		})
	}
	cli.sendNow(c, ctx.Response.GetContent())

}
func (cli *wxHandler) sendNow(c *gin.Context, i interface{}) {
	buff, err := getWSMessage(i)
	if err != nil {
		getLogger(c).Error(err)
		return
	}
	cli.send <- buff
}
func getWSMessage(i interface{}) ([]byte, error) {
	var input interface{}
	switch v := i.(type) {
	case string:
		input = map[string]interface{}{
			"data": v,
		}
	case bool:
		input = map[string]interface{}{
			"data": v,
		}
	case int, float32, float64:
		input = map[string]interface{}{
			"data": v,
		}
	case error:
		input = map[string]interface{}{
			"err": v.Error(),
		}
	default:
		input = i
	}
	return jsons.Marshal(input)
}
