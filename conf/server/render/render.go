package render

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

type Tmplt struct {
	Content string `json:"content,omitempty" valid:"required"`
	Status  string `json:"status,omitempty" valid:"required"`
}

//Render 响应模板信息
type Render struct {
	Tmplts map[string]*Tmplt `json:"tmplts,omitemptye" toml:"Tmplts,omitempty"`
	//Disable 禁用
	Disable bool `json:"disable,omitemptye" toml:"disable,omitempty"`
	*conf.Includes
}

//NewRender 构建模板
func NewRender(opts ...Option) *Render {
	r := &Render{Tmplts: make(map[string]*Tmplt)}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//GetConf 设置GetRender配置
func GetConf(cnf conf.IMainConf) (rsp *Render) {
	rsp = &Render{}
	_, err := cnf.GetSubObject("render", rsp)
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("render配置有误:%v", err))
	}
	if err == conf.ErrNoSetting {
		rsp.Disable = true
		return rsp
	}
	paths := make([]string, 0, len(rsp.Tmplts))
	for k := range rsp.Tmplts {
		paths = append(paths, k)
	}
	rsp.Includes = conf.NewInCludes(paths...)
	return rsp
}

//Get 获取转换结果
func (r *Render) Get(path string, funcs map[string]interface{}, i interface{}) (bool, int, string, error) {
	exists, service := r.Includes.In(path)
	if !exists {
		return false, 0, "", nil
	}
	var cstatus string
	var err error
	if r.Tmplts[service].Status != "" {
		cstatus, err = translate(service, r.Tmplts[service].Status, funcs, i)
		if err != nil || types.GetInt(cstatus) == 0 {
			return true, 0, "", fmt.Errorf("状态码模板%s配置有误 %w", r.Tmplts[service].Status, err)
		}
	}

	ccontent, err := translate(service, r.Tmplts[service].Content, funcs, i)
	if err != nil {
		return true, 0, "", fmt.Errorf("响应内容模板%s配置有误 %w", r.Tmplts[service].Content, err)
	}
	return true, types.GetInt(cstatus, 0), ccontent, nil
}
