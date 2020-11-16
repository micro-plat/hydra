package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

//TypeNodeName 分类节点名
const TypeNodeName = "router"

//Methods 支持的http请求类型
var Methods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodHead}

//DefMethods 普通服务包含的路由
var DefMethods = []string{http.MethodGet, http.MethodPost}

//GetWSHomeRouter 获取ws主页路由
func GetWSHomeRouter() *Router {
	return &Router{
		Path:    "/",
		Action:  Methods,
		Service: "/",
	}
}

//Routers 路由信息
type Routers struct {
	Routers []*Router `json:"routers,omitempty" toml:"routers,omitempty"`
}

func (h *Routers) String() string {
	var sb strings.Builder
	for _, v := range h.Routers {
		sb.WriteString(fmt.Sprintf("%-16s %-32s %-32s %v\n", v.Path, v.Service, strings.Join(v.Action, " "), v.Pages))
	}
	return sb.String()
}

//GetRouters 获取路由列表
func (h *Routers) GetRouters() []*Router {
	return h.Routers
}

//Router 路由信息
type Router struct {
	Path     string   `json:"path,omitempty" valid:"ascii,required" toml:"path,omitempty"`
	Action   []string `json:"action,omitempty" valid:"uppercase,in(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS)"  toml:"action,omitempty"`
	Service  string   `json:"service,omitempty" valid:"ascii,required" toml:"service,omitempty"`
	Encoding string   `json:"encoding,omitempty" toml:"encoding,omitempty"`
	Pages    []string `json:"pages,omitempty" toml:"pages,omitempty"`
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

//IsUTF8 是否是UTF8编码
func (r *Router) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
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

//GetPath 获取所有路由信息
func (h *Routers) GetPath() []string {
	list := make([]string, 0, len(h.Routers))
	for _, v := range h.Routers {
		list = append(list, v.Path)
	}
	return list
}

//GetConf 设置路由
func GetConf(cnf conf.IServerConf) (router *Routers, err error) {
	router = new(Routers)
	_, err = cnf.GetSubObject(TypeNodeName, router)
	if err == conf.ErrNoSetting || len(router.Routers) == 0 {
		return NewRouters(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取路由(%s)失败:%w", cnf.GetServerPath(), err)
	}
	if b, err := govalidator.ValidateStruct(router); !b {
		return nil, fmt.Errorf("路由(%s)配置有误:%w", cnf.GetServerPath(), err)
	}
	return router, nil
}
