package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type RequestHandler struct {
	container component.IContainer
}

func NewRequestHandler(container component.IContainer) (u *RequestHandler) {
	return &RequestHandler{container: container}
}
func (u *RequestHandler) Handle(ctx *context.Context) (r interface{}) {

	db, err := u.container.GetDB()
	if err != nil {
		return err
	}
	row, _, _, err := db.Execute("update sys_system_dictionary t set t.sort_id=0 where t.id=1", map[string]interface{}{})
	if err != nil {
		return err
	}
	return row
}

func (u *RequestHandler) Close() error {
	return nil
}
