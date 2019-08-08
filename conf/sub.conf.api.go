package conf

import "strings"

type Hosts []string

//APIServerConf api server配置信息
type APIServerConf struct {
	Address   string `json:"address,omitempty" valid:"dialstring"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)"`
	RTimeout  int    `json:"readTimeout,omitempty"`
	WTimeout  int    `json:"writeTimeout,omitempty"`
	RHTimeout int    `json:"readHeaderTimeout,omitempty"`
	Hosts     string `json:"host,omitempty"`
	Trace     bool   `json:"trace,omitempty"`
}

//NewAPIServerConf 构建api server配置信息
func NewAPIServerConf(address string) *APIServerConf {
	return &APIServerConf{
		Address: address,
	}
}

//WithTrace 构建api server配置信息
func (a *APIServerConf) WithTrace() *APIServerConf {
	a.Trace = true
	return a
}

//WithTimeout 构建api server配置信息
func (a *APIServerConf) WithTimeout(rtimeout int, wtimout int) *APIServerConf {
	a.RTimeout = rtimeout
	a.WTimeout = wtimout
	return a
}

//WithHeaderReadTimeout 构建api server配置信息
func (a *APIServerConf) WithHeaderReadTimeout(htimeout int) *APIServerConf {
	a.RHTimeout = htimeout
	return a
}

//WithHost 设置host
func (a *APIServerConf) WithHost(host ...string) *APIServerConf {
	a.Hosts = strings.Join(host, ";")
	return a
}

//WithDisable 禁用任务
func (a *APIServerConf) WithDisable() *APIServerConf {
	a.Status = "stop"
	return a
}

//WithEnable 启用任务
func (a *APIServerConf) WithEnable() *APIServerConf {
	a.Status = "start"
	return a
}
