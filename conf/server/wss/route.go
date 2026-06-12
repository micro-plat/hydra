package wss

import (
	"errors"

	"github.com/micro-plat/hydra/conf"
)

const RouteNodeName = "routes"

type Route struct {
	Host        string `json:"host,omitempty" toml:"host,omitempty"`
	PathPrefix  string `json:"pathPrefix,omitempty" toml:"pathPrefix,omitempty"`
	StripPrefix string `json:"stripPrefix,omitempty" toml:"stripPrefix,omitempty"`
	Group       string `json:"group,omitempty" toml:"group,omitempty"`
}

type Routes struct {
	Routes []Route `json:"routes,omitempty" toml:"routes,omitempty"`
}

func NewRoutes(routes ...Route) *Routes {
	return &Routes{Routes: routes}
}

func GetRoutes(cnf conf.IServerConf) (*Routes, error) {
	routes := &Routes{}
	_, err := cnf.GetSubObject(RouteNodeName, routes)
	if errors.Is(err, conf.ErrNoSetting) {
		return NewRoutes(), nil
	}
	if err != nil {
		return nil, err
	}
	return routes, nil
}
