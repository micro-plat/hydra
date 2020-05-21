package services

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry/conf/server/router"
)

var defRequestMethod = []string{"get", "post", "put", "delete"}

//处理相同路由名称不能action问题
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
	for _, m := range defRequestMethod {
		r.action[m] = m
	}
	return r
}
func (p *pathRouter) Add(service string, action []string) error {
	for _, a := range action {
		if _, ok := p.action[a]; !ok {
			return fmt.Errorf("服务%s不能重复注册", service)
		}
		delete(p.action, a)
	}
	p.routers = append(p.routers, router.NewRouter(p.path, service, action...))
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

	//只有一个需要指定，但没有action可指定
	if len(p.action) == 0 {
		return nil, fmt.Errorf("服务%v无法被使用，所有action已分配", p.routers[first])
	}
	action := make([]string, 0, len(p.action))
	for _, a := range p.action {
		action = append(action, a)
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
