package cron

import "github.com/micro-plat/hydra/conf"

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"status", "sharding"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{"task"}

//Server 服务嚣配置信息
type Server struct {
	*option
}

//New 构建cron server配置，默认为对等模式
func New(opts ...Option) *Server {
	s := &Server{option: &option{}}
	for _, opt := range opts {
		opt(s.option)
	}
	return s
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IMainConf) (s *Server, err error) {
	if _, err := cnf.GetMainObject(&s); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return s, nil
}
