package redis

import "encoding/json"

//Option 配置选项
type Option func(*Redis)

//WithConfigName 设置数据库分片索引
func WithConfigName(configName string) Option {
	return func(a *Redis) {
		a.ConfigName = configName
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *Redis) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
