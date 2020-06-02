package cache

import "github.com/micro-plat/hydra/conf"

//Cache 消息队列配置
type Cache struct {
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}

//New 构建DB连接信息
func New(proto string, raw []byte) *Cache {
	return &Cache{
		Proto: proto,
		Raw:   raw,
	}
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *Cache, err error) {
	jc, err := cnf.GetConf(tp, name)
	if err != nil {
		return nil, err
	}
	return New(jc.GetString("proto"), jc.GetRaw()), nil
}
