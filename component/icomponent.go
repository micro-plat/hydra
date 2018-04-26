package component

import "github.com/micro-plat/hydra/context"

//IComponent 提供一组或多组服务的组件
type IComponent interface {
	AddCustomerService(service string, h interface{}, groupNames ...string)
	AddTagPageService(service string, h interface{}, pages ...string)
	AddPageService(service string, h interface{}, pages ...string)
	AddAutoflowService(service string, h interface{})
	AddMicroService(service string, h interface{})
	IsMicroService(service string) bool
	IsAutoflowService(service string) bool
	IsPageService(service string) bool
	IsCustomerService(service string, group ...string) bool
	GetFallbackHandlers() map[string]interface{}
	AddFallbackHandlers(map[string]interface{})
	LoadServices() error

	GetGroupServices(group ...string) []string
	GetServices() []string

	GetGroups(service string) []string
	GetPages(service string) []string

	Fallback(name string, engine string, service string, c *context.Context) (rs interface{})
	Handling(name string, engine string, service string, c *context.Context) (rs interface{})
	Handled(name string, engine string, service string, c *context.Context) (rs interface{})
	Handle(name string, engine string, service string, c *context.Context) interface{}
	Close() error
}

type CloseHandler interface {
	Close() error
}

//Handler context handler
type Handler interface {
	Handle(name string, engine string, service string, c *context.Context) interface{}
}

type GetHandler interface {
	GetHandle(name string, engine string, service string, c *context.Context) interface{}
}
type PostHandler interface {
	PostHandle(name string, engine string, service string, c *context.Context) interface{}
}
type DeleteHandler interface {
	DeleteHandle(name string, engine string, service string, c *context.Context) interface{}
}
type PutHandler interface {
	PutHandle(name string, engine string, service string, c *context.Context) interface{}
}

//FallbackHandler context handler
type FallbackHandler interface {
	Fallback(name string, engine string, service string, c *context.Context) interface{}
}

//GetFallbackHandler context handler
type GetFallbackHandler interface {
	GetFallback(name string, engine string, service string, c *context.Context) interface{}
}

//PostFallbackHandler context handler
type PostFallbackHandler interface {
	PostFallback(name string, engine string, service string, c *context.Context) interface{}
}

//PutFallbackHandler context handler
type PutFallbackHandler interface {
	PutFallback(name string, engine string, service string, c *context.Context) interface{}
}

//DeleteFallbackHandler context handler
type DeleteFallbackHandler interface {
	DeleteFallback(name string, engine string, service string, c *context.Context) interface{}
}

type FallbackServiceFunc func(name string, engine string, service string, c *context.Context) (rs interface{})

func (h FallbackServiceFunc) Fallback(name string, engine string, service string, c *context.Context) (rs interface{}) {
	return h(name, engine, service, c)
}

type ComponentFunc func(c IContainer) error

func (h ComponentFunc) Handle(c IContainer) error {
	return h(c)
}

type ServiceFunc func(name string, engine string, service string, c *context.Context) (rs interface{})

func (h ServiceFunc) Handle(name string, engine string, service string, c *context.Context) (rs interface{}) {
	return h(name, engine, service, c)
}
