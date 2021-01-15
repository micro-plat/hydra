package wacl

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
)

type cAPI struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
	Tags    []string `json:"tags"`
}

type cList struct {
	Page string `json:"page"`
	API  []*cAPI
}

//Pages 页面信息
type Pages struct {
	Path string
	Tags []string
}
type pathMapping map[string][]*Pages

type cLists []*cList

//ToMapping 将页面配置转换为路径映射
func (c cLists) ToMapping() pathMapping {
	p := make(map[string][]*Pages)
	for _, page := range c {
		for _, api := range page.API {
			methods := api.Methods
			if len(methods) == 0 {
				methods = router.Methods
			}
			for _, method := range methods {
				key := joinPath(api.Path, method)
				pages, ok := p[key]
				if !ok {
					pages = make([]*Pages, 0, 1)
					p[key] = pages
				}
				pages = append(pages, &Pages{Path: page.Page, Tags: api.Tags})
			}
		}
	}
	return p
}
func (m pathMapping) Keys() []string {
	list := []string{}
	for k := range m {
		list = append(list, k)
	}
	return list
}
func joinPath(path string, method string) string {
	return fmt.Sprintf("%s[%s]", path, strings.ToUpper(method))
}
