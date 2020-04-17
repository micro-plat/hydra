// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"errors"
	"fmt"
	"sync"
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
type wsHandler struct {
	// The websocket connection.
	conn      *websocket.Conn
	closeChan chan struct{}
	once      sync.Once
	// Buffered channel of outbound messages.
	send     chan []byte
	jwtToken string
}

//newWSHandler 构建ws处理程序
func newWSHandler(conn *websocket.Conn) *wsHandler {
	return &wsHandler{
		conn:      conn,
		closeChan: make(chan struct{}),
		send:      make(chan []byte, 256),
	}
}

//readPump 循环从读取客户端传入数据
func (c *wsHandler) readPump(exhandler interface{}, cx *gin.Context, conn *websocket.Conn, handler servers.IExecuter, ctn context.IContainer, name string, engine string, service string, mSetting map[string]string) {
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
			if !c.wsAction(exhandler, cx, conn, handler, ctn, name, engine, service, mSetting) {
				return
			}
		}
	}
}

//wsAction 调用内部服务处理类型逻辑
func (c *wsHandler) wsAction(exhandler interface{}, ctx *gin.Context, conn *websocket.Conn, handler servers.IExecuter, ctn context.IContainer, name string, engine string, service string, mSetting map[string]string) bool {
	defer gin.Recovery()(ctx)
	defer setExt(ctx, "CONN")
	//读取传入消息
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
		getLogger(ctx).Error(err)
		return false
	}

	//消息必须为json串
	input, err := jsons.Unmarshal([]byte(msg))
	if err != nil {
		err = fmt.Errorf("请求串不是有效的json:%s(%v)", msg, err)
		getLogger(ctx).Error(err)
		c.sendNow(ctx, "init", 406, err)
		return true
	}

	//获取消息中的service字段
	var ok bool
	if service, ok = input["service"].(string); !ok {
		err = errors.New("请求未包含服务名称字段")
		getLogger(ctx).Error(err)
		c.sendNow(ctx, service, 406, err)
		return true
	}

	//调用服务执行业务逻辑
	wLogHead(ctx, service)
	nctx := context.GetContext(exhandler, name,
		engine,
		service,
		ctn,
		makeQueyStringData(ctx),
		makeMapData(input),
		makeParamsData(ctx),
		makeSettingData(ctx, mSetting),
		makeExtData(ctx),
		getLogger(ctx))
	setServiceName(ctx, nctx.Service)
	setCTX(ctx, nctx)
	defer nctx.Close()
	defer wLogTail(ctx, service, time.Now())

	if c.jwtToken == "" {
		c.jwtToken = getJWTToken(ctx)
	}

	if !wsCheckJwt(ctx, service, c.jwtToken) { //jwt 验证未通过则强制关闭客户端
		err = fmt.Errorf("请求服务:%s未通过授权", service)
		getLogger(ctx).Error(err)
		c.sendNow(ctx, service, 406, err)
		nctx.Response.MustContent(406, err)
		return false
	}

	result := handler.Execute(nctx)
	if result != nil {
		nctx.Response.ShouldContent(result)
	}

	//处理错误err,5xx
	if err := nctx.Response.GetError(); err != nil {
		err = fmt.Errorf("error:%v", err)
		getLogger(ctx).Error(err)
		if !servers.IsDebug {
			err = errors.New("error:Internal Server Error")
		}
		nctx.Response.ShouldContent(err)
	}
	//处理跳转3xx
	if url, ok := nctx.Response.IsRedirect(); ok {
		nctx.Response.MustContent(nctx.Response.GetStatus(), map[string]interface{}{
			"status":   nctx.Response.GetStatus(),
			"location": url,
		})
	}

	//设置jwt验证
	if j, ok := makeJwtToken(ctx, nctx.Response.GetParams()["__jwt_"]); ok {
		c.jwtToken = j
	}
	//向客户端写入消息
	c.sendNow(ctx, service, nctx.Response.GetStatus(), nctx.Response.GetContent())
	return true
}

//close 关闭当前连接
func (c *wsHandler) close() {
	c.once.Do(func() {
		close(c.closeChan)
		close(c.send)
	})
}
