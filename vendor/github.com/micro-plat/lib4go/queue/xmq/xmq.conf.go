package xmq

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type Conf struct {
	Address string `json:"address" valid:"dialstring,required"`
	Key     string `json:"key"`
}

func NewConf(j string) (*Conf, error) {
	conf := Conf{}
	err := json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(&conf); !b {
		err = fmt.Errorf("xmq 配置文件有误:%v", err)
		return nil, err
	}
	return &conf, nil
}
