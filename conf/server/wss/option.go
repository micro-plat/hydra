package wss

import "github.com/micro-plat/hydra/global"

type Option interface {
	Apply(*Server)
	Type() string
}

type sideOption struct {
	tp   string
	opts []func(*Server)
}

func (s sideOption) Apply(c *Server) {
	for _, opt := range s.opts {
		opt(c)
	}
}

func (s sideOption) Type() string {
	return s.tp
}

func WithServerSide(opts ...func(*Server)) Option {
	return sideOption{tp: global.WSSServer, opts: opts}
}

func WithClientSide(opts ...func(*Server)) Option {
	return sideOption{tp: global.WSSClient, opts: opts}
}

func WithAddress(address string) func(*Server) {
	return func(s *Server) { s.Address = address }
}

func WithPath(path string) func(*Server) {
	return func(s *Server) { s.Path = path }
}

func WithServer(server string) func(*Server) {
	return func(s *Server) { s.Server = server }
}

func WithGroup(group string) func(*Server) {
	return func(s *Server) { s.Group = group }
}

func WithAuth(authType string, secret string) func(*Server) {
	return func(s *Server) {
		s.AuthType = authType
		s.AuthSecret = secret
	}
}

func WithClientID(clientID string) func(*Server) {
	return func(s *Server) { s.ClientID = clientID }
}

func WithTrace() func(*Server) {
	return func(s *Server) { s.Trace = true }
}

func WithHeartbeat(pingInterval, pongTimeout, writeTimeout int) func(*Server) {
	return func(s *Server) {
		s.PingInterval = pingInterval
		s.PongTimeout = pongTimeout
		s.WriteTimeout = writeTimeout
	}
}

func WithRequestTimeout(timeout int) func(*Server) {
	return func(s *Server) { s.RequestTimeout = timeout }
}

func WithReconnect(interval int) func(*Server) {
	return func(s *Server) { s.Reconnect = interval }
}
