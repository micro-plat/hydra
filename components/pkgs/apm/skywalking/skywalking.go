package skywalking

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	ctxapm "github.com/micro-plat/hydra/context/apm"
	"github.com/micro-plat/lib4go/types"
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
	confdata := types.XMap{}
	err = json.Unmarshal([]byte(raw), &confdata)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal:%s;%+v", raw, err)
		return
	}
	m = &Client{
		instance: instance,
	}
	serverAddr := confdata.GetString("server_address")
	m.reporter, err = NewReporter(serverAddr, raw)
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
