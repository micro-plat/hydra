package pkg

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
)

type Package struct {
	URL     string `json:"url" valid:"requrl,required" toml:"url,omitempty"`
	Version string `json:"version" valid:"ascii,required" toml:"version,omitempty"`
	CRC32   uint32 `json:"crc32" valid:"required" toml:"crc32,omitempty"`
}

//NewPackage 构建CRON任务
func NewPackage(url string, version string, crc32 uint32) *Package {
	return &Package{
		URL:     url,
		Version: version,
		CRC32:   crc32,
	}
}

//GetConf 获取配置信息
func GetConf(cnf conf.IServerConf) (pkg *Package, err error) {
	_, err = cnf.GetSubObject("package", &pkg)
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("package配置有误:%v", err)
	}
	return
}
