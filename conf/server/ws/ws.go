package ws

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//Server api server配置信息
type Server struct {
	Address string `json:"address,omitempty" valid:"dialstring" toml:"address,omitempty"`
	*option
}

//New 构建websocket server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address: address,
		option:  &option{},
	}
	for _, opt := range opts {
		opt(a.option)
	}
	return a
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	if _, err := cnf.GetMainObject(&s); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("websocket主配置数据有误:%v", err)
	}
	return s, nil
}
