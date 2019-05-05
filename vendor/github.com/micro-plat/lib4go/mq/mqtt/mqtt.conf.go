package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type Conf struct {
	Address  string `json:"address" valid:"dialstring,required"`
	DumpData bool   `json:"dump"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	CertPath string `json:"cert"`
}

func NewConf(j string) (*Conf, error) {
	conf := Conf{}
	err := json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(&conf); !b {
		err = fmt.Errorf("mqtt 配置文件有误:%v", err)
		return nil, err
	}
	return &conf, nil
}
