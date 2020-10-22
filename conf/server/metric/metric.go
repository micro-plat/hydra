package metric

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

type IMetric interface {
	GetConf() (*Metric, bool)
}

type Metric struct {
	Host     string `json:"host" valid:"requrl,required" toml:"host,omitempty"`
	DataBase string `json:"dataBase" valid:"ascii,required" toml:"dataBase,omitempty"`
	Cron     string `json:"cron" valid:"ascii,required" toml:"cron,omitempty"`
	UserName string `json:"userName,omitempty" valid:"ascii" toml:"userName,omitempty"`
	Password string `json:"password,omitempty" valid:"ascii" toml:"password,omitempty"`
	Disable  bool   `json:"disable,omitempty" toml:"disable,omitempty"`
	*option
}

//New 构建api server配置信息
func New(host string, db string, cron string, opts ...Option) *Metric {
	port := ":8086"
	if !strings.Contains(host, ":") {
		host = host + port
	}
	m := &Metric{
		Host:     host,
		DataBase: db,
		Cron:     cron,
		option:   &option{},
	}
	for _, opt := range opts {
		opt(m.option)
	}
	return m
}

//GetConf 设置metric
func GetConf(cnf conf.IMainConf) (metric *Metric, err error) {
	metric = &Metric{}
	_, err = cnf.GetSubObject("metric", &metric)
	if err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	if err == conf.ErrNoSetting {
		metric.Disable = true
		return metric, nil
	}
	if b, err := govalidator.ValidateStruct(metric); !b {
		return nil, fmt.Errorf("metric配置数据有误:%v", err)
	}
	return
}
