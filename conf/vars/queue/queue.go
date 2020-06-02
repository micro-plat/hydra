package queue

import "github.com/micro-plat/hydra/conf"

//Queue 消息队列配置
type Queue struct {
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}

//New 构建DB连接信息
func New(proto string, raw []byte) *Queue {
	return &Queue{
		Proto: proto,
		Raw:   raw,
	}
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *Queue, err error) {
	jc, err := cnf.GetConf(tp, name)
	if err != nil {
		return nil, err
	}
	return New(jc.GetString("proto"), jc.GetRaw()), nil
}
