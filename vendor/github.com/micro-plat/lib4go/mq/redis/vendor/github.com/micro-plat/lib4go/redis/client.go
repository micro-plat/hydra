package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/micro-plat/lib4go/types"
)

//ClientConf redis客户端配置
type ClientConf struct {
	MasterName  string   `json:"master"`
	Address     []string `json:"addrs"`
	Password    string   `json:"password"`
	Db          int      `json:"db"`
	DialTimeout int      `json:"dial_timeout"`
	RTimeout    int      `json:"read_timeout"`
	WTimeout    int      `json:"write_timeout"`
	PoolSize    int      `json:"pool_size"`
}

//Client redis client
type Client struct {
	redis.UniversalClient
}

//ClientOption 配置选项
type ClientOption func(*ClientConf)

//WithAddress 设置哨兵服务器
func WithAddress(address []string) ClientOption {
	return func(o *ClientConf) {
		o.Address = address
	}
}

//WithPassword 设置服务器登录密码
func WithPassword(password string) ClientOption {
	return func(o *ClientConf) {
		o.Password = password
	}
}

//WithDB 设置数据库
func WithDB(db int) ClientOption {
	return func(o *ClientConf) {
		o.Db = db
	}
}

//WithDialTimeout 设置连接超时时长
func WithDialTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.DialTimeout = timeout
	}
}

//WithRTimeout 设置读写超时时长
func WithRTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.RTimeout = timeout
	}
}

//WithWTimeout 设置读写超时时长
func WithWTimeout(timeout int) ClientOption {
	return func(o *ClientConf) {
		o.WTimeout = timeout
	}
}

//NewClient 构建客户端
func NewClient(master string, option ...ClientOption) (r *Client, err error) {
	conf := &ClientConf{}
	for _, opt := range option {
		opt(conf)
	}
	conf.MasterName = master
	return NewClientByConf(conf)

}

//NewClientByJSON 根据json构建failover客户端
func NewClientByJSON(j string) (r *Client, err error) {
	conf := &ClientConf{}
	err = json.Unmarshal([]byte(j), &conf)
	if err != nil {
		return nil, err
	}
	return NewClientByConf(conf)
}

//NewClientByConf 根据配置对象构建客户端
func NewClientByConf(conf *ClientConf) (client *Client, err error) {
	conf.DialTimeout = types.DecodeInt(conf.DialTimeout, 0, 3, conf.DialTimeout)
	conf.RTimeout = types.DecodeInt(conf.RTimeout, 0, 3, conf.RTimeout)
	conf.WTimeout = types.DecodeInt(conf.WTimeout, 0, 3, conf.WTimeout)
	conf.PoolSize = types.DecodeInt(conf.PoolSize, 0, 3, conf.PoolSize)
	client = &Client{}
	opts := &redis.UniversalOptions{
		MasterName:   conf.MasterName,
		Addrs:        conf.Address,
		Password:     conf.Password,
		DB:           conf.Db,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.RTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WTimeout) * time.Second,
		PoolSize:     conf.PoolSize,
	}
	client.UniversalClient = redis.NewUniversalClient(opts)
	_, err = client.UniversalClient.Ping().Result()
	return
}
