package api

import (
	"errors"
	"fmt"

	"github.com/micro-plat/lib4go/types"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

const (

	//DefaultAPIAddress api服务默认端口号
	DefaultAPIAddress = "8080"

	//DefaultWSAddress ws服务默认端口号
	DefaultWSAddress = "8070"

	//DefaultWEBAddress web服务默认端口号
	DefaultWEBAddress = "8089"

	//DefaultRTimeOut 默认读取超时时间
	DefaultRTimeOut = 30

	//DefaultWTimeOut 默认写超时时间
	DefaultWTimeOut = 30

	//DefaultRHTimeOut 默认头读取超时时间
	DefaultRHTimeOut = 30

	//StartStatus 开启服务
	StartStatus = "start"

	//StartStop 停止服务
	StartStop = "stop"
)

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"address", "status", "rTimeout", "wTimeout", "rhTimeout", "dn"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{"router", "metric"}
var validTypes = map[string]bool{"api": true, "web": true, "ws": true}

//Server api server配置信息
type Server struct {
	Address   string `json:"address,omitempty" valid:"port,required" label:"端口号|请输入正确的端口号(1-65535)"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)"  label:"服务器状态"`
	RTimeout  int    `json:"rTimeout,omitempty" valid:"range(3|3600)" label:"请求读取超时时间|请输入正确的超时时间(3-3600)"`
	WTimeout  int    `json:"wTimeout,omitempty" valid:"range(3|3600)" label:"请求处理写入时间|请输入正确的超时时间(3-3600)"`
	RHTimeout int    `json:"rhTimeout,omitempty" valid:"range(3|3600)"`
	Domain    string `json:"dns,omitempty" valid:"dns" toml:"dns,omitempty" label:"域名"`
	Name      string `json:"name,omitempty" toml:"name,omitempty" label:"服务器名称"`
	Trace     bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建api server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address:   address,
		Status:    StartStatus,
		RTimeout:  DefaultRTimeOut,
		WTimeout:  DefaultWTimeOut,
		RHTimeout: DefaultRHTimeOut,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

//GetAPIAddress 获取api服务地址端口
func (s *Server) GetAPIAddress() string {
	if types.IsEmpty(s.Address) {
		return DefaultAPIAddress
	}
	return s.Address
}

//GetWSAddress 获取ws服务地址端口
func (s *Server) GetWSAddress() string {
	if types.IsEmpty(s.Address) {
		return DefaultWSAddress
	}
	return s.Address
}

//GetWEBAddress 获取web服务地址端口
func (s *Server) GetWEBAddress() string {
	if types.IsEmpty(s.Address) {
		return DefaultWEBAddress
	}
	return s.Address
}

//GetRTimeout 获取读取超时时间
func (s *Server) GetRTimeout() int {
	if s.RTimeout <= 0 {
		return DefaultRTimeOut
	}
	return s.RTimeout
}

//GetWTimeout 获取写超时时间
func (s *Server) GetWTimeout() int {
	if s.WTimeout <= 0 {
		return DefaultWTimeOut
	}
	return s.WTimeout
}

//GetRHTimeout 获取头读取超时时间
func (s *Server) GetRHTimeout() int {
	if s.RHTimeout <= 0 {
		return DefaultRHTimeOut
	}
	return s.RHTimeout
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	if _, ok := validTypes[cnf.GetServerType()]; !ok {
		return nil, fmt.Errorf("api主配置类型错误:%s != %+v", cnf.GetServerType(), validTypes)
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
		return nil, fmt.Errorf("api主配置数据有误:%v", err)
	}
	return s, nil
}
