package header

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry/conf"
)

//IHeader 获取header
type IHeader interface {
	GetConf() (Headers, bool)
}

//Headers http头信息
type Headers option

//NewHeader 构建请求的头配置
func NewHeader(opts ...Option) Headers {
	h := newOption()
	for _, opt := range opts {
		opt(h)
	}
	return h
}

//IsAccessControlAllowOrigin 是否允许跨域访问
func (h Headers) IsAccessControlAllowOrigin(k string) bool {
	return k == "Access-Control-Allow-Origin"
}

//AllowOrigin 是否允许跨域访问
func (h Headers) AllowOrigin(k string, v string, origin string) bool {
	return h.IsAccessControlAllowOrigin(k) && origin != "" && (v == "*" || strings.Contains(v, origin))
}

//GetConf 设置header
func GetConf(cnf conf.IMainConf) (header *Headers, err error) {
	_, err = cnf.GetSubObject("header", &header)
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("header配置有误:%v", err)
	}
	return
}
