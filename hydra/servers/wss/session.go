package wss

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type session struct {
	id     string
	group  string
	conn   *websocket.Conn
	mu     sync.Mutex
	pongAt time.Time
	closed chan struct{}
}

func newSession(id string, group string, conn *websocket.Conn) *session {
	return &session{id: id, group: group, conn: conn, pongAt: time.Now(), closed: make(chan struct{})}
}

func (s *session) write(frame *Frame) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.Marshal(frame)
	if err != nil {
		return err
	}
	return s.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (s *session) close() {
	select {
	case <-s.closed:
	default:
		close(s.closed)
		s.conn.Close()
	}
}

func (s *session) pong() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pongAt = time.Now()
}

func (s *session) alive(timeout time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.pongAt) <= timeout
}
