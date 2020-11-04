/*
处理对外提供服务的路由管理，包括注册、获取路由列表
*/

package services

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
)

//API  路由信息
var API = NewORouter()

//WEB web服务的路由信息
var WEB = NewORouter()

//WS web socket路由信息
var WS = NewORouter()

//RPC rpc服务的路由信息
var RPC = NewORouter()

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

type ORouter struct {
	routers     *router.Routers
	pathRouters map[string]*pathRouter
}

//NewORouter 构建路由管理器
func NewORouter() *ORouter {
	return &ORouter{
		routers:     router.NewRouters(),
		pathRouters: make(map[string]*pathRouter),
	}
}

func (s *ORouter) Add(path string, service string, action []string, opts ...interface{}) error {
	if _, ok := s.pathRouters[path]; !ok {
		s.pathRouters[path] = newPathRouter(path)
	}
	if err := s.pathRouters[path].Add(service, action, opts...); err != nil {
		return err
	}
	return nil
}

//GetRouters 获取所有路由配置
func (s *ORouter) GetRouters() (*router.Routers, error) {
	if len(s.routers.Routers) != 0 {
		return s.routers, nil
	}
	for _, p := range s.pathRouters {
		routers, err := p.GetRouters()
		if err != nil {
			return nil, err
		}
		s.routers.Routers = append(s.routers.Routers, routers...)
	}
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
