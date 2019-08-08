package conf

//CacheConf 消息队列配置
type CacheConf map[string]interface{}

//RedisCacheConf redis消息队列
type RedisCacheConf CacheConf

//GetProto 获取协议
func (q CacheConf) GetProto() string {
	return q["proto"].(string)
}

//NewRedisCacheConf 构建redis消息队列配置
func NewRedisCacheConf(address []string, dbIndex int) RedisCacheConf {
	conf := make(map[string]interface{})
	conf["proto"] = "redis"
	conf["address"] = address
	conf["db"] = dbIndex
	return conf
}

//WithTimeout 设置超时时长
func (m RedisCacheConf) WithTimeout(dialTimeout int, readTimeout int, writeTimeout int, poolSize int) RedisCacheConf {
	m["dial_timeout"] = dialTimeout
	m["read_timeout"] = readTimeout
	m["write_timeout"] = writeTimeout
	m["pool_size"] = poolSize
	return m
}
