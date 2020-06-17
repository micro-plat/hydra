package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

type Routers struct {
	Routers []*Router `json:"routers,omitempty"`
}

func (p *Routers) String() string {
	var sb strings.Builder
	for _, v := range p.Routers {
		sb.WriteString(fmt.Sprintf("%-16s %-32s %-32s %v\n", v.Path, v.Service, strings.Join(v.Action, " "), v.Pages))
	}
	return sb.String()
}

//Router 路由信息
type Router struct {
	Path     string   `json:"path,omitempty" valid:"ascii,required"`
	Action   []string `json:"action,omitempty" valid:"uppercase,in(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS)"`
	Service  string   `json:"service" valid:"ascii,required"`
	Encoding string   `json:"encoding,omitempty"`
	Pages    []string `json:"pages,omitempty"`
}

//NewRouter 构建路径配置
func NewRouter(path string, service string, action []string, opts ...Option) *Router {
	r := &Router{
		Path:    path,
		Action:  action,
		Service: service,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//GetEncoding 获取encoding配置，未配置时返回utf-8
func (r *Router) GetEncoding() string {
	if r.Encoding != "" {
		return r.Encoding
	}
	return "utf-8"
}

//IsUTF8 是否是UTF8编码
func (r *Router) IsUTF8() bool {
	return strings.ToLower(r.GetEncoding()) == "utf-8"
}

//NewRouters 构建路由
func NewRouters() *Routers {
	r := &Routers{
		Routers: make([]*Router, 0, 1),
	}
	return r
}

//Append 添加路由信息
func (h *Routers) Append(path string, service string, action []string, opts ...Option) *Routers {
	h.Routers = append(h.Routers, NewRouter(path, service, action, opts...))
	return h
}

//Match 根据请求路径匹配指定的路由配置
func (h *Routers) Match(path string, method string) *Router {
	if path == "" || method == http.MethodOptions || method == http.MethodHead {
		return &Router{
			Path:   path,
			Action: []string{method},
		}
	}
	for _, r := range h.Routers {
		if r.Path == path && types.StringContains(r.Action, method) {
			return r
		}
	}
	panic(fmt.Sprintf("未找到与[%s][%s]匹配的路由", path, method))
}

type ConfHandler func(cnf conf.IMainConf) *Routers

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 设置路由
func GetConf(cnf conf.IMainConf) (router *Routers) {
	router = new(Routers)
	_, err := cnf.GetSubObject("router", &router)
	if err == conf.ErrNoSetting || len(router.Routers) == 0 {
		router = NewRouters()
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("获取路由(%s)失败:%w", cnf.GetMainPath(), err))
	}
	if b, err := govalidator.ValidateStruct(router); !b {
		panic(fmt.Errorf("路由(%s)配置有误:%w", cnf.GetMainPath(), err))
	}
	return router
}
