package redis

import "encoding/json"

//Conf redis客户端配置
type option struct {
	MasterName  string   `json:"master"`
	Address     []string `json:"addrs"`
	Password    string   `json:"password"`
	Db          int      `json:"db"`
	DialTimeout int      `json:"dial_timeout"`
	RTimeout    int      `json:"read_timeout"`
	WTimeout    int      `json:"write_timeout"`
	PoolSize    int      `json:"pool_size"`
}

//Option 配置选项
type Option func(*option)

//WithMaster 设置哨兵服务器
func WithMaster(master string) Option {
	return func(o *option) {
		o.MasterName = master
	}
}

//WithAddress 设置哨兵服务器
func WithAddress(address []string) Option {
	return func(o *option) {
		o.Address = address
	}
}

//WithPassword 设置服务器登录密码
func WithPassword(password string) Option {
	return func(o *option) {
		o.Password = password
	}
}

//WithDB 设置数据库
func WithDB(db int) Option {
	return func(o *option) {
		o.Db = db
	}
}

//WithDialTimeout 设置连接超时时长
func WithDialTimeout(timeout int) Option {
	return func(o *option) {
		o.DialTimeout = timeout
	}
}

//WithRTimeout 设置读写超时时长
func WithRTimeout(timeout int) Option {
	return func(o *option) {
		o.RTimeout = timeout
	}
}

//WithWTimeout 设置读写超时时长
func WithWTimeout(timeout int) Option {
	return func(o *option) {
		o.WTimeout = timeout
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *option) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
