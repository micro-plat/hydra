package header

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName header配置节点名
const TypeNodeName = "header"

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
	return strings.EqualFold(k, HeadeAllowOrigin)
}

//GetHTTPHeaderByOrigin 根据请求origin获取http请求头
func (h Headers) GetHTTPHeaderByOrigin(origin string) http.Header {
	hd := make(http.Header)
	header := h.GetHeaderByOrigin(origin)
	for k, v := range header {
		hd[k] = []string{v}
	}
	return hd
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
	value, ok := h[HeadeAllowOrigin]
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

//GetConf 设置header
func GetConf(cnf conf.IServerConf) (header Headers, err error) {
	rawConf, err := cnf.GetSubConf(TypeNodeName)
	if errors.Is(err, conf.ErrNoSetting) {
		return Headers{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("header配置有误:%v", err)
	}

	header = rawConf.ToSMap()
	return
}
