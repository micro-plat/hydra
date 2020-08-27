package apm

import (
	"encoding/json"

	"github.com/micro-plat/hydra/conf"
)

//APM 消息队列配置
type APM struct {
	APMType          string            `json:"apmtype"`
	ServerAddress    string            `json:"server_address"`
	CheckInterval    int               `json:"check_interval"`
	InstanceProps    map[string]string `json:"instance_props"`
	MaxSendQueueSize int               `json:"max_send_queue_size"`
	Raw              []byte            `json:"-"`
}

//New 构建apm信息
func New(apmtype string, raw []byte) *APM {
	m := &APM{
		APMType: apmtype,
		Raw:     raw,
	}
	json.Unmarshal(raw, m)
	return m
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *APM, err error) {
	jc, err := cnf.GetConf(tp, name)
	if err != nil {
		return nil, err
	}
	return New(jc.GetString("apmtype"), jc.GetRaw()), nil
}
