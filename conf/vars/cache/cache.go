package cache

//Cache 消息队列配置
type Cache struct {
	Proto string `json:"proto,required"  toml:"proto,required"`
	Raw   []byte `json:"-"`
}
