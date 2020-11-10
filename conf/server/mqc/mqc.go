package mqc

import (
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
var SubConfName = []string{"queue"}

//Server mqc服务配置
type Server struct {
	Status   string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty"`
	Sharding int    `json:"sharding,omitempty" toml:"sharding,omitempty"`
	Addr     string `json:"addr,omitempty" valid:"required"  toml:"addr,omitempty"`
	Trace    bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建mqc server配置，默认为对等模式
func New(addr string, opts ...Option) *Server {
	if _, _, err := global.ParseProto(addr); err != nil {
		panic(fmt.Errorf("mqc服务器地址配置有误，必须是:proto://configName 格式 %w", err))
	}
	s := &Server{Addr: addr, Status: StartStatus}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (*Server, error) {
	s := Server{}
	if cnf.GetServerType() != global.MQC {
		return nil, fmt.Errorf("mqc主配置类型错误:%s != mqc", cnf.GetServerType())
	}

	_, err := cnf.GetMainObject(&s)

	if err == conf.ErrNoSetting {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(&s); !b {
		return nil, fmt.Errorf("mqc服务器配置数据有误:%v %v", err, s)
	}
	return &s, nil
}
