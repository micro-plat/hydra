package metric

type Metric struct {
	Host     string `json:"host" valid:"requrl,required"`
	DataBase string `json:"dataBase" valid:"ascii,required"`
	Cron     string `json:"cron" valid:"ascii,required"`
	UserName string `json:"userName,omitempty" valid:"ascii"`
	Password string `json:"password,omitempty" valid:"ascii"`
	Disable  bool   `json:"disable,omitempty"`
	*option
}

//NewMetric 构建api server配置信息
func NewMetric(host string, db string, cron string, opts ...Option) *Metric {
	m := &Metric{
		Host:     host,
		DataBase: db,
		Cron:     cron,
		option:   &option{},
	}
	for _, opt := range opts {
		opt(m.option)
	}
	return m
}
