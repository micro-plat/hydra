package render

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

type Tmplt struct {
	ContentType string `json:"content_type,omitempty"`
	Content     string `json:"content,omitempty"`
	Status      string `json:"status,omitempty" valid:"required"`
}

//Render 响应模板信息
type Render struct {
	Tmplts map[string]*Tmplt `json:"tmplts,omitemptye" toml:"Tmplts,omitempty"`
	//Disable 禁用
	Disable bool `json:"disable,omitemptye" toml:"disable,omitempty"`
	*conf.PathMatch
}

//NewRender 构建模板
func NewRender(opts ...Option) *Render {
	r := &Render{Tmplts: make(map[string]*Tmplt)}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type ConfHandler func(cnf conf.IMainConf) *Render

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
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
	rsp.PathMatch = conf.NewPathMatch(paths...)
	return rsp
}

//Get 获取转换结果
func (r *Render) Get(path string, funcs map[string]interface{}, i interface{}) (bool, int, string, string, error) {
	exists, service := r.PathMatch.Match(path)
	if !exists {
		return false, 0, "", "", nil
	}
	tmpltStatus, tmpltContentType, tmpltContent := "", "", ""
	var err error
	if r.Tmplts[service].Status != "" {
		tmpltStatus, err = translate(service, r.Tmplts[service].Status, funcs, i)
		if err != nil || types.GetInt(tmpltStatus) == 0 {
			return true, 0, "", "", fmt.Errorf("status模板%s配置有误 %w", r.Tmplts[service].Status, err)
		}
	}
	if r.Tmplts[service].ContentType != "" {
		tmpltContentType, err = translate(service, r.Tmplts[service].ContentType, funcs, i)
		if err != nil {
			return true, 0, "", "", fmt.Errorf("content_type模板%s配置有误 %w", r.Tmplts[service].ContentType, err)
		}
	}

	tmpltContent, err = translate(service, r.Tmplts[service].Content, funcs, i)
	if err != nil {
		return true, 0, "", "", fmt.Errorf("响应内容模板%s配置有误 %w", r.Tmplts[service].Content, err)
	}
	return true, types.GetInt(tmpltStatus, 0), tmpltContentType, tmpltContent, nil
}
