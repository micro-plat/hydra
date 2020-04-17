package tmpl

import (
	"fmt"
	"text/template"

	"github.com/micro-plat/hydra/conf"
)

//RspTmplItem 服务响应模板项
type RspTmplItem struct {
	*option
}

//RspTmpl 响应模板信息
type RspTmpl struct {
	//模板项
	Items []*RspTmplItem `json:"items,omitempty" valid:"required"`

	//全局参数
	Params map[string]interface{} `json:"params,omitempty"`
}

//NewRspTmpl 构建响应配置
func NewRspTmpl(opts ...Option) *RspTmpl {
	r := &RspTmpl{Items: make([]*RspTmplItem, 0, 2), Params: make(map[string]interface{})}
	for _, opt := range opts {
		nopt := &option{}
		opt(nopt)
		r.append(nopt)
	}
	return r
}

//append 追加模板配置
func (r *RspTmpl) append(o *option) *RspTmpl {
	if o.Content != "" { //检查格式
		if _, err := template.New("").Parse(o.Content); err != nil {
			panic(fmt.Errorf("RspTmpl响应模板格式错误:%v", err))
		}
	}

	if o.Status != "" { //检查格式
		if _, err := template.New("").Parse(o.Status); err != nil {
			panic(fmt.Errorf("RspTmpl响应状态码格式错误:%v", err))
		}
	}
	r.Items = append(r.Items, &RspTmplItem{option: o})
	return r
}

//SetGlobalParam 追加模板配置
func (r *RspTmpl) SetGlobalParam(k string, v interface{}) *RspTmpl {
	r.Params[k] = v
	return r
}

//GetTemplate 获取指定请求对应的模板
func (r *RspTmpl) GetTemplate(s string) (bool, *RspTmplItem) {
	var last *RspTmplItem
	for _, t := range r.Items {
		if t.Status == "" && t.Content == "" {
			continue
		}
		for _, service := range t.Services {
			if service == s {
				last = t
				goto LOOP
			}
			if service == "*" {
				last = t
			}
		}
	}
LOOP:
	if last == nil {
		return false, nil
	}
	return true, last
}

//GetRspTmpl 设置GetRspTmpl配置
func GetRspTmpl(cnf conf.IServerConf) (response *conf.Response, err error) {
	_, err = cnf.GetSubObject("rsptmpl", &response)
	if err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return response, nil
}
