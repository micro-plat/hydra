package order

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type Input struct {
	ID   string `form:"id" json:"id" valid:"int,required"` //绑定输入参数，并验证类型否是否必须输入
	Name string `form:"name" json:"name"`
}
type BindHandler struct {
	container component.IContainer
	once      sync.Once
}

func NewBindHandler(container component.IContainer) (u *BindHandler) {
	return &BindHandler{container: container}
}
func (u *BindHandler) Handle(ctx *context.Context) (r interface{}) {

	uuid := ctx.Request.GetUUID()
	go func() {
	START:
		for {
			select {
			case <-time.After(time.Second):
				if err := context.WSExchange.Notify(uuid, time.Now().Unix()); err != nil {
					ctx.Log.Error(err)
					break START
				}
			}
		}
	}()
	return "success"
}
