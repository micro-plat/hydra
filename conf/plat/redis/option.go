package redis

type option struct {
	Proto        string `json:"proto,omitempty"`
	DbIndex      int    `json:"db,omitempty"`
	DialTimeout  int    `json:"dial_timeout,omitempty"`
	ReadTimeout  int    `json:"read_timeout,omitempty"`
	WriteTimeout int    `json:"write_timeout,omitempty"`
	PoolSize     int    `json:"pool_size,omitempty"`
}

//Option 配置选项
type Option func(*option)

//WithDbIndex 设置数据库分片索引
func WithDbIndex(i int) Option {
	return func(a *option) {
		a.DbIndex = i
	}
}

//WithTimeout 设置数据库连接超时，读写超时时间
func WithTimeout(dialTimeout int, readTimeout int, writeTimeout int) Option {
	return func(a *option) {
		a.DialTimeout = dialTimeout
		a.ReadTimeout = readTimeout
		a.WriteTimeout = writeTimeout
	}
}

//WithPoolSize 设置数据库连接池大小
func WithPoolSize(i int) Option {
	return func(a *option) {
		a.PoolSize = i
	}
}
