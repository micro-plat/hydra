package cron

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

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"status", "sharding"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{"task"}

//Server 服务嚣配置信息
type Server struct {
	Status   string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty" label:"cron服务状态"`
	Sharding int    `json:"sharding,omitempty" toml:"sharding,omitempty"`
	Trace    bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建cron server配置，默认为对等模式
func New(opts ...Option) *Server {
	s := &Server{
		Status: StartStatus,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	s = &Server{}
	if cnf.GetServerType() != global.CRON {
		return nil, fmt.Errorf("cron主配置类型错误:%s != cron", cnf.GetServerType())
	}

	_, err = cnf.GetMainObject(s)

	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("cron主配置数据有误:%v", err)
	}
	return s, nil
}
