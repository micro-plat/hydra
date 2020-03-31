package conf

type CronServerConf struct {
	Status   string `json:"status,omitempty" valid:"in(start|stop)"`
	Sharding int    `json:"sharding,omitempty"`
	Trace    bool   `json:"trace,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
}

//NewCronServerConf 构建mqc server配置，默认为对等模式
func NewCronServerConf() *CronServerConf {
	return &CronServerConf{}
}

//WithTrace 构建api server配置信息
func (a *CronServerConf) WithTrace() *CronServerConf {
	a.Trace = true
	return a
}

//WitchMasterSlave 设置为主备模式
func (a *CronServerConf) WitchMasterSlave() *CronServerConf {
	a.Sharding = 1
	return a
}

//WitchSharding 设置为分片模式
func (a *CronServerConf) WitchSharding(i int) *CronServerConf {
	a.Sharding = i
	return a
}

//WitchP2P 设置为对等模式
func (a *CronServerConf) WitchP2P() *CronServerConf {
	a.Sharding = 0
	return a
}

//WithTimeout 构建api server配置信息
func (a *CronServerConf) WithTimeout(timeout int) *CronServerConf {
	a.Timeout = timeout
	return a
}

//WithDisable 禁用任务
func (a *CronServerConf) WithDisable() *CronServerConf {
	a.Status = "stop"
	return a
}

//WithEnable 启用任务
func (a *CronServerConf) WithEnable() *CronServerConf {
	a.Status = "start"
	return a
}
