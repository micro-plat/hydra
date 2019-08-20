package conf

//CacheConf 消息队列配置
type CacheConf map[string]interface{}

//RedisCacheConf redis消息队列
type RedisCacheConf = CacheConf

//GetProto 获取协议
func (q CacheConf) GetProto() string {
	return q["proto"].(string)
}

//NewRedisCacheConf 构建redis消息队列配置
func NewRedisCacheConf(dbIndex int, address ...string) RedisCacheConf {
	conf := make(map[string]interface{})
	conf["proto"] = "redis"
	conf["addrs"] = address
	conf["db"] = dbIndex
	conf["dial_timeout"] = 10
	conf["read_timeout"] = 10
	conf["write_timeout"] = 10
	conf["pool_size"] = 10
	return conf
}

//NewRedisCacheConfForProd 创建prod redis配置
func NewRedisCacheConfForProd(dbIndex int, name ...string) RedisCacheConf {
	kn := "#redisCache"
	if len(name) > 0 {
		kn = name[0]
	}
	return NewRedisCacheConf(dbIndex, kn)
}

//WithPoolSize 修改连接数
func (m RedisCacheConf) WithPoolSize(poolSize int) RedisCacheConf {
	m["pool_size"] = poolSize
	return m
}

//WithTimeout 设置超时时长
func (m RedisCacheConf) WithTimeout(dialTimeout int, readTimeout int, writeTimeout int, poolSize int) RedisCacheConf {
	m["dial_timeout"] = dialTimeout
	m["read_timeout"] = readTimeout
	m["write_timeout"] = writeTimeout
	m["pool_size"] = poolSize
	return m
}
