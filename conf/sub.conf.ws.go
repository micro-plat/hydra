package conf

//WSServerConf api server配置信息
type WSServerConf struct {
	Address   string `json:"address,omitempty" valid:"dialstring"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)"`
	RTimeout  int    `json:"readTimeout,omitempty"`
	WTimeout  int    `json:"writeTimeout,omitempty"`
	RHTimeout int    `json:"readHeaderTimeout,omitempty"`
	Trace     bool   `json:"trace,omitempty"`
}

//NewWSServerConf 构建api server配置信息
func NewWSServerConf(address string) *WSServerConf {
	return &WSServerConf{
		Address: address,
	}
}

//WithTrace 构建api server配置信息
func (a *WSServerConf) WithTrace() *WSServerConf {
	a.Trace = true
	return a
}

//WithDisable 禁用任务
func (a *WSServerConf) WithDisable() *WSServerConf {
	a.Status = "stop"
	return a
}

//WithEnable 启用任务
func (a *WSServerConf) WithEnable() *WSServerConf {
	a.Status = "start"
	return a
}
