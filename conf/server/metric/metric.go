package metric

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName metric配置节点名
const TypeNodeName = "metric"

type IMetric interface {
	GetConf() (*Metric, bool)
}

//Metric Metric
type Metric struct {
	Host     string `json:"host,omitempty" valid:"requrl,required" toml:"host,omitempty"`
	DataBase string `json:"dataBase,omitempty" valid:"ascii,required" toml:"dataBase,omitempty"`
	Cron     string `json:"cron,omitempty" valid:"ascii,required" toml:"cron,omitempty"`
	UserName string `json:"userName,omitempty" valid:"ascii" toml:"userName,omitempty"`
	Password string `json:"password,omitempty" valid:"ascii" toml:"password,omitempty"`
	Disable  bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//New 构建api server配置信息
func New(host string, db string, cron string, opts ...Option) *Metric {
	//@ljy 2020-11-03 地址+端口由配置决定，避免配置域名的时候加上端口导致错误
	// port := ":8086"
	// if !strings.Contains(host, ":") {
	// 	host = host + port
	// }
	m := &Metric{
		Host:     host,
		DataBase: db,
		Cron:     cron,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

//GetConf 设置metric
func GetConf(cnf conf.IServerConf) (metric *Metric, err error) {
	metric = &Metric{}
	_, err = cnf.GetSubObject(TypeNodeName, metric)
	if err == conf.ErrNoSetting {
		metric.Disable = true
		return metric, nil
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(metric); !b {
		return nil, fmt.Errorf("metric配置数据有误:%v", err)
	}
	return
}
