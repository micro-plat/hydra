package nfs

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
)

//NFS 网络文件系统配置
type NFS struct {
	security.ConfEncrypt
	Local   string `json:"local,omitempty" toml:"local,omitempty"`
	Host    string `json:"host,omitempty" valid:"required"  toml:"host,omitempty" label:"mqc服务地址"`
	Disable bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

//New 构建mqc NFS配置，默认为对等模式
func New(addr string, opts ...Option) *NFS {
	s := &NFS{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (*NFS, error) {
	s := NFS{}
	_, err := cnf.GetMainObject(&s)
	if errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(&s); !b {
		return nil, fmt.Errorf("mqc服务器配置数据有误:%v %v", err, s)
	}
	return &s, nil
}
