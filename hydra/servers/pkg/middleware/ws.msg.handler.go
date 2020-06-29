package middleware

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro-plat/lib4go/logger"
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

// wxHandler is a middleman between the websocket connection and the hub.
type wsHandler struct {
	engine    *wsEngine
	conn      *websocket.Conn
	closeChan chan struct{}
	once      sync.Once
	send      chan []byte
	uuid      string
	log       logger.ILogger
	clientip  string
}

//newWSHandler 使用新引擎进行业务处理
func newWSHandler(conn *websocket.Conn, uuid string, clientip string) *wsHandler {
	if wsInternalEngine == nil {
		panic("ws internal engine未初始化")
	}
	s := &wsHandler{
		conn:      conn,
		uuid:      uuid,
		closeChan: make(chan struct{}),
		send:      make(chan []byte, 256),
		log:       logger.GetSession("ws", uuid),
		clientip:  clientip,
		engine:    wsInternalEngine,
	}
	return s
}
