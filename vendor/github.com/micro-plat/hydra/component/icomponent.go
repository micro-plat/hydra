package component

import "github.com/micro-plat/hydra/context"

//IComponent 提供一组或多组服务的组件
type IComponent interface {
	AddCustomerService(service string, h interface{}, groupName string, tags ...string)
	IsCustomerService(service string, group ...string) bool
	GetFallbackHandlers() map[string]interface{}
	AddFallbackHandlers(map[string]interface{})
	LoadServices() error

	GetGroupServices(group ...string) []string
	GetServices() []string

	GetGroups(service string) []string
	GetTags(service string) []string

	Fallback(c *context.Context) (rs interface{})
	Handling(c *context.Context) (rs interface{})
	Handled(c *context.Context) (rs interface{})
	Handle(c *context.Context) interface{}

	GetMeta(key string) interface{}
	SetMeta(key string, value interface{})

	Close() error
}

type CloseHandler interface {
	Close() error
}

//Handler context handler
type Handler interface {
	Handle(c *context.Context) interface{}
}

type GetHandler interface {
	GetHandle(c *context.Context) interface{}
}
type PostHandler interface {
	PostHandle(c *context.Context) interface{}
}
type DeleteHandler interface {
	DeleteHandle(c *context.Context) interface{}
}
type PutHandler interface {
	PutHandle(c *context.Context) interface{}
}

//FallbackHandler context handler
type FallbackHandler interface {
	Fallback(c *context.Context) interface{}
}

//GetFallbackHandler context handler
type GetFallbackHandler interface {
	GetFallback(c *context.Context) interface{}
}

//PostFallbackHandler context handler
type PostFallbackHandler interface {
	PostFallback(c *context.Context) interface{}
}

//PutFallbackHandler context handler
type PutFallbackHandler interface {
	PutFallback(c *context.Context) interface{}
}

//DeleteFallbackHandler context handler
type DeleteFallbackHandler interface {
	DeleteFallback(c *context.Context) interface{}
}

type FallbackServiceFunc func(c *context.Context) (rs interface{})

func (h FallbackServiceFunc) Fallback(c *context.Context) (rs interface{}) {
	return h(c)
}

type ComponentFunc func(c IContainer) error

func (h ComponentFunc) Handle(c IContainer) error {
	return h(c)
}

type ServiceFunc func(c *context.Context) (rs interface{})

func (h ServiceFunc) Handle(c *context.Context) (rs interface{}) {
	return h(c)
}
