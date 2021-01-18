package mock

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
)

//NewContext 创建mock类型的Context包
func NewContext(content string, opts ...Option) context.IContext {
	mk := newMock(content)
	for _, opt := range opts {
		opt(mk)
	}
	return ctx.NewCtx(mk, mk.serverType)
}
