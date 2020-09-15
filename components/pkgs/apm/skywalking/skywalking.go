package skywalking

import (
	"encoding/json"
	"fmt"

	conf "github.com/micro-plat/hydra/conf/vars/apm"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	ctxapm "github.com/micro-plat/hydra/context/apm"
)

const (
	APMType = "skywalking"
)

//Client x
type Client struct {
	reporter ctxapm.Reporter
	instance string
}

// New 根据配置文件创建一个APM对象
func New(instance, raw string) (m *Client, err error) {
	conf := conf.APM{}
	err = json.Unmarshal([]byte(raw), &conf)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal:%s;%+v", raw, err)
		return
	}
	m = &Client{
		instance: instance,
	}
	m.reporter, err = NewReporter(conf.ServerAddress, &conf)
	if err != nil {
		err = fmt.Errorf("apm.NewReporter:%s;%+v", raw, err)
		return
	}
	return
}

func (c *Client) CreateTracer(service string) (tracer ctxapm.Tracer, err error) {
	return NewTracer(service, WithReporter(c.reporter), WithInstance(c.instance))
}

type skywalkingResolver struct {
}

func (s *skywalkingResolver) Resolve(instance, conf string) (apm.IAPM, error) {
	return New(instance, conf)
}
func init() {
	apm.Register(APMType, &skywalkingResolver{})
}
