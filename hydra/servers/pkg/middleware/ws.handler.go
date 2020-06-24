package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
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
	engine *dispatcher.Engine
	// The websocket connection.
	conn      *websocket.Conn
	closeChan chan struct{}
	once      sync.Once
	// Buffered channel of outbound messages.
	send chan []byte
	uuid string
}

//newWSHandler 使用新引擎进行业务处理
func newWSHandler(conn *websocket.Conn, uuid string, routers ...*router.Router) *wsHandler {
	s := &wsHandler{
		conn:      conn,
		uuid:      uuid,
		closeChan: make(chan struct{}),
		send:      make(chan []byte, 256),
		engine:    dispatcher.New(),
	}
	s.engine.Use(Recovery().DispFunc(global.WS))
	s.engine.Use(Logging().DispFunc())   //记录请求日志
	s.engine.Use(BlackList().DispFunc()) //黑名单控制
	s.engine.Use(WhiteList().DispFunc()) //白名单控制
	s.engine.Use(Trace().DispFunc())     //跟踪信息
	s.engine.Use(Delay().DispFunc())     //
	s.engine.Use(Options().DispFunc())   //处理option响应
	s.engine.Use(Static().DispFunc())    //处理静态文件
	s.engine.Use(Header().DispFunc())    //设置请求头
	s.engine.Use(BasicAuth().DispFunc()) //
	s.engine.Use(APIKeyAuth().DispFunc())
	s.engine.Use(RASAuth().DispFunc())
	s.engine.Use(JwtAuth().DispFunc())   //jwt安全认证
	s.engine.Use(Render().DispFunc())    //响应渲染组件
	s.engine.Use(JwtWriter().DispFunc()) //设置jwt回写
	// s.engine.Use(s.metric.Handle().DispFunc()) //生成metric报表
	s.addWSRouter(routers...)
	return s
}
func (s *wsHandler) addWSRouter(routers ...*router.Router) {
	for _, router := range routers {
		for _, method := range router.Action {
			s.engine.Handle(strings.ToUpper(method), router.Path, ExecuteHandler(router.Service).DispFunc())
		}
	}
}
