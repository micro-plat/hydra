package conf

//RPCServerConf api server配置信息
type RPCServerConf struct {
	Address   string `json:"address,omitempty" valid:"dialstring"`
	Status    string `json:"status,omitempty" valid:"in(start|stop)"`
	RTimeout  int    `json:"readTimeout,omitempty"`
	WTimeout  int    `json:"writeTimeout,omitempty"`
	RHTimeout int    `json:"readHeaderTimeout,omitempty"`
	Trace     bool   `json:"trace,omitempty"`
	Host      string `json:"host,omitempty"`
	Domain    string `json:"dn,omitempty"`
}

//NewRPCServerConf 构建api server配置信息
func NewRPCServerConf(address string) *RPCServerConf {
	return &RPCServerConf{
		Address: address,
	}
}

//WithTrace 构建api server配置信息
func (a *RPCServerConf) WithTrace() *RPCServerConf {
	a.Trace = true
	return a
}

//WithDisable 禁用任务
func (a *RPCServerConf) WithDisable() *RPCServerConf {
	a.Status = "stop"
	return a
}

//WithEnable 启用任务
func (a *RPCServerConf) WithEnable() *RPCServerConf {
	a.Status = "start"
	return a
}

//WithDNS 设置请求域名
func (a *RPCServerConf) WithDNS(host string, ip string) *RPCServerConf {
	a.Host = host
	a.Domain = ip
	return a
}
