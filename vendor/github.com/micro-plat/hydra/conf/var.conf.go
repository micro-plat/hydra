package conf

type DBConf struct {
	Provider   string `json:"provider" valid:"required"`
	ConnString string `json:"connString" valid:"required"`
	MaxOpen    int    `json:"maxOpen" valid:"required"`
	MaxIdle    int    `json:"maxIdle" valid:"required"`
	LefeTime   int    `json:"lifeTime" valid:"required"`
}
type CacheConf struct {
	Proto string `json:"proto" valid:"ascii,required"`
}
type QueueConf struct {
	Proto string `json:"proto" valid:"ascii,required"`
}
