/*
处理对外提供服务的路由管理，包括注册、获取路由列表
*/

package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
)

//API  路由信息
var API = NewORouter("API")

//WEB web服务的路由信息
var WEB = NewORouter("WEB")

//WS web socket路由信息
var WS = NewORouter("WS")

//RPC rpc服务的路由信息
var RPC = NewORouter("RPC")

//GetRouter 获取服务器的路由配置
func GetRouter(tp string) *ORouter {
	switch tp {
	case global.API:
		return API
	case global.Web:
		return WEB
	case global.WS:
		return WS
	case global.RPC:
		return RPC
	default:
		panic(fmt.Sprintf("无法获取服务%s的路由配置", tp))
	}
}

//ORouter ORouter
type ORouter struct {
	name        string
	routers     *router.Routers
	pathRouters map[string]*pathRouter
}

//NewORouter 构建路由管理器
func NewORouter(name string) *ORouter {
	return &ORouter{
		name:        name,
		pathRouters: make(map[string]*pathRouter),
	}
}

//Add 添加路由
func (s *ORouter) Add(path string, service string, action []string, opts ...interface{}) error {
	if _, ok := s.pathRouters[path]; !ok {
		s.pathRouters[path] = newPathRouter(path)
	}
	if err := s.pathRouters[path].Add(service, action, opts...); err != nil {
		return err
	}
	return nil
}

//BuildRouters 根据前缀处理路由数据
func (s *ORouter) BuildRouters(prefix string) (*router.Routers, error) {
	s.routers = router.NewRouters()
	s.routers.ServicePrefix = prefix
	for _, prouter := range s.pathRouters {
		routers, err := prouter.GetRouters()
		if err != nil {
			return nil, err
		}
		tmplist := make([]*router.Router, len(routers))
		for i, r := range routers {
			var t = *r
			tmplist[i] = &t
			tmplist[i].Path = fmt.Sprintf("%s%s", prefix, r.Path)
			s.fillActs(r)
		}
		s.routers.Routers = append(s.routers.Routers, tmplist...)
	}

	for p, acts := range s.routers.MapPath {
		if _, ok := acts[http.MethodOptions]; !ok {
			service := fmt.Sprintf("%s$%s", p, http.MethodOptions)
			s.routers.Routers = append(s.routers.Routers, router.NewRouter(
				p,
				service,
				[]string{http.MethodOptions},
			))
			acts[http.MethodOptions] = service
			s.routers.MapPath[p] = acts
		}
	}

	return s.routers, nil
}

func (s *ORouter) fillActs(r *router.Router) {
	array, ok := s.routers.MapPath[r.Path]
	if !ok {
		array = make(map[string]string)
	}
	for i := range r.Action {
		array[r.Action[i]] = r.Service
	}
	s.routers.MapPath[r.Path] = array
}

//GetRouters 获取所有路由配置
func (s *ORouter) GetRouters() (*router.Routers, error) {
	return s.routers, nil
}

//处理管理相同服务名称的所有路由，当action重复时报错，对未指定时进行排除选择
type pathRouter struct {
	path    string
	action  map[string]string
	routers []*router.Router
}

func newPathRouter(path string) *pathRouter {
	r := &pathRouter{
		path:    path,
		action:  make(map[string]string),
		routers: make([]*router.Router, 0, 1),
	}
	for _, m := range router.Methods {
		r.action[m] = m
	}
	return r
}

func (p *pathRouter) Add(service string, action []string, opts ...interface{}) error {
	for _, a := range action {
		if _, ok := p.action[a]; !ok {
			return fmt.Errorf("服务%s不能重复注册", service)
		}
		delete(p.action, a)
	}

	//转换输入参数
	nopts := make([]router.Option, 0, len(opts))
	for _, p := range opts {
		v, ok := p.(router.Option)
		if !ok {
			panic(fmt.Errorf("%s注册的服务类型必须是router.Option", service))
		}
		nopts = append(nopts, v)
	}
	p.routers = append(p.routers, router.NewRouter(p.path, service, action, nopts...))
	return nil
}
func (p *pathRouter) GetRouters() ([]*router.Router, error) {
	//检查配置中是否有且只有一个路由未指定action
	var first = -1
	var hasRepeat bool
	for i, r := range p.routers {
		if len(r.Action) == 0 {
			if first != -1 {
				hasRepeat = true
			}
			if !hasRepeat {
				first = i
			}
		}
	}

	//没有未指定的action,直接返回
	if first == -1 {
		return p.routers, nil
	}

	//有多个需要指定action
	if hasRepeat {
		return nil, fmt.Errorf("重复注册的服务%v", p.routers[first])
	}

	action := make([]string, 0, len(p.action))
	for _, a := range p.action {
		for _, b := range router.DefMethods {
			if a == b {
				action = append(action, a)
			}
		}
	}

	if len(action) == 0 {
		return nil, fmt.Errorf("服务%s无法注册，所有action已分配", p.routers[first].Path)
	}
	p.routers[first].Action = action
	return p.routers, nil

}
func (p *pathRouter) String() string {
	var sb strings.Builder
	for _, v := range p.routers {
		sb.WriteString(fmt.Sprintf("%-32s %-32s\n", v.Service, strings.Join(v.Action, " ")))
	}
	return sb.String()
}
