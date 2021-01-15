package wacl

import (
	"encoding/json"

	"github.com/micro-plat/hydra/conf"
)

//WebACL 访问控制列表
type WebACL struct {
	p    *conf.PathMatch
	list pathMapping
}

//New WEB权限访问控制列表
func New(data []byte) (*WebACL, error) {
	acl := &WebACL{}
	if err := acl.load(data); err != nil {
		return nil, err
	}
	acl.p = conf.NewPathMatch(acl.list.Keys()...)
	return acl, nil
}

//GetPages 获取路径对应的页面列表
func (a *WebACL) GetPages(apiPath string, method string) []*Pages {
	ok, path := a.p.Match(joinPath(apiPath, method))
	if !ok {
		return nil
	}
	return a.list[path]
}

func (a *WebACL) load(data []byte) error {
	list := make(cLists, 0, 1)
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	a.list = list.ToMapping()
	return nil
}
