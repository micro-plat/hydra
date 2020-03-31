package conf

//MQCServerConf mqc服务配置
type MQCServerConf struct {
	Status   string `json:"status,omitempty" valid:"in(start|stop)"`
	Sharding int    `json:"sharding,omitempty"`
	Trace    bool   `json:"trace,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
}

//NewMQCServerConf 构建mqc server配置，默认为对等模式
func NewMQCServerConf() *MQCServerConf {
	return &MQCServerConf{Sharding: 1}
}

//WithTrace 构建api server配置信息
func (a *MQCServerConf) WithTrace() *MQCServerConf {
	a.Trace = true
	return a
}

//WithMasterSlave 设置为主备模式
func (a *MQCServerConf) WithMasterSlave() *MQCServerConf {
	a.Sharding = 1
	return a
}

//WithSharding 设置为分片模式
func (a *MQCServerConf) WithSharding(i int) *MQCServerConf {
	a.Sharding = i
	return a
}

//WithP2P 设置为对等模式
func (a *MQCServerConf) WithP2P() *MQCServerConf {
	a.Sharding = 0
	return a
}

//WithTimeout 构建api server配置信息
func (a *MQCServerConf) WithTimeout(timeout int) *MQCServerConf {
	a.Timeout = timeout
	return a
}

//WithDisable 禁用任务
func (a *MQCServerConf) WithDisable() *MQCServerConf {
	a.Status = "stop"
	return a
}

//WithEnable 启用任务
func (a *MQCServerConf) WithEnable() *MQCServerConf {
	a.Status = "start"
	return a
}
