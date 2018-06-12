package context

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//SetHeader 设置http头
func (r *Response) SetHeader(name string, value string) {
	r.Params[name] = value
}
func (r *Response) SetHeaders(h map[string]string) {
	for k, v := range h {
		r.Params[k] = v
	}
}
func (r *Response) GetParams() map[string]interface{} {
	return r.Params
}
func (r *Response) SetParams(v map[string]interface{}) {
	r.Params = v
}
func (r *Response) SetParam(key string, v interface{}) {
	r.Params[key] = v
}

//GetHeaders 获取http头配置
func (r *Response) GetHeaders() map[string]string {
	header := make(map[string]string)
	for k, v := range r.Params {
		if !strings.HasPrefix(k, "__") && v != nil && k != "Status" && k != "Location" {
			switch v.(type) {
			case []string:
				list := v.([]string)
				for _, i := range list {
					header[k] = i
				}
			default:
				header[k] = types.GetString(v)
			}
		}
	}
	return header
}
