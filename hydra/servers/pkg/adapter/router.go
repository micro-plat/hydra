package adapter

//IRouter IRouter
type IRouter interface {
	GetActions() []string
	GetPath() string
	GetService() string
}
