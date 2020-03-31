package conf

//QueueConf 消息队列配置
type QueueConf map[string]interface{}

//LMQConf 本地消息队列
type LMQConf = QueueConf

//MQTTConf MQTT配置
type MQTTConf = QueueConf

//RedisQueueConf redis消息队列
type RedisQueueConf = QueueConf

//GetProto 获取协议
func (q QueueConf) GetProto() string {
	return q["proto"].(string)
}

//NewLMQConf 本地消息队列
func NewLMQConf() LMQConf {
	conf := make(map[string]interface{})
	conf["proto"] = "lmq"
	return conf
}

//NewMQTTConf 构建MQTT消息队列配置
func NewMQTTConf(address string, userName string, password string) MQTTConf {
	conf := make(map[string]interface{})
	conf["proto"] = "mqtt"
	conf["address"] = address
	conf["userName"] = userName
	conf["password"] = password
	return conf
}

//WithCert 添加证书
func (m MQTTConf) WithCert(cert string) MQTTConf {
	m["cert"] = cert
	return m
}

//NewRedisQueueConf 构建redis消息队列配置
func NewRedisQueueConf(address []string, dbIndex int) RedisQueueConf {
	conf := make(map[string]interface{})
	conf["proto"] = "redis"
	conf["addrs"] = address
	conf["db"] = dbIndex
	return conf
}

//NewRedisQueueConfForProd 创建prod redis配置
func NewRedisQueueConfForProd(dbIndex int, name ...string) RedisQueueConf {
	kn := "#redisCache"
	if len(name) > 0 {
		kn = name[0]
	}
	return NewRedisQueueConf([]string{kn}, dbIndex)
}

//WithTimeout 设置超时时长
func (m RedisQueueConf) WithTimeout(dialTimeout int, readTimeout int, writeTimeout int, poolSize int) RedisQueueConf {
	m["dial_timeout"] = dialTimeout
	m["read_timeout"] = readTimeout
	m["write_timeout"] = writeTimeout
	m["pool_size"] = poolSize
	return m
}
