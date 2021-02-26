package processor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/router"
)

//TypeNodeName processor配置节点名
const TypeNodeName = "processor"

//IProcessor IProcessor
type IProcessor interface {
	GetConf() (*Processor, bool)
}

//Processor Processor
type Processor struct {
	ServicePrefix string `json:"servicePrefix,omitempty" valid:"lowercase,maxstringlength(16),matches(^(/[a-z0-9]+)$)"  toml:"servicePrefix,omitempty" label:"服务前缀"`
}

//New 构建api server配置信息
func New(opts ...Option) *Processor {
	m := &Processor{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

//TreatRouters TreatRouters
func (p *Processor) TreatRouters(routers []*router.Router) {
	for i := range routers {
		routers[i].RealPath = fmt.Sprintf("%s%s", p.ServicePrefix, routers[i].Path)
		if !strings.Contains(strings.Join(routers[i].Action, ","), http.MethodOptions) {
			routers[i].Action = append(routers[i].Action, http.MethodOptions)
		}
	}
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
