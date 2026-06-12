package wss

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

const (
	DefaultAddress        = "8443"
	DefaultPath           = "/hydra/wss"
	DefaultPingInterval   = 25
	DefaultPongTimeout    = 75
	DefaultWriteTimeout   = 10
	DefaultRequestTimeout = 60
	DefaultReconnect      = 10
	StartStatus           = "start"
	StartStop             = "stop"
)

var MainConfName = []string{"address", "server", "status", "path", "group", "clientID", "authType", "authSecret", "pingInterval", "pongTimeout", "writeTimeout", "requestTimeout", "reconnectInterval", "trace"}
var SubConfName = []string{"routes", "metric", "processor"}

type Server struct {
	security.ConfEncrypt
	Address        string `json:"address,omitempty" toml:"address,omitempty"`
	Server         string `json:"server,omitempty" toml:"server,omitempty"`
	Status         string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty"`
	Path           string `json:"path,omitempty" toml:"path,omitempty"`
	Group          string `json:"group,omitempty" toml:"group,omitempty"`
	ClientID       string `json:"clientID,omitempty" toml:"clientID,omitempty"`
	AuthType       string `json:"authType,omitempty" toml:"authType,omitempty"`
	AuthSecret     string `json:"authSecret,omitempty" toml:"authSecret,omitempty"`
	PingInterval   int    `json:"pingInterval,omitempty" toml:"pingInterval,omitempty"`
	PongTimeout    int    `json:"pongTimeout,omitempty" toml:"pongTimeout,omitempty"`
	WriteTimeout   int    `json:"writeTimeout,omitempty" toml:"writeTimeout,omitempty"`
	RequestTimeout int    `json:"requestTimeout,omitempty" toml:"requestTimeout,omitempty"`
	Reconnect      int    `json:"reconnectInterval,omitempty" toml:"reconnectInterval,omitempty"`
	Trace          bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

func NewServerSide(opts ...Option) *Server {
	s := newDefault()
	s.Address = DefaultAddress
	for _, opt := range opts {
		opt.Apply(s)
	}
	return s
}

func NewClientSide(opts ...Option) *Server {
	s := newDefault()
	for _, opt := range opts {
		opt.Apply(s)
	}
	return s
}

func newDefault() *Server {
	return &Server{
		Status:         StartStatus,
		Path:           DefaultPath,
		PingInterval:   DefaultPingInterval,
		PongTimeout:    DefaultPongTimeout,
		WriteTimeout:   DefaultWriteTimeout,
		RequestTimeout: DefaultRequestTimeout,
		Reconnect:      DefaultReconnect,
	}
}

func (s *Server) GetAddress() string {
	return types.GetString(s.Address, DefaultAddress)
}

func (s *Server) GetPath() string {
	return types.GetString(s.Path, DefaultPath)
}

func (s *Server) GetPingInterval() int {
	if s.PingInterval <= 0 {
		return DefaultPingInterval
	}
	return s.PingInterval
}

func (s *Server) GetPongTimeout() int {
	if s.PongTimeout <= 0 {
		return DefaultPongTimeout
	}
	return s.PongTimeout
}

func (s *Server) GetWriteTimeout() int {
	if s.WriteTimeout <= 0 {
		return DefaultWriteTimeout
	}
	return s.WriteTimeout
}

func (s *Server) GetRequestTimeout() int {
	if s.RequestTimeout <= 0 {
		return DefaultRequestTimeout
	}
	return s.RequestTimeout
}

func (s *Server) GetReconnect() int {
	if s.Reconnect <= 0 {
		return DefaultReconnect
	}
	return s.Reconnect
}

func GetConf(cnf conf.IServerConf) (*Server, error) {
	tp := cnf.GetServerType()
	if tp != global.WSSServer && tp != global.WSSClient {
		return nil, fmt.Errorf("wss主配置类型错误:%s", tp)
	}
	s := &Server{}
	_, err := cnf.GetMainObject(s)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}
	if s.Status == "" {
		s.Status = StartStatus
	}
	if s.Path == "" {
		s.Path = DefaultPath
	}
	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("wss主配置数据有误:%v", err)
	}
	return s, nil
}
