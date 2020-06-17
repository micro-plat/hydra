package header

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
)

//IHeader 获取header
type IHeader interface {
	GetConf() (Headers, bool)
}

//Headers http头信息
type Headers map[string]string

//New 构建请求的头配置
func New(opts ...Option) Headers {
	h := make(map[string]string)
	for _, opt := range opts {
		opt(h)
	}
	return h
}

//IsAccessControlAllowOrigin 是否是Access-Control-Allow-Origin
func (h Headers) IsAccessControlAllowOrigin(k string) bool {
	return strings.EqualFold(k, "Access-Control-Allow-Origin")
}

//GetHeaderByOrigin 根据origin选择合适的header
func (h Headers) GetHeaderByOrigin(origin string) Headers {
	hd := New()
	if h.hasCross(origin) {
		for k, v := range h {
			hd[k] = v
			if h.IsAccessControlAllowOrigin(k) {
				hd[k] = origin
			}
		}
		return hd
	}
	for k, v := range h {
		if !strings.Contains(k, "Access-Control-Allow-") {
			hd[k] = v
		}
	}
	return hd
}

//hasCross 是否允许跨域访问
func (h Headers) hasCross(origin string) bool {
	value, ok := h["Access-Control-Allow-Origin"]
	if !ok {
		return false
	}
	if value == "*" {
		return true
	}
	vv := strings.Split(value, ",")
	for _, v := range vv {
		if v == origin { //allow
			return true
		}
	}
	return false
}

type ConfHandler func(cnf conf.IMainConf) Headers

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 设置header
func GetConf(cnf conf.IMainConf) (header Headers) {
	_, err := cnf.GetSubObject("header", &header)
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("header配置有误:%v", err))
	}
	if err == conf.ErrNoSetting {
		return nil
	}
	return
}
