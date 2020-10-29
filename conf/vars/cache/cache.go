package cache

//Cache 消息队列配置
type Cache struct {
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}
