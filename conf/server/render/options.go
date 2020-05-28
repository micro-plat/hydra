package render

import "github.com/micro-plat/lib4go/types"

type Option func(*Render)

//WithTmplt 添加模板
func WithTmplt(path string, content string, status ...string) Option {
	return func(a *Render) {
		if _, ok := a.Tmplts[path]; !ok {
			a.Tmplts[path] = &Tmplt{Content: content, Status: types.GetStringByIndex(status, 0, "")}
		}
	}
}
