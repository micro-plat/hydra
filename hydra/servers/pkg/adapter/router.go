package adapter

//IRouter IRouter
type IRouter interface {
	GetActions() []string
	GetPath() string
	GetService() string
}

//RouteInfo RouteInfo
type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

type RoutesInfo []RouteInfo
