package processor

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName processor配置节点名
const TypeNodeName = "processor"

//IProcessor IProcessor
type IProcessor interface {
	GetConf() (*Processor, bool)
}

//Processor Processor
type Processor struct {
	ServicePrefix string `json:"servicePrefix,omitempty" valid:"lowercase,maxstringlength(16),matches(^/?[a-z0-9]+$)"  toml:"servicePrefix,omitempty" label:"服务前缀"`
}

//New 构建api server配置信息
func New(opts ...Option) *Processor {
	m := &Processor{}
	for _, opt := range opts {
		opt(m)
	}
	if ok, err := govalidator.ValidateStruct(m); !ok {
		panic(fmt.Errorf("Processor配置数据有误:%v", err))
	}
	return m
}

//GetConf 设置Processor
func GetConf(cnf conf.IServerConf) (p *Processor, err error) {
	p = New()
	_, err = cnf.GetSubObject(TypeNodeName, p)
	if errors.Is(err, conf.ErrNoSetting) {
		return p, nil
	}
	if err != nil {
		return nil, err
	}

	if ok, err := govalidator.ValidateStruct(p); !ok {
		return nil, fmt.Errorf("Processor配置数据有误:%v", err)
	}
	return
}
