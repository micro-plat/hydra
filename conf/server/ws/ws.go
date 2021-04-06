package ws

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
)

const (
	//StartStatus 开启服务
	StartStatus = "start"
	//StartStop 停止服务
	StartStop = "stop"
)

//Server api server配置信息
type Server struct {
	security.ConfEncrypt
	Address   string `json:"address,omitempty" valid:"dialstring" toml:"address,omitempty" label:"ws服务地址"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty"`
	RTimeout  int    `json:"rTimeout,omitempty" toml:"rTimeout,omitzero"`
	WTimeout  int    `json:"wTimeout,omitempty" toml:"wTimeout,omitzero"`
	RHTimeout int    `json:"rhTimeout,omitempty" toml:"rhTimeout,omitzero"`
	Host      string `json:"host,omitempty" toml:"host,omitempty"`
	Domain    string `json:"dns,omitempty" toml:"dns,omitempty"`
	Trace     bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建websocket server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address: address,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	_, err = cnf.GetMainObject(&s)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("websocket主配置数据有误:%v", err)
	}
	return s, nil
}
