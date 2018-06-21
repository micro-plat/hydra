package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type QueryHandler struct {
	container component.IContainer
}

func NewQueryHandler(container component.IContainer) (u *QueryHandler) {
	return &QueryHandler{container: container}
}

func (u *QueryHandler) GetHandle(ctx *context.Context) (r interface{}) {
	// return `<?xml version="1.0" encoding="ISO-8859-1"?>
	// <books><book><author>Jack Herrington</author><title>PHP Hacks</title><publisher>O'Reilly</publisher></book><book><author>Jack Herrington</author><title>Podcasting Hacks</title><publisher>O'Reilly</publisher></book><book><author>王小为</author><title>深入在线工具</title><publisher>aTool.org组织</publisher></book></books>`
	//return "success"
	// return map[string]interface{}{
	// 	"a": "success",
	// }
	return []string{"a", "b"}
}
func (u *QueryHandler) Handle(ctx *context.Context) (r interface{}) {
	return "success"
}
