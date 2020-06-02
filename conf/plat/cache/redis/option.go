package redis

type option struct {
}

//Option 配置选项
type Option func(*Redis)

//WithDbIndex 设置数据库分片索引
func WithDbIndex(i int) Option {
	return func(a *Redis) {
		a.DbIndex = i
	}
}

//WithTimeout 设置数据库连接超时，读写超时时间
func WithTimeout(dialTimeout int, readTimeout int, writeTimeout int) Option {
	return func(a *Redis) {
		a.DialTimeout = dialTimeout
		a.ReadTimeout = readTimeout
		a.WriteTimeout = writeTimeout
	}
}

//WithPoolSize 设置数据库连接池大小
func WithPoolSize(i int) Option {
	return func(a *Redis) {
		a.PoolSize = i
	}
}
