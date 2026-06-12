package aigw

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/lib4go/types"
)

const (
	// DefaultAddress AIGW默认端口
	DefaultAddress = "8082"

	// DefaultRTimeOut 默认读取超时时间（0=不限制，由上游client和流式空闲超时控制）
	DefaultRTimeOut = 0

	// DefaultWTimeOut 默认写入超时时间（0=不限制，由上游client和流式空闲超时控制）
	DefaultWTimeOut = 0

	// DefaultRHTimeOut 默认请求头读取超时时间
	DefaultRHTimeOut = 30

	// DefaultStreamTimeOut 默认流式响应超时时间
	DefaultStreamTimeOut = 1800

	// StartStatus 开启服务
	StartStatus = "start"

	// StartStop 停止服务
	StartStop = "stop"
)

// MainConfName 主配置中的关键配置名
var MainConfName = []string{"address", "status", "rTimeout", "wTimeout", "rhTimeout", "streamTimeout", "dns"}

// SubConfName 子配置中的关键配置名
var SubConfName = []string{"router", "metric", "processor"}

var validTypes = map[string]bool{"aigw": true}

// Server AIGW server配置信息
type Server struct {
	security.ConfEncrypt
	Address       string `json:"address,omitempty" valid:"port,required" label:"端口号"`
	Status        string `json:"status,omitempty" valid:"in(start|stop)" label:"服务器状态"`
	RTimeout      int    `json:"rTimeout,omitempty" valid:"range(0|3600)" label:"请求读取超时时间"`
	WTimeout      int    `json:"wTimeout,omitempty" valid:"range(0|7200)" label:"请求处理写入时间"`
	RHTimeout     int    `json:"rhTimeout,omitempty" valid:"range(3|3600)"`
	StreamTimeout int    `json:"streamTimeout,omitempty" valid:"range(3|7200)"`
	Domain        string `json:"dns,omitempty" valid:"dns" toml:"dns,omitempty" label:"域名"`
	Name          string `json:"name,omitempty" toml:"name,omitempty" label:"服务器名称"`
	Trace         bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

// New 构建AIGW server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address:       address,
		Status:        StartStatus,
		RTimeout:      DefaultRTimeOut,
		WTimeout:      DefaultWTimeOut,
		RHTimeout:     DefaultRHTimeOut,
		StreamTimeout: DefaultStreamTimeOut,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// GetAddress 获取AIGW服务地址端口
func (s *Server) GetAddress() string {
	if types.IsEmpty(s.Address) {
		return DefaultAddress
	}
	return s.Address
}

// GetRTimeout 获取读取超时时间（0=不限制）
func (s *Server) GetRTimeout() int {
	return s.RTimeout
}

// GetWTimeout 获取写超时时间（0=不限制）
func (s *Server) GetWTimeout() int {
	return s.WTimeout
}

// GetRHTimeout 获取头读取超时时间
func (s *Server) GetRHTimeout() int {
	if s.RHTimeout <= 0 {
		return DefaultRHTimeOut
	}
	return s.RHTimeout
}

// GetStreamTimeout 获取流式响应超时时间
func (s *Server) GetStreamTimeout() int {
	if s.StreamTimeout <= 0 {
		return DefaultStreamTimeOut
	}
	return s.StreamTimeout
}

// GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	if _, ok := validTypes[cnf.GetServerType()]; !ok {
		return nil, fmt.Errorf("aigw主配置类型错误:%s != %+v", cnf.GetServerType(), validTypes)
	}
	s = &Server{}
	_, err = cnf.GetMainObject(s)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("aigw主配置数据有误:%v", err)
	}
	return s, nil
}
