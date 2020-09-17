package apm

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//APM 消息队列配置
type APM struct {
	APMName string `json:"-"`
	//ServerAddress 单个地址（不支持多个）
	ServerAddress       string            `json:"server_address" toml:"server_address"`
	ReportCheckInterval int               `json:"report_check_interval,omitempty" toml:"report_check_interval,omitempty"`
	InstanceProps       map[string]string `json:"instance_props,omitempty" toml:"instance_props,omitempty"`
	MaxSendQueueSize    int               `json:"max_send_queue_size,omitempty" toml:"max_send_queue_size,omitempty"`
	Credentials         *Credential       `json:"transport_credentials,omitempty"  toml:"transport_credentials,omitempty"`
	AuthenticationKey   string            `json:"authentication_key,omitempty"  toml:"authentication_key,omitempty" `
	Raw                 []byte            `json:"-"`
}

type Credential struct {
	CertFile   string `json:"cert_file,omitempty" toml:"cert_file,omitempty"`
	ServerName string `json:"server_name,omitempty" toml:"server_name,omitempty"`
}

//New 构建apm信息
func New(apmname string, raw []byte) *APM {
	m := &APM{
		APMName: apmname,
		Raw:     raw,
	}
	err := json.Unmarshal(raw, m)
	if err != nil {
		panic(fmt.Sprintf("新建apm配置失败：%s,content:%s", err, string(raw)))
	}
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
