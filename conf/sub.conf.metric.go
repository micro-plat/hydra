package conf

type Metric struct {
	Host     string `json:"host" valid:"requrl,required"`
	DataBase string `json:"dataBase" valid:"ascii,required"`
	Cron     string `json:"cron" valid:"ascii,required"`
	UserName string `json:"userName,omitempty" valid:"ascii"`
	Password string `json:"password,omitempty" valid:"ascii"`
	Disable  bool   `json:"disable,omitempty"`
}

//NewMetric 构建api server配置信息
func NewMetric(host string, db string, cron string) *Metric {
	return &Metric{
		Host:     host,
		DataBase: db,
		Cron:     cron,
	}
}

//WithUserName 设置用户名密码
func (m *Metric) WithUserName(uname string, p string) *Metric {
	m.UserName = uname
	m.Password = p
	return m
}

//WithEnable 启用配置
func (m *Metric) WithEnable() *Metric {
	m.Disable = false
	return m
}

//WithDisable 禁用配置
func (m *Metric) WithDisable() *Metric {
	m.Disable = false
	return m
}
