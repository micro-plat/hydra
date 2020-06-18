package api

import (
	"github.com/micro-plat/hydra/conf"
)

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"address", "status", "rTimeout", "wTimeout", "rhTimeout", "dn"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{"router", "metric"}

//Server api server配置信息
type Server struct {
	Address   string `json:"address,omitempty" valid:"dialstring" toml:"address,omitempty"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty"`
	RTimeout  int    `json:"rTimeout,omitempty" toml:"rTimeout,omitzero"`
	WTimeout  int    `json:"wTimeout,omitempty" toml:"wTimeout,omitzero"`
	RHTimeout int    `json:"rhTimeout,omitempty" toml:"rhTimeout,omitzero"`
	Host      string `json:"host,omitempty" toml:"host,omitempty"`
	Domain    string `json:"dn,omitempty" toml:"dn,omitempty"`
	Trace     bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建api server配置信息
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
func GetConf(cnf conf.IMainConf) (s *Server, err error) {
	if _, err := cnf.GetMainObject(&s); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return s, nil
}
