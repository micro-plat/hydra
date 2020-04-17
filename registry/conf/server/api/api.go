package api

import "github.com/micro-plat/hydra/conf"

//Server api server配置信息
type Server struct {
	Address string `json:"address,omitempty" valid:"dialstring"`
	*option
}

//New 构建api server配置信息
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

//GetHosts 获取hosts
func GetHosts(set ISetHosts, cnf conf.IMainConf) (hosts []string, err error) {
	hosts = cnf.GetStrings("host")
	return hosts, nil
}
