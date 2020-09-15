package apm

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/asaskevich/govalidator"
)

//APM 消息队列配置
type APM struct {
	APMName string `json:"-"`
	//ServerAddress 单个地址（不支持多个）
	ServerAddress       string            `json:"server_address,omitempty" toml:"server_address,omitempty"`
	ReportCheckInterval int               `json:"report_check_interval,omitempty" toml:"report_check_interval,omitempty"`
	InstanceProps       map[string]string `json:"instance_props,omitempty" toml:"instance_props,omitempty"`
	MaxSendQueueSize    int               `json:"max_send_queue_size,omitempty" toml:"max_send_queue_size,omitempty"`
	Raw                 []byte            `json:"-"`
}

//New 构建apm信息
func New(apmname string, raw []byte) *APM {
	m := &APM{
		APMName: apmname,
		Raw:     raw,
	}
	json.Unmarshal(raw, m)
	return m
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *APM, err error) {
	s = &APM{}
	_, err = cnf.GetObject(tp, name, s)

	if err != nil {
		panic(fmt.Errorf("读取./var/%s/%s 配置发生错误 %w", tp, name, err))
	}

	if b, err := govalidator.ValidateStruct(s); !b {
		panic(fmt.Errorf("./var/%s/%s 配置有误 %w", tp, name, err))
	}

	return
}
