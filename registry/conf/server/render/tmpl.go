package render

import (
	"fmt"
	"text/template"

	"github.com/micro-plat/hydra/registry/conf"
)

//RenderItem 服务响应模板项
type RenderItem struct {
	*option
}

//Render 响应模板信息
type Render struct {
	//模板项
	Items []*RenderItem `json:"items,omitempty" valid:"required" toml:"items,omitempty"`

	//全局参数
	Params map[string]interface{} `json:"params,omitempty"`
}

//New 构建响应配置
func New(opts ...Option) *Render {
	r := &Render{Items: make([]*RenderItem, 0, 2), Params: make(map[string]interface{})}
	for _, opt := range opts {
		nopt := &option{}
		opt(nopt)
		r.append(nopt)
	}
	return r
}

//append 追加模板配置
func (r *Render) append(o *option) *Render {
	if o.Content != "" { //检查格式
		if _, err := template.New("").Parse(o.Content); err != nil {
			panic(fmt.Errorf("Render响应模板格式错误:%v", err))
		}
	}

	if o.Status != "" { //检查格式
		if _, err := template.New("").Parse(o.Status); err != nil {
			panic(fmt.Errorf("Render响应状态码格式错误:%v", err))
		}
	}
	r.Items = append(r.Items, &RenderItem{option: o})
	return r
}

//SetGlobalParam 追加模板配置
func (r *Render) SetGlobalParam(k string, v interface{}) *Render {
	r.Params[k] = v
	return r
}

//GetTemplate 获取指定请求对应的模板
func (r *Render) GetTemplate(s string) (bool, *RenderItem) {
	var last *RenderItem
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

//GetConf 设置GetRender配置
func GetConf(cnf conf.IMainConf) (rsp *Render, err error) {
	if _, err = cnf.GetSubObject("tmpl", &rsp); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return rsp, nil
}
