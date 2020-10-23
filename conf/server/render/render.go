package render

import (
	"fmt"

	"github.com/asaskevich/govalidator"
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

//GetConf 设置GetRender配置
func GetConf(cnf conf.IMainConf) (rsp *Render, err error) {
	rsp = &Render{}
	_, err = cnf.GetSubObject("render", rsp)
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("render配置格式有误:%v", err)
	}
	if err == conf.ErrNoSetting {
		rsp.Disable = true
		return rsp, nil
	}
	paths := make([]string, 0, len(rsp.Tmplts))
	for k := range rsp.Tmplts {
		if b, err := govalidator.ValidateStruct(k); !b {
			return nil, fmt.Errorf("render Tmplt配置数据有误:%v", err)
		}
		paths = append(paths, k)
	}
	rsp.PathMatch = conf.NewPathMatch(paths...)
	return rsp, nil
}

//Get 获取转换结果
func (r *Render) Get(path string, funcs map[string]interface{}, i interface{}) (bool, int, string, string, error) {
	exists, service := r.PathMatch.Match(path, "/")
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
