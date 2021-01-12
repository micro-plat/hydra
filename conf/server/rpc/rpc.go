package rpc

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
)

const (
	//StartStatus 开启服务
	StartStatus = "start"
	//StartStop 停止服务
	StartStop = "stop"
)

//DefaultMaxRecvMsgSize 最大默认接收字节数
const DefaultMaxRecvMsgSize = 1024 * 1024 * 20

//DefaultMaxSendMsgSize 最大默认发送字节数
const DefaultMaxSendMsgSize = 1024 * 1024 * 20

//DefaultRPCAddress rpc服务默认地址
const DefaultRPCAddress = ":8090"

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"address", "status", "rTimeout", "wTimeout", "rhTimeout", "dn"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{"router", "metric"}

//Server rpc server配置信息
type Server struct {
	Address        string `json:"address,omitempty" toml:"address,omitempty"`
	Status         string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty" label:"rpc服务状态"`
	Host           string `json:"host,omitempty" toml:"host,omitempty"`
	Domain         string `json:"dns,omitempty" toml:"dns,omitempty"`
	Trace          bool   `json:"trace,omitempty" toml:"trace,omitempty"`
	MaxRecvMsgSize int    `json:"maxRecvMsgSize,omitempty" toml:"maxRecvMsgSize,omitempty"`
	MaxSendMsgSize int    `json:"maxSendMsgSize,omitempty" toml:"maxSendMsgSize,omitempty"`
}

//New 构建rpc server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address: address,
		Status:  StartStatus,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

//GetMaxRecvMsgSize 获取最大接收字节数
func (c *Server) GetMaxRecvMsgSize() int {
	if c.MaxRecvMsgSize <= 0 {
		return DefaultMaxRecvMsgSize
	}

	return c.MaxRecvMsgSize
}

//GetMaxSendMsgSize 获取最大发送字节数
func (c *Server) GetMaxSendMsgSize() int {
	if c.MaxSendMsgSize <= 0 {
		return DefaultMaxSendMsgSize
	}

	return c.MaxSendMsgSize
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	s = &Server{}
	if cnf.GetServerType() != global.RPC {
		return nil, fmt.Errorf("rpc主配置类型错误:%s != rpc", cnf.GetServerType())
	}

	_, err = cnf.GetMainObject(s)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("rpc主配置数据有误:%v", err)
	}
	return s, nil
}
