package header

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
)

//Headers http头信息
type Headers = option

//NewHeader 构建请求的头配置
func NewHeader(opts ...Option) Headers {
	h := newOption()
	for _, opt := range opts {
		opt(h)
	}
	return h
}

//GetHeaders 设置header
func GetHeaders(cnf conf.IServerConf) (header *conf.Headers, err error) {
	_, err = cnf.GetSubObject("header", &header)
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("header配置有误:%v", err)
	}
	return
}
